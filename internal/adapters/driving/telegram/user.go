package telegram

import (
	"github.com/SneedusSnake/Reservations/internal/application"
	"github.com/SneedusSnake/Reservations/internal/domain/users"
)

type TelegramUser struct {
	TelegramId int64
	users.User
}

type TelegramUsersRepository interface {
	Add(u TelegramUser) error
	Get(tgId int64) (TelegramUser, error)
}

type CreateUser struct {
	Id int64
	Name string
}

type TelegramUserService struct {
	store TelegramUsersRepository
	userService application.UserService
}

func NewTelegramUserService(
	store TelegramUsersRepository,
	userService application.UserService,
) *TelegramUserService {
	return &TelegramUserService{store: store, userService: userService}
}

func (s *TelegramUserService) Get(id int64) (TelegramUser, error) {
	return s.store.Get(id)
}

func (s *TelegramUserService) Create(cmd CreateUser) (TelegramUser, error) {

	user, err := s.userService.Create(application.CreateUser{
		Name: cmd.Name,
	})
	if err != nil {
		return TelegramUser{}, err
	}

	tgUser := TelegramUser{
		TelegramId: cmd.Id,
		User: user,
	}
	err = s.store.Add(tgUser)

	if err != nil {
		return TelegramUser{}, err
	}

	return tgUser, nil
}
