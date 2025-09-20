package inmemory_test

import (
	"testing"
	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/SneedusSnake/Reservations/adapters/driven/persistence/inmemory"
)

func TestInMemoryReservationsStore(t *testing.T) {
	contract := reservations.ReservationsRegistryContract{
		NewRegistry:  func() reservations.ReservationsRegistry {
			return inmemory.NewReservationStore();
		},
	}
	contract.Test(t);
}
