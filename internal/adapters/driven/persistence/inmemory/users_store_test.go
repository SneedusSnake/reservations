package inmemory_test

import (
	"testing"
	"github.com/SneedusSnake/Reservations/internal/domain/users"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
)

func TestInMemoryUsersStore(t *testing.T) {
	contract := users.UsersStoreContract{
		NewStore:  func() users.UsersStore {
			return inmemory.NewUsersStore();
		},
	}
	contract.Test(t);
}
