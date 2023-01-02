package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rizkimul/gorilla-begin/v2/entity"
)

type Repository interface {
	Getuserall() ([]entity.Person, error)
	GetUserById(id string) ([]entity.Person, error)
	InsertUser(person *entity.Person) (*entity.Person, error)
	UpdateUser(id string, person *entity.Person) (int64, error)
	DeleteUser(id string) (int64, error)
}

type repo struct {
	DB *sqlx.DB
}

const (
	getUserAll  = "SELECT * FROM person"
	getUserById = "SELECT * FROM person WHERE id=$1"
	insertUser  = "INSERT INTO person (name, email, password, phonenumber) VALUES ($1, $2, $3, $4)"
	updateUser  = "UPDATE person SET (name, email, password, phonenumber) = ($1, $2, $3, $4) WHERE id=$5"
	deleteUser  = "DELETE FROM person WHERE id=$1"
)

func NewRepository(db *sqlx.DB) Repository {
	return &repo{
		DB: db,
	}
}

func (r *repo) Getuserall() ([]entity.Person, error) {
	person := []entity.Person{}
	err := r.DB.Select(&person, getUserAll)

	return person, err
}

func (r *repo) GetUserById(id string) ([]entity.Person, error) {
	person := []entity.Person{}

	err := r.DB.Select(&person, getUserById, id)
	return person, err
}

func (r *repo) InsertUser(person *entity.Person) (*entity.Person, error) {
	var id string
	err := r.DB.QueryRow(insertUser, person.Name, person.Email, person.Password, person.Phonenumber).Scan(&id)

	return person, err
}

func (r *repo) UpdateUser(id string, person *entity.Person) (int64, error) {
	res, err := r.DB.Exec(updateUser, person.Name, person.Email, person.Password, person.Phonenumber, id)

	rowsAfffected, err := res.RowsAffected()

	return rowsAfffected, err
}

func (r *repo) DeleteUser(id string) (int64, error) {
	res, err := r.DB.Exec(deleteUser, id)

	RowsAffected, err := res.RowsAffected()

	return RowsAffected, err
}
