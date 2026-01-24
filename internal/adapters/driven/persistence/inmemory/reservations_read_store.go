package inmemory

import (
	"slices"
	"time"

	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
	reservationPorts "github.com/SneedusSnake/Reservations/internal/ports/reservations"
	userPorts "github.com/SneedusSnake/Reservations/internal/ports/users"
	readmodel "github.com/SneedusSnake/Reservations/internal/read_model"
)

type ReservationsReadStore struct {
	reservationsStore reservationPorts.ReservationsRepository
	users        userPorts.UsersRepository
	subjects 	reservationPorts.SubjectsRepository
}

func NewReservationReadStore(reservationsStore *ReservationsStore, users *UsersStore, subjects *SubjectsStore) *ReservationsReadStore {
	return &ReservationsReadStore{
		reservationsStore: reservationsStore,
		users: users,
		subjects: subjects,
	}
}

func (r *ReservationsReadStore) Get(id int) (readmodel.Reservation, error) {
	reservation, err := r.reservationsStore.Get(id)

	if err != nil {
		return readmodel.Reservation{}, err
	}
	result, err := r.make(reservation)
	if err != nil {
		return readmodel.Reservation{}, err
	}

	return result, nil
}

func (r *ReservationsReadStore) Active(t time.Time, tags ...string) ([]readmodel.Reservation, error) {
	var result []readmodel.Reservation
	list := r.reservationsStore.List()
	if len(tags) > 0 {
		filterSubjects, err := r.subjects.GetByTags(tags)
		if err != nil {
			return []readmodel.Reservation{}, err
		}
		list = filterBySubjects(list, filterSubjects)
	}

	for _, reservation := range list {
		model, err := r.make(reservation)
		if err != nil {
			return []readmodel.Reservation{}, err
		}

		if (!reservation.End.Before(t)) {
			result = append(result, model)
		}
	}
	return result, nil
}

func (r *ReservationsReadStore) make(reservation reservations.Reservation) (readmodel.Reservation, error) {
	user, err := r.users.Get(reservation.UserId)
	if err != nil {
		return readmodel.Reservation{}, err
	}

	subject, err := r.subjects.Get(reservation.SubjectId)
	if err != nil {
		return readmodel.Reservation{}, err
	}

	return readmodel.Reservation{Id: reservation.Id, Subject: subject.Name, User: user.Name, Start: reservation.Start, End: reservation.End}, nil
}

func filterBySubjects(rs reservations.Reservations, subjects reservations.Subjects) reservations.Reservations {
	var subjectIds []int
	var filtered reservations.Reservations

	for _, subject := range subjects {
		subjectIds = append(subjectIds, subject.Id)
	}

	for _, r := range rs {
		if slices.Contains(subjectIds, r.SubjectId) {
			filtered = append(filtered, r)
		}
	}

	return filtered
}
