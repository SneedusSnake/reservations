package application_test

import (
	"testing"
	"time"

	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
	"github.com/SneedusSnake/Reservations/internal/application"
	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	"github.com/SneedusSnake/Reservations/internal/domain/users"
	reservationsPort "github.com/SneedusSnake/Reservations/internal/ports/reservations"
	usersPort "github.com/SneedusSnake/Reservations/internal/ports/users"
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

func (c *FakeClock) TimeTravel(minutes int) time.Time {
	return c.Current().Add(time.Minute*time.Duration(minutes))
}

var subjectsStore reservationsPort.SubjectsRepository
var reservationsStore reservationsPort.ReservationsRepository
var usersStore usersPort.UsersRepository
var clock *FakeClock

func TestCreateReservation(t *testing.T) {
	handler := getSUT()
	subjects := createTestSubjects(subjectsStore, t)
	users := createTestUsers(usersStore, t)
	futurePeriod := [2]time.Time{
		clock.Current().Add(time.Hour),
		clock.Current().Add(time.Hour * 2),
	}

	t.Run("it returns an error if subject does not exist", func(t *testing.T) {
		cmd := application.CreateReservation{1234, users[0].Id, futurePeriod[0], futurePeriod[1]}

		_, err := handler.Create(cmd)

		assert.Error(t, err)
	})

	t.Run("it returns an error if user does not exist", func(t *testing.T) {
		cmd := application.CreateReservation{subjects[0].Id, 1234, futurePeriod[0], futurePeriod[1]}

		_, err := handler.Create(cmd)

		assert.Error(t, err)
	})

	t.Run("it returns an error on attempt to create a reservation in the past", func(t *testing.T) {
		current := clock.Current()
		t.Cleanup(func() {
			clock.Set(current)
		})
		cmd := application.CreateReservation{subjects[0].Id, users[0].Id, futurePeriod[0], futurePeriod[1]}
		clock.Set(cmd.From.Add(time.Minute * 2))

		_, err := handler.Create(cmd)

		assert.Error(t, err)
	})

	t.Run("it creates a reservation, if subject is available at given time period", func(t *testing.T) {
		user := users[0]
		subject := subjects[0]
		cmd := application.CreateReservation{subject.Id, user.Id, futurePeriod[0], futurePeriod[1]}

		result, err := handler.Create(cmd)
		assert.NoError(t, err)
		r, err := reservationsStore.Get(result.Id)

		assert.NoError(t, err)
		assert.True(t, result.SubjectId == subject.Id)
		assert.True(t, result.UserId == user.Id)
		assert.True(t, result.Start.Equal(futurePeriod[0]))
		assert.True(t, result.End.Equal(futurePeriod[1]))
		assert.Equal(t, result, r)
		t.Cleanup(func() {
			reservationsStore.Remove(result.Id)
		})
	})

	t.Run("it cannot create a reservation given a conflicting reservation exists", func(t *testing.T) {
		user := users[0]
		subject := subjects[0]
		r1 := createReservation(t, subject.Id, users[1].Id, futurePeriod[0].Add(time.Minute), futurePeriod[1].Add(-time.Minute*5))
		r2 := createReservation(t, subject.Id, users[2].Id, futurePeriod[1].Add(-time.Minute*4), futurePeriod[1].Add(time.Minute))
		cmd := application.CreateReservation{subject.Id, user.Id, futurePeriod[0], futurePeriod[1]}

		_, err := handler.Create(cmd)

		assertAlreadyReservedError(t, err, []int{r1.Id, r2.Id})
	})

	t.Run("it can create a reservation, given the conflicting reservation belongs to another subject", func(t *testing.T) {
		user := users[0]
		subject := subjects[0]
		createReservation(t, subjects[1].Id, users[1].Id, futurePeriod[0].Add(time.Minute), futurePeriod[1].Add(time.Minute))

		cmd := application.CreateReservation{subject.Id, user.Id, futurePeriod[0], futurePeriod[1]}
		reservation, err := handler.Create(cmd)
		assert.NoError(t, err)
		t.Cleanup(func() {
			reservationsStore.Remove(reservation.Id)
		})
	})
}

func TestRemoveReservation(t *testing.T) {
	handler := getSUT()
	subjects := createTestSubjects(subjectsStore, t)
	users := createTestUsers(usersStore, t)

	t.Run("it returns error if no reservation for subject exists", func(t *testing.T) {
		cmd := application.RemoveReservations{users[0].Id, subjects[0].Id}

		err := handler.Remove(cmd)

		assert.Error(t, err)
	})

	t.Run("it removes all user reservations for subject", func(t *testing.T) {
		createReservation(t, subjects[0].Id, users[0].Id, clock.TimeTravel(-10), clock.TimeTravel(5))
		createReservation(t, subjects[0].Id, users[0].Id, clock.TimeTravel(5), clock.TimeTravel(10))
		cmd := application.RemoveReservations{users[0].Id, subjects[0].Id}

		err := handler.Remove(cmd)

		assert.NoError(t, err)
		assert.Equal(t, 0, len(reservationsStore.List()))
	})

	t.Run("it does not remove user's past reservations", func(t *testing.T) {
		createReservation(t, subjects[0].Id, users[0].Id, clock.TimeTravel(-60), clock.TimeTravel(-30))
		createReservation(t, subjects[0].Id, users[0].Id, clock.TimeTravel(-10), clock.TimeTravel(-5))
		cmd := application.RemoveReservations{users[0].Id, subjects[0].Id}

		err := handler.Remove(cmd)

		assert.Error(t, err)
		assert.Equal(t, 2, len(reservationsStore.List()))
	})

	t.Run("it does not remove other users' reservations", func(t *testing.T) {
		createReservation(t, subjects[0].Id, users[1].Id, clock.TimeTravel(5), clock.TimeTravel(10))
		cmd := application.RemoveReservations{users[0].Id, subjects[0].Id}

		err := handler.Remove(cmd)

		assert.Error(t, err)
		assert.Equal(t, 1, len(reservationsStore.List()))
	})

	t.Run("it does not remove other user's subjects reservations", func(t *testing.T) {
		createReservation(t, subjects[1].Id, users[0].Id, clock.TimeTravel(5), clock.TimeTravel(10))
		cmd := application.RemoveReservations{users[0].Id, subjects[0].Id}

		err := handler.Remove(cmd)

		assert.Error(t, err)
		assert.Equal(t, 1, len(reservationsStore.List()))
	})
}

func getSUT() *application.ReservationService {
	subjectsStore = inmemory.NewSubjectsStore()
	usersStore = inmemory.NewUsersStore()
	reservationsStore = inmemory.NewReservationStore()
	clock = &FakeClock{}
	clock.Set(time.Now())
	return application.NewReservationService(
		subjectsStore,
		reservationsStore,
		usersStore,
		clock,
	)
}

func createTestSubjects(store reservationsPort.SubjectsRepository, t *testing.T) reservations.Subjects {
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

func createTestUsers(store users.UsersStore, t *testing.T) []users.User {
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
		Id:        reservationsStore.NextIdentity(),
		UserId:    userId,
		SubjectId: subjectId,
		Start:     start,
		End:       end,
	}
	err := reservationsStore.Add(r)
	assert.NoError(t, err)

	t.Cleanup(func() {
		reservationsStore.Remove(r.Id)
	})

	return r
}

func assertAlreadyReservedError(t *testing.T, err error, reservationIds []int) {
	assert.Error(t, err)
	reservationError, ok := err.(application.AlreadyReservedError)
	assert.True(t, ok)
	assert.Equal(t, reservationIds, reservationError.ReservationIds)
}
