package services

import (
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/repository"
)

type CartServices interface {
	Getall() ([]entity.Cart, error)
	GetById(id string) ([]entity.Cart, error)
	Insert(cart *entity.Cart) (*entity.Cart, error)
	Update(id string, cart *entity.Cart) (int64, error)
	Delete(id string) (int64, error)
}

type cartsvc struct {
	cartrepo repository.CartRepository
}

// var repo repository.Repository = repository.NewRepository()

func NewCartServices(cartrepo repository.CartRepository) CartServices {
	return &cartsvc{
		cartrepo: cartrepo,
	}
}

func (s *cartsvc) Getall() ([]entity.Cart, error) {
	return s.cartrepo.Getcartall()
}

func (s *cartsvc) GetById(id string) ([]entity.Cart, error) {
	return s.cartrepo.GetcartById(id)
}

func (s *cartsvc) Insert(person *entity.Cart) (*entity.Cart, error) {
	return s.cartrepo.Insertcart(person)
}

func (s *cartsvc) Update(id string, person *entity.Cart) (int64, error) {
	return s.cartrepo.Updatecart(id, person)
}

func (s *cartsvc) Delete(id string) (int64, error) {
	return s.cartrepo.Deletecart(id)
}
