package users

type User struct {
	Id int
	Name string
	Email string
	Password string
}

type UsersStore interface {
	NextIdentity() int
	Add(u User) error
	Get(id int) (User, error)
	Remove(id int) error
}
