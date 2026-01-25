package reservations

import (
	"strings"
)

type Subject struct {
		Id int
		Name string
}
type Subjects []Subject

func (subjects Subjects) Names() string {
	subjectNames := []string{}
	for _, subject := range subjects {
		subjectNames = append(subjectNames, subject.Name)
	}
	return strings.Join(subjectNames, "\n")
}
