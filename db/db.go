package database

import "github.com/jmoiron/sqlx"

type Connect interface {
	Connection() (*sqlx.DB, error)
}

type connect struct{}

func NewConnection() Connect {
	return &connect{}
}

func (*connect) Connection() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=postgres password=root dbname=db_golang sslmode=disable")
	return db, err
}
