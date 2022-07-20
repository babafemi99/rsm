package psql

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type psql struct {
	log  *logrus.Logger
	conn *pgx.Conn
}

func NewPsqlStore(log *logrus.Logger) (*psql, error) {
	log.Info("Database connection starting")
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:mysecretpassword@localhost:5432/rsmstore")
	if err != nil {
		log.Fatalf("error connecting to db %v", err)
	}
	query, queryErr := ioutil.ReadFile("./../../migration/create/create_table.sql")
	if queryErr != nil {
		log.Fatalf("Error reading migratison file %v", queryErr)
	}
	_, err = conn.Exec(context.Background(), string(query))
	if err != nil {
		log.Fatalf("Error Executing migration code %v", err)
	}
	log.Info("Database connected successfully")
	return &psql{
		conn: conn,
	}, nil
}

func (p *psql) GetConnection() *pgx.Conn {
	return p.conn
}
