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
	subjectsStorage.AddTag(1, "test")
	subjectsStorage.AddTag(3, "test")
	usersStorage.Add(domain.User{Id: 1, Name: "Alice"})
	usersStorage.Add(domain.User{Id: 2, Name: "Bob"})

	t.Run("It returns error when no reservation is found", func(t *testing.T) {
		_, err := store.Get(12345)

		assert.Error(t, err)
	})

	t.Run("It fetches reservation read model by Id", func(t *testing.T) {
		t.Cleanup(func() {deleteReservations(t, reservationsStorage)})
		expected := reservations.Reservation{Id: 1, UserId: 1, SubjectId: 1, Start: time.Now(), End: time.Now().Add(time.Hour)}
		addReservation(t, reservationsStorage, expected)
		actual, err := store.Get(1)
		assert.NoError(t, err)

		assert.Equal(t, "Alice", actual.User)
		assert.Equal(t, "Subject#1", actual.Subject)
		assert.Equal(t, expected.Start, actual.Start)
		assert.Equal(t, expected.End, actual.End)
	})

	t.Run("It fetches active reservations list", func(t *testing.T) {
		t.Cleanup(func() {deleteReservations(t, reservationsStorage)})
		rs := []reservations.Reservation{
			{Id: 1, UserId: 1, SubjectId: 1, Start: time.Now(), End: time.Now().Add(time.Hour)},
			{Id: 2, UserId: 1, SubjectId: 2, Start: time.Now(), End: time.Now().Add(time.Hour)},
			{Id: 3, UserId: 2, SubjectId: 3, Start: time.Now(), End: time.Now().Add(time.Hour)},
			{Id: 4, UserId: 2, SubjectId: 1, Start: time.Now().Add(time.Hour*-1), End: time.Now().Add(-time.Second)},
		}
		addReservations(t, reservationsStorage, rs)
		list, err := store.Active(time.Now())
		assert.NoError(t, err)
		rs = rs[:3]
		assert.Equal(t, len(rs), len(list))

		for _, r := range rs {
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
		t.Cleanup(func() {deleteReservations(t, reservationsStorage)})
		rs := []reservations.Reservation{
			{Id: 1, UserId: 1, SubjectId: 1, Start: time.Now(), End: time.Now().Add(time.Hour)},
			{Id: 2, UserId: 1, SubjectId: 2, Start: time.Now(), End: time.Now().Add(time.Hour)},
			{Id: 3, UserId: 2, SubjectId: 3, Start: time.Now(), End: time.Now().Add(time.Hour)},
		}
		addReservations(t, reservationsStorage, rs)
		list, err := store.Active(time.Now(), "test")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(list))
		assert.Equal(t, "Subject#1", list[0].Subject)
		assert.Equal(t, "Subject#3", list[1].Subject)

	})
}

func deleteReservations(t *testing.T, store ReservationsRepository) {
	rs, err := store.List()
	assert.NoError(t, err)
	for _, r := range rs {
		err := store.Remove(r.Id)
		assert.NoError(t, err)
	}
}

func addReservation(t *testing.T, store ReservationsRepository, r reservations.Reservation) {
	err := store.Add(r)
	assert.NoError(t, err)
}

func addReservations(t *testing.T, store ReservationsRepository, rs reservations.Reservations) {
	for _, r := range rs {
		addReservation(t, store, r)
	}
}

