package inmemory_test

import (
	"testing"
	"github.com/SneedusSnake/Reservations/domain/users"
	"github.com/SneedusSnake/Reservations/adapters/driven/persistence/inmemory"
)

func TestInMemoryUsersStore(t *testing.T) {
	contract := users.UsersStoreContract{
		NewStore:  func() users.UsersStore {
			return inmemory.NewUsersStore();
		},
	}
	contract.Test(t);
}
