package users

import (
	"testing"

	"github.com/SneedusSnake/Reservations/internal/domain/users"
	"github.com/alecthomas/assert/v2"
)

type UsersRepositoryContract struct {
	NewStore func() UsersRepository
}

var store UsersRepository

func (s UsersRepositoryContract) Test(t *testing.T) {
	store = s.NewStore()

	t.Run("it returns error when user was not found", func(t *testing.T) {
		_, err := store.Get(1234)
		assert.Error(t, err)
	})

	t.Run("it adds a new user to the store", func(t *testing.T) {
		user, err := makeUser("Adam")
		assert.NoError(t, err)
		t.Cleanup(func() {
			store.Remove(user.Id)
		})

		err = store.Add(user)
		assert.NoError(t, err)

		foundUser, err := store.Get(user.Id)
		assert.NoError(t, err)

		assert.Equal(t, user, foundUser)
	})

	t.Run("It cannot add user with same id twice", func(t *testing.T) {
		user, err := makeUser("Eve")
		assert.NoError(t, err)
		t.Cleanup(func() {
			store.Remove(user.Id)
		})

		err = store.Add(user)
		assert.NoError(t, err)

		err = store.Add(user)
		assert.Error(t, err)
	})
}

func makeUser(name string) (users.User, error) {
	id, err := store.NextIdentity()
	if err != nil {
		return users.User{}, err
	}

	return users.User{Id: id, Name: name}, nil
}
