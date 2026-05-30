package resources

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	assert "github.com/selyukovn/go-wm-assert"
)

type Sql = struct {
	Db                    *sql.DB
	FnIsDuplicateKeyError func(error) bool
	FnIsDeadlockError     func(error) bool
}

func OpenMysql(host string, user string, password string, dbName string) Sql {
	assert.Str().NotEmpty().Must(host)
	assert.Str().NotEmpty().Must(user)
	assert.Str().NotEmpty().Must(password)
	assert.Str().NotEmpty().Must(dbName)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, password, host, dbName))
	if err != nil {
		panic(err)
	}

	fnIsDuplicateKeyError := func(err error) bool {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			return mysqlErr.Number == 1062
		}
		return false
	}

	fnIsDeadlockError := func(err error) bool {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			return mysqlErr.Number == 1213
		}
		return false
	}

	return Sql{
		Db:                    db,
		FnIsDuplicateKeyError: fnIsDuplicateKeyError,
		FnIsDeadlockError:     fnIsDeadlockError,
	}
}
