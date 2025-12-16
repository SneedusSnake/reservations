package inmemory_test

import (
	"testing"
	"github.com/SneedusSnake/Reservations/internal/ports/reservations"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
)

func TestInMemoryReservationsReadStore(t *testing.T) {
	subjects := inmemory.NewSubjectsStore()
	users := inmemory.NewUsersStore()
	reservationsStore := inmemory.NewReservationStore()

	contract := reservations.ReservationsReadRepositoryContract{
		NewRepository:  func() reservations.ReservationsReadRepository {
			return inmemory.NewReservationReadStore(reservationsStore, users, subjects);
		},
	}
	contract.Test(t, reservationsStore, users, subjects);
}
