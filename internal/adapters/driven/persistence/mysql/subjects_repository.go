package mysql

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
)

type SubjectsRepository struct {
	connection *sql.DB
	sequence *sequence
}

func NewSubjectsRepository(connection *sql.DB) *SubjectsRepository {
	return &SubjectsRepository{
		connection: connection,
		sequence: &sequence{
			name: "subject_seq",
			connection: connection,
		},
	}
}

func (s *SubjectsRepository) NextIdentity() (int, error) {
	return s.sequence.Next()
}

func (s *SubjectsRepository) Add(subject reservations.Subject) error {
	_, err := s.connection.Exec("INSERT INTO subjects(id, name) VALUES (?, ?)", subject.Id, subject.Name)

	return err
}

func (s *SubjectsRepository) Get(id int) (reservations.Subject, error) {
	subject := reservations.Subject{}

	row := s.connection.QueryRow("SELECT*FROM subjects WHERE id = ?", id)

	if err := row.Scan(&subject.Id, &subject.Name); err != nil {
		if err == sql.ErrNoRows {
			return reservations.Subject{}, fmt.Errorf("Subject with id %d was not found", id)
		}

		return reservations.Subject{}, err
	}

	return subject, nil
}

func (s *SubjectsRepository) List() (reservations.Subjects, error) {
	var subjects reservations.Subjects

	rows, err := s.connection.Query("SELECT*FROM subjects")
	if err != nil {
		return subjects, err
	}
	
	for rows.Next() {
		var subject reservations.Subject
		if err = rows.Scan(&subject.Id, &subject.Name); err != nil {
			return subjects, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (s *SubjectsRepository) Remove(id int) error {
	_, err := s.connection.Exec("DELETE FROM subjects WHERE id = ?", id)

	return err
}

func (s *SubjectsRepository) AddTag(id int, tag string) error {
	_, err := s.connection.Exec("INSERT INTO subject_tags VALUES (?, ?)", id, tag)

	return err
}

func (s *SubjectsRepository) GetTags(id int) ([]string, error) {
	var tags []string
	rows, err := s.connection.Query("SELECT tag FROM subject_tags WHERE subject_id = ?", id)
	if err != nil {
		return tags, err
	}

	for rows.Next() {
		var tag string
		if err = rows.Scan(&tag); err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *SubjectsRepository) GetByTags(tags []string) (reservations.Subjects, error) {
	var subjects reservations.Subjects
	slices.Sort(tags)

	rows, err := s.connection.Query(`
		SELECT s.id, s.name, t.tags FROM subjects AS s
		JOIN (SELECT subject_id, GROUP_CONCAT(tag ORDER BY tag) AS tags FROM subject_tags GROUP BY subject_id) t ON t.subject_id = s.id
		WHERE t.tags LIKE ?`,
		"%" + strings.Join(tags, ",") + "%",
	)

	if err != nil {
		return reservations.Subjects{}, err
	}

	for rows.Next() {
		var subject reservations.Subject
		var rowTags string
		err = rows.Scan(&subject.Id, &subject.Name, &rowTags)
		if err != nil {
			return reservations.Subjects{}, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (s *SubjectsRepository) GetByName(name string) (reservations.Subject, error) {
	subject := reservations.Subject{}

	row := s.connection.QueryRow("SELECT*FROM subjects WHERE name = ? FOR UPDATE", name)

	if err := row.Scan(&subject.Id, &subject.Name); err != nil {
		if err == sql.ErrNoRows {
			return reservations.Subject{}, fmt.Errorf("Subject with name %s was not found", name)
		}

		return reservations.Subject{}, err
	}

	return subject, nil
}
