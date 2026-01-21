package reservations

import "github.com/SneedusSnake/Reservations/internal/domain/reservations"

type SubjectsRepository interface {
	NextIdentity() (int, error)
	Add(s reservations.Subject) error
	Get(id int) (reservations.Subject, error)
	List() reservations.Subjects
	Remove(id int) error
	AddTag(id int, tag string) error
	GetTags(id int) ([]string, error)
	GetByTags(tags []string) reservations.Subjects
}
