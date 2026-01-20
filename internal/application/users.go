package application

import (
	"github.com/SneedusSnake/Reservations/internal/domain/users"
	ports "github.com/SneedusSnake/Reservations/internal/ports/users"
)

type CreateUser struct {
	Name string
	Email string
	Password string
}

type UserService struct {
	store ports.UsersRepository 
}

func NewUserService(store ports.UsersRepository) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Get(id int) (users.User, error) {
	return s.store.Get(id)
}

func (s *UserService) Create(cmd CreateUser) (users.User, error) {
	id, err := s.store.NextIdentity()
	if err != nil {
		return users.User{}, err
	}

	user := users.User{
		Id: id,
		Name: cmd.Name,
		Email: cmd.Email,
		Password: cmd.Password,
	}

	return user, s.store.Add(user)
}
