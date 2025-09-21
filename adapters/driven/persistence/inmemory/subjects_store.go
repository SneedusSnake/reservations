package inmemory

import (
	"errors"
	"fmt"

	"github.com/SneedusSnake/Reservations/domain/reservations"
)

type SubjectsStore struct {
	counter int
	subjects []reservations.Subject
}

func NewSubjectsStore() *SubjectsStore {
	return &SubjectsStore{}
}

func (s *SubjectsStore) NextIdentity() int {
	s.counter++
	return s.counter;
}

func (s *SubjectsStore) Add(subject reservations.Subject) error {
	s.subjects = append(s.subjects, subject)
	return nil;
}

func (s *SubjectsStore) Get(id int) (reservations.Subject, error) {
	for _, subject := range s.subjects {
		if subject.Id == id {
			return subject, nil
		}
	}
	return reservations.Subject{}, errors.New(fmt.Sprintf("Subject with id %d not found", id))
}

func (s *SubjectsStore) Remove(id int) error {
	for index, subject := range s.subjects {
		if subject.Id == id {
			s.subjects = append(s.subjects[:index], s.subjects[index+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Subject with id %d was not found", id))
}

func (s *SubjectsStore) AddTag(id int, tag string) error {
	return nil;
}

func (s *SubjectsStore) GetByTag(tag string) []reservations.Subject {
	return nil;
}
