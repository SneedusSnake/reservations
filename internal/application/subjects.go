package application

import "github.com/SneedusSnake/Reservations/internal/domain/reservations"

type subjectsHandler struct {
	store reservations.SubjectsStore
}

type AddSubjectHandler struct {
	subjectsHandler
}

func (h *AddSubjectHandler) Handle(name string) (reservations.Subject, error) {
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

type AddSubjectTags struct {
	subjectId int
	tags []string
}

type AddSubjectTagsHandler struct {
	subjectsHandler
}

func (h *AddSubjectTagsHandler) Handle(cmd AddSubjectTags) error {
	for _, tag := range cmd.tags {
		err := h.store.AddTag(cmd.subjectId, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

type ListSubjectsHandler struct {
	subjectsHandler
}

func (h *ListSubjectsHandler) Handle() reservations.Subjects {
	return h.store.List()
}

type ListSubjectTagsHandler struct {
	subjectsHandler
}

func (h *ListSubjectTagsHandler) Handle(subjectId int) ([]string, error) {
	return h.store.GetTags(subjectId)
}
