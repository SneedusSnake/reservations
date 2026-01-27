package mysql

import (
	"database/sql"
	"fmt"

	"github.com/SneedusSnake/Reservations/internal/adapters/driving/telegram"
)

type TelegramUsersRepository struct {
	connection *sql.DB
	sequence *sequence
}

func NewTelegramUsersRepository(connection *sql.DB) *TelegramUsersRepository {
	return &TelegramUsersRepository{
		connection: connection,
	}
}

func (s *TelegramUsersRepository) Add(u telegram.TelegramUser) error {
	_, err := s.connection.Exec("INSERT INTO telegram_users(telegram_id, user_id) VALUES (?, ?)", u.TelegramId, u.Id)

	return err
}

func (s *TelegramUsersRepository) Get(tgId int64) (telegram.TelegramUser, error) {
	var u telegram.TelegramUser

	row := s.connection.QueryRow(`
		SELECT u.*, tg.telegram_id FROM users u 
		JOIN telegram_users tg ON u.id = tg.user_id
		WHERE tg.telegram_id = ?
	`, tgId)

	if err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.TelegramId); err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("User with telegram id %d was not found", tgId)
		}

		return u, err
	}

	return u, nil
}
