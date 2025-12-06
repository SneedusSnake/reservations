package inmemory_test

import (
	"testing"
	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
)

func TestInMemorySubjectsStore(t *testing.T) {
	contract := reservations.SubjectStoreContract{
		NewStore:  func() reservations.SubjectsStore {
			return inmemory.NewSubjectsStore();
		},
	}
	contract.Test(t);
}
