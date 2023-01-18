package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rizkimul/gorilla-begin/v2/entity"
)

type Repository interface {
	Getuserall() ([]entity.Person, error)
	GetUserById(id string) (entity.Person, error)
	InsertUser(person *entity.Person) error
	UpdateUser(id string, person *entity.Person) error
	DeleteUser(id string) error
	Login(name string) (person *entity.Person, err error)
}

type repo struct {
	DB *sqlx.DB
}

const (
	getUserAll     = "SELECT * FROM person"
	getUserById    = "SELECT * FROM person WHERE id=$1"
	insertUser     = "INSERT INTO person (name, email, password, created_at) VALUES ($1, $2, $3, $4)"
	updateUser     = "UPDATE person SET (name, email, password, updated_at) = ($1, $2, $3, $4) WHERE id=$5"
	deleteUser     = "DELETE FROM person WHERE id=$1"
	getUserbyEmail = "SELECT * FROM person WHERE email=$1"
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

func (r *repo) GetUserById(id string) (entity.Person, error) {
	person := entity.Person{}

	err := r.DB.Get(&person, getUserById, id)
	return person, err
}

func (r *repo) InsertUser(person *entity.Person) error {
	_, err := r.DB.Exec(insertUser, person.Name, person.Email, person.Password, person.CreatedAt)

	return err
}

func (r *repo) UpdateUser(id string, person *entity.Person) error {
	_, err := r.DB.Exec(updateUser, person.Name, person.Email, person.Password, person.UpdatedAt, id)

	return err
}

func (r *repo) DeleteUser(id string) error {
	_, err := r.DB.Exec(deleteUser, id)

	return err
}

func (r *repo) Login(name string) (person *entity.Person, err error) {
	person = new(entity.Person)
	err = r.DB.Get(person, getUserbyEmail, name)
	return person, err
}
