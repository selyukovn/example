package code

import "context"

type HasherInterface interface {
	// Hash
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Hash(ctx context.Context, code Code) (Hash, error)

	// Compare
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Compare(ctx context.Context, code Code, hash Hash) (bool, error)
}
