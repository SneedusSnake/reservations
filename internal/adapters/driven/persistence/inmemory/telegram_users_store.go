package inmemory

import (
	"fmt"

	 "github.com/SneedusSnake/Reservations/internal/adapters/driving/telegram"
	"github.com/SneedusSnake/Reservations/internal/ports/users"
)

type TelegramUsersStore struct {
	users users.UsersRepository
	links map[int64]int
}

func NewTelegramUsersStore(s users.UsersRepository) *TelegramUsersStore {
	return &TelegramUsersStore{s, make(map[int64]int)}
}

func (s *TelegramUsersStore) Add(u telegram.TelegramUser) error {
	s.links[u.TelegramId] = u.Id

	return nil
}

func (s *TelegramUsersStore) Get(tgId int64) (telegram.TelegramUser, error) {
	userId, ok := s.links[tgId]

	if !ok {
		return telegram.TelegramUser{}, fmt.Errorf("No user with telegram id %d was found", tgId)
	}

	u, err := s.users.Get(userId)

	if err != nil {
		return telegram.TelegramUser{}, err
	}

	return telegram.TelegramUser{TelegramId: tgId, User: u}, nil
}

