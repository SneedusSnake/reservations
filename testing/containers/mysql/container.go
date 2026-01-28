package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/SneedusSnake/Reservations/testing/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type MysqlContainer struct {
	*mysql.MySQLContainer
	connectionString string
}

func (container *MysqlContainer) migrate() error {
	conn, err := container.Connection()
    if err != nil {
		return err
    }

	goose.SetDialect("mysql")
	return goose.Up(conn, utils.TestsRootDir() + "/../migrations")
}

func (container *MysqlContainer) Connection() (*sql.DB, error) {
	db, err := sql.Open("mysql", container.connectionString)
    if err != nil {
		return nil, err
    }

    pingErr := db.Ping()
    if pingErr != nil {
		return nil, err
    }

	return db, nil
}

func (container *MysqlContainer) ExternalConnectionString(ctx context.Context, args ...string) (string, error) {
	endpoint, err := container.ContainerIP(ctx)
	if err != nil {
		return "", err
	}

	extraArgs := ""
	args = append(args, "parseTime=true")
	if len(args) > 0 {
		extraArgs = strings.Join(args, "&")
	}
	if extraArgs != "" {
		extraArgs = "?" + extraArgs
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s%s", "root", "root", endpoint, "app", extraArgs), nil
}

func Start(ctx context.Context, network string, logs ...testcontainers.LogConsumer) (*MysqlContainer, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Networks: []string{network},
		},
	}
	container, err := mysql.Run(
		ctx,
		"mysql:8.4",
		mysql.WithDatabase("app"),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
		testcontainers.CustomizeRequest(req),
	)

	if err != nil {
		return nil, err
	}

	connStr, err := container.ConnectionString(ctx, "parseTime=true")
	if err != nil {
		return nil, err
	}
	sqlContainer := &MysqlContainer{
		MySQLContainer: container,
		connectionString: connStr,
	}
	err = sqlContainer.migrate()
	if err != nil {
		return nil, err
	}

	return sqlContainer, nil
}


