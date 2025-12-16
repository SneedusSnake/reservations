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
	subjectsStorage.Add(reservations.Subject{Id: 1, Name: "Subject#1"})
	subjectsStorage.Add(reservations.Subject{Id: 2, Name: "Subject#2"})
	subjectsStorage.Add(reservations.Subject{Id: 3, Name: "Subject#3"})
	usersStorage.Add(domain.User{Id: 1, Name: "Alice"})
	usersStorage.Add(domain.User{Id: 2, Name: "Bob"})
	reservations := []reservations.Reservation{
		{Id: 1, UserId: 1, SubjectId: 1, Start: time.Now(), End: time.Now().Add(time.Hour)},
		{Id: 2, UserId: 1, SubjectId: 2, Start: time.Now(), End: time.Now().Add(time.Hour)},
		{Id: 3, UserId: 2, SubjectId: 3, Start: time.Now(), End: time.Now().Add(time.Hour)},
		{Id: 4, UserId: 2, SubjectId: 1, Start: time.Now().Add(time.Hour*-1), End: time.Now().Add(-time.Second)},
	}
	for _, reservation := range reservations {
		reservationsStorage.Add(reservation)
	}

	t.Run("It returns error when no reservation is found", func(t *testing.T) {
		_, err := store.Get(12345)

		assert.Error(t, err)
	})

	t.Run("It fetches reservation read model by Id", func(t *testing.T) {
		reservation, err := store.Get(1)
		assert.NoError(t, err)

		assert.Equal(t, "Alice", reservation.User)
		assert.Equal(t, "Subject#1", reservation.Subject)
		assert.Equal(t, reservations[0].Start, reservation.Start)
		assert.Equal(t, reservations[0].End, reservation.End)
	})

	t.Run("It fetches active reservations list", func(t *testing.T) {
		list, err := store.Active(time.Now())
		assert.NoError(t, err)
		reservations := reservations[:3]
		assert.Equal(t, len(reservations), len(list))

		for _, reservation := range reservations {
			u, _ := usersStorage.Get(reservation.UserId)
			s, _ := subjectsStorage.Get(reservation.SubjectId)
			assert.SliceContains(t, list, readmodel.Reservation{
				Id: reservation.Id,
				User: u.Name,
				Subject: s.Name,
				Start: reservation.Start,
				End: reservation.End,
			})
		}
	})
}

