package reservations

import (
	"testing"
	"time"

	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	domain "github.com/SneedusSnake/Reservations/internal/domain/users"
	"github.com/SneedusSnake/Reservations/internal/ports/users"
	readmodel "github.com/SneedusSnake/Reservations/internal/read_model"
	"github.com/alecthomas/assert/v2"
)


type ReservationsReadRepositoryContract struct {
	NewRepository func() ReservationsReadRepository
}

func (r ReservationsReadRepositoryContract) Test (t *testing.T, reservationsStorage ReservationsRepository, usersStorage users.UsersRepository, subjectsStorage SubjectsRepository) {
	store := r.NewRepository()
	factory := builder(t, reservationsStorage)
	now, err := time.Parse(time.DateTime, "2025-02-01 14:00:00")
	assert.NoError(t, err)

	subjects := reservations.Subjects{
		reservations.Subject{Id: 1, Name: "Subject#1"},
		reservations.Subject{Id: 2, Name: "Subject#2"},
		reservations.Subject{Id: 3, Name: "Subject#3"},
	}

	for _, s := range subjects {
		err := subjectsStorage.Add(s)
		assert.NoError(t, err)
	}

	subjectsStorage.AddTag(subjects[0].Id, "test")
	subjectsStorage.AddTag(subjects[2].Id, "test")

	users := []domain.User{
		{Id: 1, Name: "Alice"},
		{Id: 2, Name: "Bob"},
	}

	for _, u := range users {
		err := usersStorage.Add(u)
		assert.NoError(t, err)
	}

	cleanUp := factory.CleanUp

	t.Run("It returns error when no reservation is found", func(t *testing.T) {
		_, err := store.Get(12345)

		assert.Error(t, err)
	})

	t.Run("It fetches reservation read model by Id", func(t *testing.T) {
		cleanUp(t)
		user := users[1]
		subject := subjects[0]
		expected := factory.UserId(user.Id).SubjectId(subject.Id).Persist()

		actual, err := store.Get(1)
		assert.NoError(t, err)

		assert.Equal(t, subject.Name, actual.Subject)
		assert.Equal(t, user.Name, actual.User)
		assert.Equal(t, expected.Start, actual.Start)
		assert.Equal(t, expected.End, actual.End)
	})

	t.Run("It fetches active reservations list", func(t *testing.T) {
		cleanUp(t)
		blueprint := factory.UserId(users[0].Id).SubjectId(subjects[0].Id).
			StartsAt(now).
			EndsAt(now.Add(time.Hour))
		expectedReservations := reservations.Reservations{
				blueprint.Persist(),
				blueprint.SubjectId(subjects[1].Id).Persist(),
				blueprint.UserId(users[1].Id).SubjectId(subjects[2].Id).Persist(),
				blueprint.UserId(users[1].Id).Persist(),
		}
		blueprint.StartsAt(now.Add(-time.Hour)).EndsAt(now.Add(-time.Minute)).Persist()

		list, err := store.Active(now)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedReservations), len(list))

		for _, r := range expectedReservations {
			u, _ := usersStorage.Get(r.UserId)
			s, _ := subjectsStorage.Get(r.SubjectId)
			assert.SliceContains(t, list, readmodel.Reservation{
				Id: r.Id,
				User: u.Name,
				Subject: s.Name,
				Start: r.Start,
				End: r.End,
			})
		}
	})

	t.Run("It fetches active reservations list filtered by tags", func(t *testing.T) {
		cleanUp(t)
		blueprint := factory.UserId(users[0].Id).StartsAt(now).EndsAt(now.Add(time.Hour))
		blueprint.SubjectId(subjects[0].Id).Persist()
		blueprint.SubjectId(subjects[1].Id).Persist()
		blueprint.SubjectId(subjects[2].Id).Persist()

		list, err := store.Active(now, "test")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(list))
		assert.Equal(t, "Subject#1", list[0].Subject)
		assert.Equal(t, "Subject#3", list[1].Subject)
	})
}
