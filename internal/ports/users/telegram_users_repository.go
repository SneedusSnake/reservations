package users

import "github.com/SneedusSnake/Reservations/internal/domain/users"

type TelegramUsersRepository interface {
	Add(u users.TelegramUser) error
	Get(tgId int64) (users.TelegramUser, error)
}
