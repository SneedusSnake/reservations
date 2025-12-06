package inmemory_test

import (
	"testing"
	"github.com/SneedusSnake/Reservations/internal/ports/reservations"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
)

func TestInMemoryReservationsStore(t *testing.T) {
	contract := reservations.ReservationsRepositoryContract{
		NewRegistry:  func() reservations.ReservationsRepository {
			return inmemory.NewReservationStore();
		},
	}
	contract.Test(t);
}
