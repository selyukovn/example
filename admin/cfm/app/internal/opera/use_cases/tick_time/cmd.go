package tick_time

import (
	"context"
	"errors"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/opera/domain_facades"
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
	grt       goroutiner.Goroutiner
	cfmDomFac domain_facades.CfmDomFac
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewCommand
//
// Паникует при нулевых аргументах.
func NewCommand(
	grt goroutiner.Goroutiner,
	cfmDomFac domain_facades.CfmDomFac,
) Command {
	assert.NotZeroMust(grt)
	assert.Cmp[domain_facades.CfmDomFac]().NotEq(domain_facades.CfmDomFacNil).Must(cfmDomFac)

	return Command{
		grt:       grt,
		cfmDomFac: cfmDomFac,
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
func (c Command) Execute(ctx context.Context, limit uint) error {
	assert.NotNilDeepMust(ctx)
	assert.Num[uint]().Positive().Must(limit)

	// Находим конфирмации на тик
	cfmIds, err := c.cfmDomFac.GetIdsGoingToExpire(ctx, limit)
	switch err.(type) {
	case nil:
	case std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, c, "Execute")
	default:
		panic(err)
	}
	logger.InfoFf(ctx, "Конфирмации: %v", cfmIds)

	cfmIdsCount := len(cfmIds)
	if cfmIdsCount == 0 {
		return nil
	}

	// --

	workersCount := cfmIdsCount/10 + std.Ternary[int](cfmIdsCount%10 > 0, 1, 0)
	errs := c.grt.
		Batch(ctx).
		AddRange(workersCount, func(i int) (goroutiner.Goroutine, []goroutiner.Middleware) {
			from := i * 10
			to := min(from+10, cfmIdsCount)
			wCfmIds := cfmIds[from:to]
			logger.DebugFf(ctx, "Конфирмации воркера #%d: %v", i, wCfmIds)
			return func(ctx context.Context) error {
				ctx = logger.AddAttrToCtx(ctx, "worker", strconv.Itoa(i))
				return errors.Join(c.executeWorker(ctx, wCfmIds)...)
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

func (c Command) executeWorker(ctx context.Context, wCfmIds []cfm.Id) []error {
	errs := make([]error, 0)
	for _, cfmId := range wCfmIds {
		err := c.cfmDomFac.TickTime(ctx, cfmId)
		switch err.(type) {
		case nil:
		case std.ErrorNotFound, cfm.ErrorFinished, std.ErrorAlreadyDone, std.ErrorRuntime:
			// NotFound -- тоже баг: id есть, а конфирмации нет.
			// Closed -- тоже баг: нашли как незавершенную, а она завершена.
			// AlreadyDone -- тоже баг: нашли как протухшую незавершенную, а делать с ней нечего.
			errs = append(errs, std.WrapErrorToRuntime(err, c, "Execute", fmt.Sprintf("cfmId: %q", cfmId)))
		default:
			panic(err)
		}
	}
	return errs
}

// ---------------------------------------------------------------------------------------------------------------------
