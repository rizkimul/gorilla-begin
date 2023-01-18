package services

import (
	"time"

	"github.com/rizkimul/gorilla-begin/v2/config"
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/helper"
	"github.com/rizkimul/gorilla-begin/v2/repository"
	"github.com/rizkimul/gorilla-begin/v2/response"
	"github.com/rizkimul/gorilla-begin/v2/utils"
)

type Services interface {
	Getall() ([]entity.Person, error)
	GetById(id string) (entity.Person, error)
	Insert(person *entity.Person) error
	Update(id string, person *entity.Person) error
	Delete(id string) error
	Login(email string, pass string) (map[string]string, error)
	Print(result []response.Response) ([]byte, error)
}

type svc struct {
	repo   repository.Repository
	helper helper.Helper
	token  utils.Token
}

// var repo repository.Repository = repository.NewRepository()

func NewServices(repo repository.Repository, helper helper.Helper, token utils.Token) Services {
	return &svc{
		repo:   repo,
		helper: helper,
		token:  token,
	}
}

func (s *svc) Getall() ([]entity.Person, error) {
	return s.repo.Getuserall()
}

func (s *svc) GetById(id string) (entity.Person, error) {
	return s.repo.GetUserById(id)
}

func (s *svc) Insert(person *entity.Person) error {
	passhash, _ := s.helper.HashPass(person.Password)
	person.Password = passhash
	person.CreatedAt = time.Now()
	err := s.repo.InsertUser(person)

	return err
}

func (s *svc) Update(id string, person *entity.Person) error {
	person.UpdatedAt = time.Now()
	passhash, _ := s.helper.HashPass(person.Password)
	person.Password = passhash
	err := s.repo.UpdateUser(id, person)
	return err
}

func (s *svc) Delete(id string) error {
	return s.repo.DeleteUser(id)
}

func (s *svc) Login(email string, pass string) (map[string]string, error) {
	loginRes, err := s.repo.Login(email)
	if err != nil {
		return nil, err
	}

	if err := s.helper.MatchPass(pass, loginRes.Password); err != nil {
		return nil, err
	}

	conf, err := config.LoadConfig(".")
	if err != nil {
		return nil, err
	}
	accessToken, err := s.token.CreateToken(conf.AccessTokenExp)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.token.CreateToken(conf.RefreshTokenExp)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Access Token": accessToken, "Refresh Token": refreshToken}, nil

}

func (s *svc) Print(result []response.Response) ([]byte, error) {
	pdfBytes, err := s.helper.PrintData(&result)
	return pdfBytes, err
}
