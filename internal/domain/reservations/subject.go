package reservations

import (
	"fmt"
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

func (subjects Subjects) Find(name string) (Subject, error) {
	for _, s := range subjects {
		if s.Name == name {
			return s, nil
		}
	}
	return Subject{}, fmt.Errorf("No subject with name %s found", name)
}
