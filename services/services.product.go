package services

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/rizkimul/gorilla-begin/v2/config"
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/repository"
)

type ProductServices interface {
	Getallproduct() ([]entity.Product, error)
	GetproductById(id int) (entity.Product, error)
	Insertproduct(file multipart.File, product *entity.Product) (*entity.Product, error)
	Updateproduct(id int, file multipart.File, product *entity.Product) (*entity.Product, error)
	Deleteproduct(id int) error
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

func (s *prdctsvc) GetproductById(id int) (entity.Product, error) {
	return s.prdctrepo.GetProductById(id)
}

func (s *prdctsvc) Insertproduct(file multipart.File, product *entity.Product) (*entity.Product, error) {
	conf, _ := config.LoadConfig(".")
	cld, _ := cloudinary.NewFromURL(conf.CloudSecretKey)
	ctx := context.Background()
	resp, _ := cld.Upload.Upload(ctx, file, uploader.UploadParams{})

	product.ProductImage = resp.SecureURL
	product.CreatedAt = time.Now()

	err := s.prdctrepo.InsertProduct(product)

	return product, err
}

func (s *prdctsvc) Updateproduct(id int, file multipart.File, product *entity.Product) (*entity.Product, error) {
	conf, _ := config.LoadConfig(".")
	cld, _ := cloudinary.NewFromURL(conf.CloudSecretKey)
	ctx := context.Background()
	resp, _ := cld.Upload.Upload(ctx, file, uploader.UploadParams{})

	product.ProductImage = resp.SecureURL
	product.UpdatedAt = time.Now()

	err := s.prdctrepo.UpdateProduct(id, product)
	return product, err
}

func (s *prdctsvc) Deleteproduct(id int) error {
	return s.prdctrepo.DeleteProduct(id)
}
