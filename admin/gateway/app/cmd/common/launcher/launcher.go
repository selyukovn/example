package launcher

import (
	"context"
	"errors"
	"fmt"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std/logger"
	"os"
	"os/signal"
	"syscall"
)

type Server = struct {
	Name    string
	FnStart func(context.Context) error
	FnStop  func(context.Context) error
}

func LaunchServers(servers []Server) {
	// "На старт..."
	// --------------------------------

	noPanicGrt := goroutiner.New(goroutiner.MwPanicToError(func(pv any, ds []byte, ctx context.Context) error {
		logger.PanicFf(ctx, pv, ds, "launcher")
		return fmt.Errorf("panic: %#v; stack: %s", pv, ds)
	}))

	bgCtx := context.Background()

	// "Внимание..."
	// --------------------------------

	// Стартер приложения
	starter := noPanicGrt.
		Batch(bgCtx).
		Add(func(ctx context.Context) error {
			return noPanicGrt.
				Batch(ctx).
				AddRange(len(servers), func(i int) (goroutiner.Goroutine, []goroutiner.Middleware) {
					return func(ctx context.Context) error {
						name := servers[i].Name
						logger.InfoFf(ctx, "Запуск %s ...", name)
						// `CancelOnError` отменяет "соседние" горутины с помощью отмены контекста.
						// Если отмену контекста не обработать, `CancelOnError` не вернет результат.
						// Обработка отмены контекста не должна выполняться внутри `Server.FnStart()` --
						// это ответственность лаунчера, поскольку он контролирует graceful-shutdown-процесс.
						select {
						case <-ctx.Done():
							return ctx.Err()
						// `FnStart` блокирующий -- необходима отдельная горутина, иначе select не начнет выполнение.
						case err := <-noPanicGrt.SingleAsync(ctx, servers[i].FnStart):
							if err != nil {
								logger.ErrorFf(ctx, "Ошибка при запуске %s: %s - %#v", name, err, err)
								return err
							}
							logger.InfoFf(ctx, "%s больше не запушен!", name)
							return nil
						}
					}, nil
				}).
				CancelOnError()
		})

	// Перехватчик стоп-сигнала
	stopSigCh := make(chan os.Signal, 1)
	defer close(stopSigCh)
	signal.Notify(stopSigCh, syscall.SIGTERM, syscall.SIGINT)

	fnGracefulShutdown := func() {
		errs := noPanicGrt.
			Batch(bgCtx).
			AddRange(len(servers), func(i int) (goroutiner.Goroutine, []goroutiner.Middleware) {
				return func(ctx context.Context) error {
					name := servers[i].Name
					logger.InfoFf(ctx, "Остановка %s ...", name)
					if err := servers[i].FnStop(ctx); err != nil {
						logger.ErrorFf(ctx, "Ошибка при остановке %s: %s - %#v", name, err, err)
					}
					logger.InfoFf(ctx, "%s остановлен!", name)
					return nil
				}, nil
			}).
			Wait()

		if err := errors.Join(errs...); err != nil {
			logger.ErrorFf(bgCtx, "Ошибки завершения: %s - %#v", err, err)
		}
	}

	// "Марш!"
	// --------------------------------

	startErrCh := starter.Async()

	select {
	case startErr := <-startErrCh:
		if startErr == nil {
			panic("Запуск серверов не заблокировал main, или сервера не вернули ошибки запуска!")
		}

		logger.ErrorFf(bgCtx, "Ошибка запуска: %s - %#v", startErr, startErr)
		// Старт валится при первой ошибке, а значит какие-то компоненты могут быть запущенными.
		// Их нужно корректно остановить -- поэтому shutdown выполняется и в этом случае.
		// Возможно, какой-то из упавших / незапущенных компонентов запаникует при попытке его остановить --
		// в этом случае main не завалится из-за отдельных горутин с `goroutiner.MwPanicToError()`.
		fnGracefulShutdown()
	case stopSig := <-stopSigCh:
		logger.InfoFf(bgCtx, "Получен стоп-сигнал %q", stopSig.String())
		fnGracefulShutdown()
	}

	// "Финиш!"
	// --------------------------------

	logger.InfoFf(bgCtx, "Конец!")
}
