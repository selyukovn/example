package code

import (
	"context"
	"errors"
	"example/admin/cfm/internal/domain/cfm/code"
	"example/admin/cfm/internal/infra/logger"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"math"
	"math/rand"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type SenderImplDummy struct {
	l *logger.Logger
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewSenderImplDummy(l *logger.Logger) *SenderImplDummy {
	return &SenderImplDummy{
		l: l,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Send
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (s *SenderImplDummy) Send(ctx context.Context, code code.Code, email std.Email) error {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(code.IsNil())
	assert.FalseMust(email.IsNil())

	traceId, _ := s.l.GetTraceIdFromCtx(ctx)

	var rMin uint32 = 100
	var rMax uint32 = 5000
	sendingDuration := time.Duration(rMin+rand.Uint32()%(rMax-rMin+1)) * time.Millisecond

	isSuccess := rand.Uint32() > math.MaxUint32/2

	message := fmt.Sprintf(
		"[%T] [traceId=%s] Send(%v, %v) -- %s (%s)",
		s,
		traceId,
		code,
		email,
		std.Ternary[string](isSuccess, "success", "failed"),
		sendingDuration.String(),
	)

	time.Sleep(sendingDuration)

	println(message)

	if !isSuccess {
		return std.WrapErrorToRuntime(errors.New(message), s, "Send")
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
