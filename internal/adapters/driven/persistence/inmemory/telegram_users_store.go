package inmemory

import (
	"fmt"

	"github.com/SneedusSnake/Reservations/internal/domain/users"
)

type TelegramUsersStore struct {
	users users.UsersStore
	links map[int64]int
}

func NewTelegramUsersStore(s users.UsersStore) *TelegramUsersStore {
	return &TelegramUsersStore{s, make(map[int64]int)}
}

func (s *TelegramUsersStore) Add(u users.TelegramUser) error {
	s.links[u.TelegramId] = u.Id

	return nil
}

func (s *TelegramUsersStore) Get(tgId int64) (users.TelegramUser, error) {
	userId, ok := s.links[tgId]

	if !ok {
		return users.TelegramUser{}, fmt.Errorf("No user with telegram id %d was found", tgId)
	}

	u, err := s.users.Get(userId)

	if err != nil {
		return users.TelegramUser{}, err
	}

	return users.TelegramUser{TelegramId: tgId, User: u}, nil
}

