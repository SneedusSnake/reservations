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

func TestMysqlReservationsReadRepository(t *testing.T) {
	container, err := mysqlContainer.Start(context.Background(), containers.Stdout("Mysql"))
	if  err != nil {
		assert.NoError(t, err)
	}
	connection, err := container.Connection()
	if  err != nil {
		assert.NoError(t, err)
	}

	contract := reservations.ReservationsReadRepositoryContract{
		NewRepository: func() reservations.ReservationsReadRepository {
			return mysql.NewReservationsReadRepository(connection)
		},
	}

	contract.Test(
		t,
		mysql.NewReservationsRepository(connection),
		mysql.NewUsersRepository(connection),
		mysql.NewSubjectsRepository(connection),
	)
}
