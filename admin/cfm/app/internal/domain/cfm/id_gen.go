package cfm

import "context"

type IdGeneratorInterface interface {
	// Generate
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Generate(ctx context.Context) (Id, error)
}
