package launcher

import (
	"context"
	"errors"
	infra_logger "example/admin/bff/internal/infra/logger"
	"fmt"
	goroutiner "github.com/selyukovn/go-routiner"
	"os"
	"os/signal"
	"syscall"
)

type Server = struct {
	Name    string
	FnStart func(context.Context) error
	FnStop  func(context.Context) error
}

func LaunchServers(infraLogger *infra_logger.Logger, servers []Server) {
	// "На старт..."
	// --------------------------------

	noPanicGrt := goroutiner.New(goroutiner.MwPanicToError(func(pv any, ds []byte, ctx context.Context) error {
		infraLogger.CtxPanicFf(ctx, pv, ds, "launcher")
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
						infraLogger.GeneralInfoFf("Запуск %s ...", name)
						if err := servers[i].FnStart(ctx); err != nil {
							infraLogger.GeneralErrorFf("Ошибка при запуске %s: %s - %#v", name, err, err)
							return err
						}
						infraLogger.GeneralInfoFf("%s больше не запушен!", name)
						return nil
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
					infraLogger.GeneralInfoFf("Остановка %s ...", name)
					if err := servers[i].FnStop(ctx); err != nil {
						infraLogger.GeneralErrorFf("Ошибка при остановке %s: %s - %#v", name, err, err)
					}
					infraLogger.GeneralInfoFf("%s остановлен!", name)
					return nil
				}, nil
			}).
			Wait()

		if err := errors.Join(errs...); err != nil {
			infraLogger.GeneralErrorFf("Ошибки завершения: %s - %#v", err, err)
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

		infraLogger.GeneralErrorFf("Ошибка запуска: %s - %#v", startErr, startErr)
		// Старт валится при первой ошибке, а значит какие-то компоненты могут быть запущенными.
		// Их нужно корректно остановить -- поэтому shutdown выполняется и в этом случае.
		// Возможно, какой-то из упавших / незапущенных компонентов запаникует при попытке его остановить --
		// в этом случае main не завалится из-за отдельных горутин с `goroutiner.MwPanicToError()`.
		fnGracefulShutdown()
	case stopSig := <-stopSigCh:
		infraLogger.GeneralInfoFf("Получен стоп-сигнал %q", stopSig.String())
		fnGracefulShutdown()
	}

	// "Финиш!"
	// --------------------------------

	infraLogger.GeneralInfoFf("Конец!")
}
