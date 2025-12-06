package users

type TelegramUser struct {
	TelegramId int64
	User
}

type TelegramUsersStore interface {
	Add(u TelegramUser) error
	Get(tgId int64) (TelegramUser, error)
}
