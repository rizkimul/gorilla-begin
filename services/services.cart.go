package services

import (
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/repository"
)

type CartServices interface {
	Getall() ([]entity.Cart, error)
	GetById(id int) (entity.Cart, error)
	Insert(cart *entity.Cart) error
	Update(id int, cart *entity.Cart) error
	Delete(id int) error
}

type cartsvc struct {
	cartrepo   repository.CartRepository
	spCartRepo repository.SPCartRepository
}

// var repo repository.Repository = repository.NewRepository()

func NewCartServices(cartrepo repository.CartRepository, spCartRepo repository.SPCartRepository) CartServices {
	return &cartsvc{
		cartrepo:   cartrepo,
		spCartRepo: spCartRepo,
	}
}

func (s *cartsvc) Getall() ([]entity.Cart, error) {
	return s.cartrepo.Getcartall()
}

func (s *cartsvc) GetById(id int) (entity.Cart, error) {
	return s.cartrepo.GetcartById(id)
}

func (s *cartsvc) Insert(cart *entity.Cart) error {
	return s.cartrepo.Insertcart(cart)
}

func (s *cartsvc) Update(id int, person *entity.Cart) error {
	return s.cartrepo.Updatecart(id, person)
}

func (s *cartsvc) Delete(id int) error {
	return s.cartrepo.Deletecart(id)
}
