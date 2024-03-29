package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rizkimul/gorilla-begin/v2/entity"
)

type RepositoryProduct interface {
	GetProductall() ([]entity.Product, error)
	GetProductById(id int) (entity.Product, error)
	InsertProduct(product *entity.Product) error
	UpdateProduct(id int, person *entity.Product) error
	DeleteProduct(id int) error
}

type repoProduct struct {
	DB *sqlx.DB
}

const (
	getProductAll  = "SELECT * FROM product"
	getProductById = "SELECT * FROM product WHERE id=$1"
	insertProduct  = "INSERT INTO product (product_name, product_description, price, product_image) VALUES ($1, $2, $3, $4)"
	updateProduct  = "UPDATE product SET (product_name, product_description, price, product_image, updated_at) = ($1, $2, $3, $4, $5) WHERE id=$6"
	deleteProduct  = "DELETE FROM product WHERE id=$1"
)

func NewProductRepository(db *sqlx.DB) RepositoryProduct {
	return &repoProduct{
		DB: db,
	}
}

func (rp *repoProduct) GetProductall() ([]entity.Product, error) {
	product := []entity.Product{}
	err := rp.DB.Select(&product, getProductAll)

	return product, err
}

func (rp *repoProduct) GetProductById(id int) (entity.Product, error) {
	product := entity.Product{}

	err := rp.DB.Get(&product, getProductById, id)
	return product, err
}

func (rp *repoProduct) InsertProduct(product *entity.Product) error {
	_, err := rp.DB.Exec(insertProduct, product.ProductName, product.ProductDescription, product.Price, product.ProductImage)

	return err
}

func (rp *repoProduct) UpdateProduct(id int, product *entity.Product) error {
	_, err := rp.DB.Exec(updateProduct, product.ProductName, product.ProductDescription, product.Price, product.ProductImage, product.UpdatedAt, id)

	return err
}

func (rp *repoProduct) DeleteProduct(id int) error {
	_, err := rp.DB.Exec(deleteProduct, id)

	return err
}
