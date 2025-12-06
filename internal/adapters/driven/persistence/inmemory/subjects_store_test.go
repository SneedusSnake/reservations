package inmemory_test

import (
	"testing"
	"github.com/SneedusSnake/Reservations/internal/ports/reservations"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
)

func TestInMemorySubjectsStore(t *testing.T) {
	contract := reservations.SubjectsRepositoryContract{
		NewStore:  func() reservations.SubjectsRepository {
			return inmemory.NewSubjectsStore();
		},
	}
	contract.Test(t);
}
