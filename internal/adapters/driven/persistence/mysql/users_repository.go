package mysql

import (
	"fmt"
	"database/sql"
	"github.com/SneedusSnake/Reservations/internal/domain/users"
)

type UsersRepository struct {
	connection *sql.DB
	sequence *sequence
}

func NewUsersRepository(connection *sql.DB) *UsersRepository {
	return &UsersRepository{
		connection: connection,
		sequence: &sequence{
			name: "user_seq",
			connection: connection,
		},
	}
}

func (s *UsersRepository) NextIdentity() (int, error) {
	return s.sequence.Next()
}

func (s *UsersRepository) Add(u users.User) error {
	_, err := s.connection.Exec("INSERT INTO users(id, name, email, password) VALUES (?, ?, ?, ?)", u.Id, u.Name, u.Email, u.Password)

	if err != nil {
		return err
	}

	return nil
}

func (s *UsersRepository) Get(id int) (users.User, error) {
	u := users.User{}

	row := s.connection.QueryRow("SELECT*FROM users WHERE id = ?", id)

	if err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Password); err != nil {
		if err == sql.ErrNoRows {
			return users.User{}, fmt.Errorf("User with id %d was not found", id)
		}

		return users.User{}, err
	}

	return u, nil
}

func (s *UsersRepository) Remove(id int) error {
	return nil;
}
