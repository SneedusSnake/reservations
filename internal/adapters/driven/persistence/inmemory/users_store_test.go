package inmemory_test

import (
	"testing"
	"github.com/SneedusSnake/Reservations/internal/ports/users"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
)

func TestInMemoryUsersStore(t *testing.T) {
	contract := users.UsersRepositoryContract{
		NewStore:  func() users.UsersRepository {
			return inmemory.NewUsersStore();
		},
	}
	contract.Test(t);
}
