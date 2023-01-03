package services

import (
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/repository"
	"github.com/rizkimul/gorilla-begin/v2/response"
)

type SPCartServices interface {
	Getall() ([]response.ShoppingCart, error)
	GetById(id string) ([]entity.ShoppingCart, error)
	Insert(spcart *entity.ShoppingCart) (*entity.ShoppingCart, error)
	Update(id string, spcart *entity.ShoppingCart) (int64, error)
	Delete(id string) (int64, error)
}

type spcartsvc struct {
	spcartrepo repository.SPCartRepository
}

// var repo repository.Repository = repository.NewRepository()

func NewSPCartServices(spcartrepo repository.SPCartRepository) SPCartServices {
	return &spcartsvc{
		spcartrepo: spcartrepo,
	}
}

func (s *spcartsvc) Getall() ([]response.ShoppingCart, error) {
	return s.spcartrepo.Getspcartall()
}

func (s *spcartsvc) GetById(id string) ([]entity.ShoppingCart, error) {
	return s.spcartrepo.GetspcartById(id)
}

func (s *spcartsvc) Insert(spcart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	return s.spcartrepo.Insertspcart(spcart)
}

func (s *spcartsvc) Update(id string, spcart *entity.ShoppingCart) (int64, error) {
	return s.spcartrepo.Updatespcart(id, spcart)
}

func (s *spcartsvc) Delete(id string) (int64, error) {
	return s.spcartrepo.Deletespcart(id)
}
