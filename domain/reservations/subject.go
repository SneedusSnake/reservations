package reservations

type Subject struct {
		Id int
		Name string
}

type SubjectsStore interface {
	NextIdentity() int
	Add(s Subject) error
	Get(id int) (Subject, error)
	List() []Subject
	Remove(id int) error
	AddTag(id int, tag string) error
	GetByTags(tags []string) []Subject
}
