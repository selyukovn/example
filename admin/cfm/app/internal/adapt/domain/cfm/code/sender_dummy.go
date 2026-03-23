package code

import (
	"context"
	"errors"
	"example/admin/cfm/internal/domain/cfm/code"
	"fmt"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"math"
	"math/rand"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ code.SenderInterface = SenderImplDummy{}

type SenderImplDummy struct{}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewSenderImplDummy() SenderImplDummy {
	return SenderImplDummy{}
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
func (s SenderImplDummy) Send(ctx context.Context, code code.Code, email std.Email) error {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(code.IsNil())
	assert.FalseMust(email.IsNil())

	var rMin uint32 = 100
	var rMax uint32 = 5000
	sendingDuration := time.Duration(rMin+rand.Uint32()%(rMax-rMin+1)) * time.Millisecond

	isSuccess := rand.Uint32() > math.MaxUint32/2

	message := fmt.Sprintf(
		"[%T] Send(%v, %v) -- %s (%s)",
		s,
		code,
		email,
		std.Ternary[string](isSuccess, "success", "failed"),
		sendingDuration.String(),
	)

	time.Sleep(sendingDuration)

	logger.InfoFf(ctx, message)

	if !isSuccess {
		return std.WrapErrorToRuntime(errors.New(message), s, "Send")
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
