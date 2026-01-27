package mysql

import (
	"context"
	"testing"

	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/mysql"
	"github.com/SneedusSnake/Reservations/internal/ports/reservations"
	"github.com/SneedusSnake/Reservations/testing/containers"
	mysqlContainer "github.com/SneedusSnake/Reservations/testing/containers/mysql"
	"github.com/alecthomas/assert/v2"
)

func TestMysqlSubjectsRepository(t *testing.T) {
	container, err := mysqlContainer.Start(context.Background(), "", containers.Stdout("Mysql"))
	if  err != nil {
		assert.NoError(t, err)
	}
	connection, err := container.Connection()
	if  err != nil {
		assert.NoError(t, err)
	}

	contract := reservations.SubjectsRepositoryContract{
		NewStore: func() reservations.SubjectsRepository {
			return mysql.NewSubjectsRepository(connection)
		},
	}

	contract.Test(t)
}
