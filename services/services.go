package services

import (
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/repository"
)

type Services interface {
	Getall() ([]entity.Person, error)
	GetById(id string) ([]entity.Person, error)
	Insert(person *entity.Person) (*entity.Person, error)
	Update(id string, person *entity.Person) (int64, error)
	Delete(id string) (int64, error)
}

type svc struct {
	repo repository.Repository
}

// var repo repository.Repository = repository.NewRepository()

func NewServices(repo repository.Repository) Services {
	return &svc{
		repo: repo,
	}
}

func (s *svc) Getall() ([]entity.Person, error) {
	return s.repo.Getuserall()
}

func (s *svc) GetById(id string) ([]entity.Person, error) {
	return s.repo.GetUserById(id)
}

func (s *svc) Insert(person *entity.Person) (*entity.Person, error) {
	return s.repo.InsertUser(person)
}

func (s *svc) Update(id string, person *entity.Person) (int64, error) {
	return s.repo.UpdateUser(id, person)
}

func (s *svc) Delete(id string) (int64, error) {
	return s.repo.DeleteUser(id)
}
