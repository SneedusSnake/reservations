package reservations

type Subject struct {
		Id int
		Name string
}

type SubjectsStore interface {
	NextIdentity() int
	Add(s Subject) error
	Get(id int) (Subject, error)
	Remove(id int) error
	AddTag(id int, tag string) error
	GetByTag(tag string) []Subject
}
