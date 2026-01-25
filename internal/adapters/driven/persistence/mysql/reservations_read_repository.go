package mysql

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"

	readmodel "github.com/SneedusSnake/Reservations/internal/read_model"
)

type ReservationsReadRepository struct {
	connection *sql.DB
}

func NewReservationsReadRepository(connection *sql.DB) *ReservationsReadRepository {
	return &ReservationsReadRepository{connection: connection}
}

func (r *ReservationsReadRepository) Get(id int) (readmodel.Reservation, error) {
	var result readmodel.Reservation
	
	row := r.connection.QueryRow(baseQuery() + " WHERE r.id = ?", id)

	if err := row.Scan(
		&result.Id,
		&result.Subject,
		&result.User,
		&result.Start,
		&result.End,
	); err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("Reservation with id %d was not found", id)
		}
		return result, err
	}
	
	return result, nil
}

func (r *ReservationsReadRepository) Active(t time.Time, tags ...string) ([]readmodel.Reservation, error) {
	var result []readmodel.Reservation
	var params []any
	params = append(params, t, t)
	query := baseQuery()
	conditions := ` WHERE 1=1 AND r.start <= ? AND r.end >= ?`

	if len(tags) > 0 {
		slices.Sort(tags)
		query += `
			JOIN (SELECT subject_id, GROUP_CONCAT(tag ORDER BY tag) AS tags 
			FROM subject_tags GROUP BY subject_id) t ON t.subject_id = s.id
		`
		params = append(params, strings.Join(tags, ","))
		conditions += ` AND t.tags LIKE ?`
	}
	rows, err := r.connection.Query(
		query + conditions,
		params...,
	)

	if err != nil {
		return result, err
	}

	for rows.Next() {
		var model readmodel.Reservation
		err = rows.Scan(
			&model.Id,
			&model.Subject,
			&model.User,
			&model.Start,
			&model.End,
		)
		if err != nil {
			return result, err
		}
		result = append(result, model)
	}

	return result, nil
}

func baseQuery() string {
	return `
		SELECT r.id, s.name, u.name, r.start, r.end FROM reservations r
		JOIN users u on u.id = r.user_id
		JOIN subjects s on s.id = r.subject_id
	`
}
