package mysql

import (
	"database/sql"
	"fmt"
)

type sequence struct {
	name string
	connection *sql.DB
}

func (seq *sequence) Next() (int, error) {
	var id int
	_, err := seq.connection.Exec(fmt.Sprintf("UPDATE %s SET value = LAST_INSERT_ID(value+1)", seq.name))
	if err != nil {
		return 0, nil
	}

	row := seq.connection.QueryRow("SELECT LAST_INSERT_ID() as id")
	if err = row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
