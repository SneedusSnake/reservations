package inmemory

import (
	"fmt"
	"slices"
	"sync"

	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/SneedusSnake/Reservations/utils"
)

type SubjectsStore struct {
	counter int
	subjects reservations.Subjects
	tags map[string][]int
	mu sync.Mutex
}

func NewSubjectsStore() *SubjectsStore {
	return &SubjectsStore{counter: 0, subjects: reservations.Subjects{}, tags: make(map[string][]int)}
}

func (s *SubjectsStore) NextIdentity() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	return s.counter;
}

func (s *SubjectsStore) Add(subject reservations.Subject) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subjects = append(s.subjects, subject)
	return nil;
}

func (s *SubjectsStore) Get(id int) (reservations.Subject, error) {
	for _, subject := range s.subjects {
		if subject.Id == id {
			return subject, nil
		}
	}
	return reservations.Subject{}, fmt.Errorf("Subject with id %d not found", id)
}

func (s *SubjectsStore) List() reservations.Subjects {
	return s.subjects
}

func (s *SubjectsStore) Remove(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for index, subject := range s.subjects {
		if subject.Id == id {
			s.subjects = append(s.subjects[:index], s.subjects[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Subject with id %d was not found", id)
}

func (s *SubjectsStore) AddTag(id int, tag string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	subject, _ := s.Get(id)
	s.tags[tag] = append(s.tags[tag], subject.Id)
	return nil;
}

func (s *SubjectsStore) GetTags(id int) ([]string, error) {
	_, err := s.Get(id)
	if err != nil {
		return []string{}, err
	}
	tags := []string{}

	for tag, subjectIds := range s.tags {
		if slices.Contains(subjectIds, id) {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

func (s *SubjectsStore) GetByTags(tags []string) reservations.Subjects {
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
