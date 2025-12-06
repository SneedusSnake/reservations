package inmemory_test

import (
	"testing"
	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
)

func TestInMemoryReservationsStore(t *testing.T) {
	contract := reservations.ReservationsRegistryContract{
		NewRegistry:  func() reservations.ReservationsRegistry {
			return inmemory.NewReservationStore();
		},
	}
	contract.Test(t);
}
