package services

import (
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/repository"
)

type ProductServices interface {
	Getallproduct() ([]entity.Product, error)
	GetproductById(id string) ([]entity.Product, error)
	Insertproduct(product *entity.Product) (*entity.Product, error)
	Updateproduct(id string, product *entity.Product) (int64, error)
	Deleteproduct(id string) (int64, error)
}

type prdctsvc struct {
	prdctrepo repository.RepositoryProduct
}

// var repo repository.Repository = repository.NewRepository()

func NewProductServices(prdctrepo repository.RepositoryProduct) ProductServices {
	return &prdctsvc{
		prdctrepo: prdctrepo,
	}
}

func (s *prdctsvc) Getallproduct() ([]entity.Product, error) {
	return s.prdctrepo.GetProductall()
}

func (s *prdctsvc) GetproductById(id string) ([]entity.Product, error) {
	return s.prdctrepo.GetProductById(id)
}

func (s *prdctsvc) Insertproduct(product *entity.Product) (*entity.Product, error) {
	return s.prdctrepo.InsertProduct(product)
}

func (s *prdctsvc) Updateproduct(id string, product *entity.Product) (int64, error) {
	return s.prdctrepo.UpdateProduct(id, product)
}

func (s *prdctsvc) Deleteproduct(id string) (int64, error) {
	return s.prdctrepo.DeleteProduct(id)
}
