package session_tick_time

import (
	"context"
	"errors"
	"example/admin/auth/internal/domain/session"
	"example/admin/auth/internal/opera/domain_facades"
	"fmt"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"strconv"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Command struct {
	grt        goroutiner.Goroutiner
	sessDomFac domain_facades.SessionDomFac
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewCommand
//
// Паникует при нулевых аргументах.
func NewCommand(
	grt goroutiner.Goroutiner,
	sessDomFac domain_facades.SessionDomFac,
) Command {
	assert.NotZeroMust(grt)
	assert.Cmp[domain_facades.SessionDomFac]().NotEq(domain_facades.SessionDomFacNil).Must(sessDomFac)

	return Command{
		grt:        grt,
		sessDomFac: sessDomFac,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Execute
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (c Command) Execute(args Args) error {
	assert.FalseMust(args.IsNil())

	ctx := args.Ctx()
	limit := args.Limit()

	// --

	// Находим сессии на тик
	sessIds, err := c.sessDomFac.GetIdsGoingToExpire(ctx, limit)
	switch err.(type) {
	case nil:
	case std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, c, "Execute")
	default:
		panic(err)
	}
	logger.InfoFf(ctx, "Сессии: %v", sessIds)

	sessIdsCount := len(sessIds)
	if sessIdsCount == 0 {
		return nil
	}

	// --

	workersCount := sessIdsCount/10 + std.Ternary[int](sessIdsCount%10 > 0, 1, 0)
	errs := c.grt.
		Batch(ctx).
		AddRange(workersCount, func(i int) (goroutiner.Goroutine, []goroutiner.Middleware) {
			from := i * 10
			to := min(from+10, sessIdsCount)
			wSessIds := sessIds[from:to]
			logger.DebugFf(ctx, "Сессии воркера #%d: %v", i, wSessIds)
			return func(ctx context.Context) error {
				ctx = logger.AddAttrToCtx(ctx, "worker", strconv.Itoa(i))
				return errors.Join(c.executeWorker(ctx, wSessIds)...)
			}, nil
		}).
		Wait()

	if err = errors.Join(errs...); err != nil {
		err = std.WrapErrorToRuntime(err, c, "Execute")
		logger.ErrorFf(ctx, "Ошибки: %v", err)
		return err
	}

	return nil
}

func (c Command) executeWorker(ctx context.Context, wSessIds []session.Id) []error {
	errs := make([]error, 0)
	for _, sessId := range wSessIds {
		err := c.sessDomFac.TickTime(ctx, sessId)
		switch err.(type) {
		case nil:
		case std.ErrorNotFound, session.ErrorClosed, std.ErrorAlreadyDone, std.ErrorRuntime:
			// NotFound -- тоже баг: id есть, а сессии нет.
			// Closed -- тоже баг: нашли как незакрытую, а она закрыта.
			// AlreadyDone -- тоже баг: нашли как протухшую незакрытую, а делать с ней нечего.
			errs = append(errs, std.WrapErrorToRuntime(err, c, "Execute", fmt.Sprintf("sessId: %q", sessId)))
		default:
			panic(err)
		}
	}
	return errs
}

// ---------------------------------------------------------------------------------------------------------------------
