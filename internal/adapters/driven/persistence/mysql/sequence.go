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
	result, err := seq.connection.Exec(fmt.Sprintf("UPDATE %s SET value = LAST_INSERT_ID(value+1)", seq.name))
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()

	return int(id), err
}
