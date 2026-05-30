package kvdb

import (
	"context"
)

type KvDbInterface interface {
	// Ключи можно было бы тоже принимать как `[]byte` для универсальности,
	// но это лишь приведет к постоянным приведениям типов `[]byte(keyStr)` при использовании интерфейса,
	// поскольку в большинстве случаев ключами все же будут являться строки.

	// Insert
	//
	// Ошибки:
	// 	- std.ErrorAlreadyDone
	// 	- std.ErrorRuntime
	Insert(ctx context.Context, key string, value []byte) error

	// Update
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	Update(ctx context.Context, key string, value []byte) error

	// Delete
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	Delete(ctx context.Context, key string) error

	// Get
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	Get(ctx context.Context, key string) ([]byte, error)
}
