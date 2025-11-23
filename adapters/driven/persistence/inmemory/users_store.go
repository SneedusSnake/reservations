package inmemory

import (
	"fmt"

	"github.com/SneedusSnake/Reservations/domain/users"
)

type UsersStore struct {
	counter int
	users  []users.User
}

func NewUsersStore() *UsersStore {
	return &UsersStore{}
}

func (s *UsersStore) NextIdentity() int {
	s.counter++
	return s.counter
}

func (s *UsersStore) Add(u users.User) error {
	existingUser, err := s.Get(u.Id)
	if err == nil {
		return fmt.Errorf("User with id %d already exists: %v", u.Id, existingUser)
	}

	s.users = append(s.users, u)

	return nil
}

func (s *UsersStore) Get(id int) (users.User, error) {
	for _, u := range s.users {
		if u.Id == id {
			return u, nil
		}
	}

	return users.User{}, fmt.Errorf("User with id %d was not found", id)
}

func (s *UsersStore) Remove(id int) error {
	return nil
}
