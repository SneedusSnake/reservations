package mysql

import (
	"fmt"
	"database/sql"
	"github.com/SneedusSnake/Reservations/internal/domain/users"
)

type UsersRepository struct {
	connection *sql.DB
}

func NewUsersRepository(connection *sql.DB) *UsersRepository {
	return &UsersRepository{connection: connection}
}

func (s *UsersRepository) NextIdentity() int {
	var id int
	s.connection.Exec("UPDATE user_seq SET value = LAST_INSERT_ID(value+1)")
	row := s.connection.QueryRow("SELECT LAST_INSERT_ID() as id")
	if err := row.Scan(&id); err != nil {
		return 0
	}

	return id
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
