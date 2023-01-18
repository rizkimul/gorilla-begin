package services

import (
	"errors"
	"reflect"

	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/helper"
	"github.com/rizkimul/gorilla-begin/v2/repository"
)

type SPCartServices interface {
	Getall() ([]entity.SpCart, error)
	GetById(id int) (entity.ShoppingCart, error)
	Insert(spcart *entity.ShoppingCart) error
	Update(id int, spcart *entity.ShoppingCart) error
	Delete(id int) error
}

type spcartsvc struct {
	spcartrepo repository.SPCartRepository
	helper     helper.Helper
	prodsrvc   ProductServices
	cartsrvc   CartServices
}

func NewSPCartServices(spcartrepo repository.SPCartRepository, helper helper.Helper, prodsrvc ProductServices, cartsrvc CartServices) SPCartServices {
	return &spcartsvc{
		spcartrepo: spcartrepo,
		helper:     helper,
		prodsrvc:   prodsrvc,
		cartsrvc:   cartsrvc,
	}
}

func (s *spcartsvc) Getall() ([]entity.SpCart, error) {
	return s.spcartrepo.Getspcartall()
}

func (s *spcartsvc) GetById(id int) (entity.ShoppingCart, error) {
	return s.spcartrepo.GetspcartById(id)
}

func (s *spcartsvc) Insert(spcart *entity.ShoppingCart) error {
	p, err := s.prodsrvc.GetproductById(spcart.ProductId)
	c, _ := s.cartsrvc.GetById(spcart.CartId)
	v := reflect.ValueOf(c)
	if err != nil {
		err := errors.New("can't find product")
		return err
	} else if v.IsZero() {
		err := errors.New("can't find cart")
		return err
	}
	spcart.TotalPrice = s.helper.CountTotal(spcart.QtyProduct, int(p.Price))
	err = s.spcartrepo.Insertspcart(spcart)
	return err
}

func (s *spcartsvc) Update(id int, spcart *entity.ShoppingCart) error {
	p, err := s.prodsrvc.GetproductById(spcart.ProductId)
	c, _ := s.cartsrvc.GetById(spcart.CartId)
	v := reflect.ValueOf(c)
	if err != nil {
		err := errors.New("can't find product")
		return err
	} else if v.IsZero() {
		err := errors.New("can't find cart")
		return err
	}
	spcart.TotalPrice = s.helper.CountTotal(spcart.QtyProduct, int(p.Price))
	err = s.spcartrepo.Updatespcart(id, spcart)
	return err
}

func (s *spcartsvc) Delete(id int) error {
	return s.spcartrepo.Deletespcart(id)
}
