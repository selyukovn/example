package code

import "context"

type GeneratorInterface interface {
	// Generate
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Generate(ctx context.Context) (Code, error)
}
