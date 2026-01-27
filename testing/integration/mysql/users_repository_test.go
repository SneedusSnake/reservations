package mysql

import (
	"context"
	"testing"

	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/mysql"
	"github.com/SneedusSnake/Reservations/internal/ports/users"
	"github.com/SneedusSnake/Reservations/testing/containers"
	mysqlContainer "github.com/SneedusSnake/Reservations/testing/containers/mysql"
	"github.com/alecthomas/assert/v2"
)

func TestMysqlUsersRepository(t *testing.T) {
	container, err := mysqlContainer.Start(context.Background(), "", containers.Stdout("Mysql"))
	if  err != nil {
		assert.NoError(t, err)
	}
	connection, err := container.Connection()
	if  err != nil {
		assert.NoError(t, err)
	}

	contract := users.UsersRepositoryContract{
		NewStore: func() users.UsersRepository {
			return mysql.NewUsersRepository(connection)
		},
	}

	contract.Test(t)
}
