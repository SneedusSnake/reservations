package application

import (
	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	reservationsPort "github.com/SneedusSnake/Reservations/internal/ports/reservations"
)

type SubjectService struct {
	store reservationsPort.SubjectsRepository
}

func NewSubjectService(store reservationsPort.SubjectsRepository) *SubjectService {
	return &SubjectService{store: store}
}

func (h *SubjectService) Create(name string) (reservations.Subject, error) {
	subject := reservations.Subject{
		Id: h.store.NextIdentity(),
		Name: name,
	}
	err := h.store.Add(subject)

	if err != nil {
		return subject, err
	}

	return subject, nil
}

type AddTags struct {
	SubjectId int
	Tags []string
}

func (h *SubjectService) AddTags(cmd AddTags) error {
	for _, tag := range cmd.Tags {
		err := h.store.AddTag(cmd.SubjectId, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *SubjectService) List() reservations.Subjects {
	return h.store.List()
}

func (h *SubjectService) ListTags(subjectId int) ([]string, error) {
	return h.store.GetTags(subjectId)
}
