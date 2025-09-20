package reservations

type Subject struct {
		Id int
		Name string
}

type SubjectsStore interface {
	add(s Subject)
	
}
