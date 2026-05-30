package dlq

import "context"

type TopicHolderInterface interface {
	// Hold
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorAlreadyDone
	// 	- std.ErrorRuntime
	Hold(ctx context.Context, topic string) error

	// UnHold
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorAlreadyDone
	// 	- std.ErrorRuntime
	UnHold(ctx context.Context, topic string) error

	// WaitTillOnHold
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	WaitTillOnHold(ctx context.Context, topic string) error
}
