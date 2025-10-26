package inmemory

import (
	"errors"
	"fmt"

	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/SneedusSnake/Reservations/utils"
)

type SubjectsStore struct {
	counter int
	subjects []reservations.Subject
	tags map[string][]int
}

func NewSubjectsStore() *SubjectsStore {
	return &SubjectsStore{counter: 0, subjects: []reservations.Subject{}, tags: make(map[string][]int)}
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
	subject, _ := s.Get(id)
	s.tags[tag] = append(s.tags[tag], subject.Id)
	return nil;
}

func (s *SubjectsStore) GetByTags(tags []string) []reservations.Subject {
	subjects := make([]reservations.Subject, 0)
	subjectIds := s.tags[tags[0]]

	for _, tag := range tags[1:] {
		subjectIds = utils.Intersect(subjectIds, s.tags[tag])
	}

	for _, id := range subjectIds {
		subject, _ := s.Get(id)
		subjects = append(subjects, subject)
	}

	return subjects
}
