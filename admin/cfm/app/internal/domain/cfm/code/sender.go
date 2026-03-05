package code

import (
	"context"
	"github.com/selyukovn/go-std"
)

type SenderInterface interface {
	// Send
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Send(ctx context.Context, code Code, email std.Email) error
}
