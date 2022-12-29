package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rizkimul/gorilla-begin/v2/entity"
)

type Repository interface {
	Getuserall() ([]entity.Person, error)
	GetUserById(id string) ([]entity.Person, error)
	InsertUser(person *entity.Person) (*entity.Person, error)
	UpdateUser(id string, person *entity.Person) (int64, error)
	DeleteUser(id string) (int64, error)
}

type repo struct{}

const (
	getUserAll  = "SELECT * FROM person"
	getUserById = "SELECT * FROM person WHERE id=$1"
	insertUser  = "INSERT INTO person (name, email, password, phonenumber) VALUES ($1, $2, $3, $4)"
	updateUser  = "UPDATE person SET (name, email, password, phonenumber) = ($1, $2, $3, $4) WHERE id=$5"
	deleteUser  = "DELETE FROM person WHERE id=$1"
)

func NewRepository() Repository {
	return &repo{}
}

var db, err = sqlx.Connect("postgres", "user=postgres password=root dbname=db_golang sslmode=disable")

func (*repo) Getuserall() ([]entity.Person, error) {
	defer db.Close()
	person := []entity.Person{}
	err = db.Select(&person, getUserAll)

	return person, err
}

func (*repo) GetUserById(id string) ([]entity.Person, error) {
	defer db.Close()

	person := []entity.Person{}

	err = db.Select(&person, getUserById, id)
	return person, err
}

func (*repo) InsertUser(person *entity.Person) (*entity.Person, error) {
	defer db.Close()
	var id string
	err = db.QueryRow(insertUser, person.Name, person.Email, person.Password, person.Phonenumber).Scan(&id)

	return person, err
}

func (*repo) UpdateUser(id string, person *entity.Person) (int64, error) {
	defer db.Close()
	res, err := db.Exec(updateUser, person.Name, person.Email, person.Password, person.Phonenumber, id)

	rowsAfffected, err := res.RowsAffected()

	return rowsAfffected, err
}

func (*repo) DeleteUser(id string) (int64, error) {
	defer db.Close()
	res, err := db.Exec(deleteUser, id)

	RowsAffected, err := res.RowsAffected()

	return RowsAffected, err
}
