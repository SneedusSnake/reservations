package application_test

import (
	"testing"
	"time"

	"github.com/SneedusSnake/Reservations/adapters/driven/persistence/inmemory"
	"github.com/SneedusSnake/Reservations/application"
	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/SneedusSnake/Reservations/domain/users"
	"github.com/alecthomas/assert/v2"
)

type FakeClock struct {
	now time.Time
}

func (c *FakeClock) Current() time.Time {
	return c.now
}

func (c *FakeClock) Set(t time.Time) {
	c.now = t
}

var registry reservations.ReservationsRegistry

func TestCreateReservationHandler(t *testing.T) {
	subjectsStore := inmemory.NewSubjectsStore()
	usersStore := inmemory.NewUsersStore()
	registry = inmemory.NewReservationStore()
	clock := &FakeClock{}
	clock.Set(time.Now())
	handler := application.NewCreateReservationHandler(
		subjectsStore,
		registry,
		usersStore,
		clock,
	)
	subjects := createTestSubjects(subjectsStore, t)
	users := createTestUsers(usersStore, t)
	futurePeriod := [2]time.Time{
		clock.Current().Add(time.Hour),
		clock.Current().Add(time.Hour*2),
	}

	t.Run("it returns an error if subject does not exist", func(t *testing.T) {
		cmd := application.CreateReservation{1234, users[0].Id, futurePeriod[0], futurePeriod[1]}

		_, err := handler.Handle(cmd)

		assert.Error(t, err)
	})

	t.Run("it returns an error if user does not exist", func(t *testing.T) {
		cmd := application.CreateReservation{subjects[0].Id, 1234, futurePeriod[0], futurePeriod[1]}

		_, err := handler.Handle(cmd)

		assert.Error(t, err)
	})

	t.Run("it returns an error on attempt to create a reservation in the past", func(t *testing.T) {
		current := clock.Current()
		t.Cleanup(func() {
			clock.Set(current)
		})
		cmd := application.CreateReservation{subjects[0].Id, users[0].Id, futurePeriod[0], futurePeriod[1]}
		clock.Set(cmd.From.Add(time.Minute*2))

		_, err := handler.Handle(cmd)

		assert.Error(t, err)
	})

	t.Run("it creates a reservation, if subject is available at given time period", func(t *testing.T) {
		user := users[0]
		subject := subjects[0]
		cmd := application.CreateReservation{subject.Id, user.Id, futurePeriod[0], futurePeriod[1]}

		result, err := handler.Handle(cmd)
		assert.NoError(t,err)
		r, err := registry.Get(result.Id)

		assert.NoError(t,err)
		assert.True(t, result.SubjectId == subject.Id)
		assert.True(t, result.UserId == user.Id)
		assert.True(t, result.Start.Equal(futurePeriod[0]))
		assert.True(t, result.End.Equal(futurePeriod[1]))
		assert.Equal(t, result, r)
		t.Cleanup(func() {
			registry.Remove(result.Id)
		})
	})

	t.Run("it cannot create a reservation given a conflicting reservation exists", func(t *testing.T) {
		user := users[0]
		subject := subjects[0]
		r1 := createReservation(t, subject.Id, users[1].Id, futurePeriod[0].Add(time.Minute), futurePeriod[1].Add(-time.Minute*5))
		r2 := createReservation(t, subject.Id, users[2].Id, futurePeriod[1].Add(-time.Minute*4), futurePeriod[1].Add(time.Minute))
		cmd := application.CreateReservation{subject.Id, user.Id, futurePeriod[0], futurePeriod[1]}

		_, err := handler.Handle(cmd)

		assertAlreadyReservedError(t, err, []int{r1.Id, r2.Id})
	})

	t.Run("it can create a reservation, given the conflicting reservation belongs to another subject", func(t *testing.T) {
		user := users[0]
		subject := subjects[0]
		createReservation(t, subjects[1].Id, users[1].Id, futurePeriod[0].Add(time.Minute), futurePeriod[1].Add(time.Minute))

		cmd := application.CreateReservation{subject.Id, user.Id, futurePeriod[0], futurePeriod[1]}
		reservation, err := handler.Handle(cmd)
		assert.NoError(t, err)
		t.Cleanup(func() {
			registry.Remove(reservation.Id)
		})
	})
}

func createTestSubjects(store reservations.SubjectsStore, t *testing.T) reservations.Subjects {
	subjects := reservations.Subjects{
		reservations.Subject{Id: 1, Name: "Subject#1"},
		reservations.Subject{Id: 2, Name: "Subject#2"},
	}

	for _, s := range subjects {
		err := store.Add(s)
		if err != nil {
			t.Fatal(err)
		}
	}

	return subjects
}

func createTestUsers(store users.UsersStore, t *testing.T) []users.User  {
	users := []users.User{
		{Id: 1, Name: "Test 1"},
		{Id: 2, Name: "Test 2"},
		{Id: 3, Name: "Test 3"},
	}

	for _, u := range users {
		err := store.Add(u)
		if err != nil {
			t.Fatal(err)
		}
	}

	return users
}

func createReservation(t *testing.T, subjectId int, userId int, start time.Time, end time.Time) reservations.Reservation {
	r := reservations.Reservation{
		Id: registry.NextIdentity(),
		UserId: userId,
		SubjectId: subjectId,
		Start: start,
		End: end,
	}
	err := registry.Add(r)
	assert.NoError(t, err)

	t.Cleanup(func() {
		registry.Remove(r.Id)
	})

	return r
}

func assertAlreadyReservedError(t *testing.T, err error, reservationIds []int) {
	assert.Error(t, err)
	reservationError, ok := err.(application.AlreadyReservedError)
	assert.True(t, ok)
	assert.Equal(t, reservationIds, reservationError.ReservationIds)
}
