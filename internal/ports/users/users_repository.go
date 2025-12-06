package users

import "github.com/SneedusSnake/Reservations/internal/domain/users"

type UsersRepository interface {
	NextIdentity() int
	Add(u users.User) error
	Get(id int) (users.User, error)
	Remove(id int) error
}
