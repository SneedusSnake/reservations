package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/SneedusSnake/Reservations/internal/domain/reservations"
)

type ReservationsRepository struct {
	connection *sql.DB
	sequence *sequence
}

func NewReservationsRepository(connection *sql.DB) *ReservationsRepository {
	return &ReservationsRepository{
		connection: connection,
		sequence: &sequence{
			name: "reservation_seq",
			connection: connection,
		},
	}
}

func (r *ReservationsRepository) NextIdentity() (int, error) {
	return r.sequence.Next()
}

func (r *ReservationsRepository) List() reservations.Reservations {
	var result reservations.Reservations 

	rows, err := r.connection.Query("SELECT*FROM reservations")

	if err != nil {
		return reservations.Reservations{}
	}

	for rows.Next() {
		var record reservations.Reservation
		if err = rows.Scan(
			&record.Id,
			&record.UserId,
			&record.SubjectId,
			&record.Start,
			&record.End,
		); err != nil {
			return reservations.Reservations{}
		}
		result = append(result, record)
	}

	return result
}

func (r *ReservationsRepository) Add(record reservations.Reservation) error {
	_, err := r.connection.Exec("INSERT INTO reservations(id, user_id, subject_id, start, end) VALUES(?,?,?,?,?)", record.Id, record.UserId, record.SubjectId, record.Start, record.End)

	return err
}

func (r *ReservationsRepository) Get(id int) (reservations.Reservation, error) {
	var result reservations.Reservation
	var err error

	row := r.connection.QueryRow("SELECT*FROM reservations WHERE id = ?", id)

	if err = row.Scan(&result.Id, &result.UserId, &result.SubjectId, &result.Start, &result.End); err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("Reservation with id %d was not found", id)
		}
	}

	return result, err
}

func (r *ReservationsRepository) Remove(id int) error {
	_, err := r.connection.Exec("DELETE FROM reservations WHERE id = ?", id)

	return err
}

func (r *ReservationsRepository) ForPeriod(from time.Time, to time.Time) (reservations.Reservations, error) {
	var result reservations.Reservations 

	rows, err := r.connection.Query(
		"SELECT*FROM reservations where (start BETWEEN ? AND ?) OR (end BETWEEN ? AND ?)",
		from,
		to,
		from,
		to,
	)

	if err != nil {
		return result, err
	}

	for rows.Next() {
		var record reservations.Reservation
		if err = rows.Scan(
			&record.Id,
			&record.UserId,
			&record.SubjectId,
			&record.Start,
			&record.End,
		); err != nil {
			return result, err
		}
		result = append(result, record)
	}

	return result, nil
}
