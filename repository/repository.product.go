package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rizkimul/gorilla-begin/v2/entity"
)

type RepositoryProduct interface {
	GetProductall() ([]entity.Product, error)
	GetProductById(id string) ([]entity.Product, error)
	InsertProduct(product *entity.Product) (*entity.Product, error)
	UpdateProduct(id string, person *entity.Product) (int64, error)
	DeleteProduct(id string) (int64, error)
}

type repoProduct struct {
	DB *sqlx.DB
}

const (
	getProductAll  = "SELECT * FROM product"
	getProductById = "SELECT * FROM product WHERE id=$1"
	insertProduct  = "INSERT INTO product (product_name, product_description, price, product_image) VALUES ($1, $2, $3, $4)"
	updateProduct  = "UPDATE product SET (product_name, product_description, price, product_image) = ($1, $2, $3, $4) WHERE id=$5"
	deleteProduct  = "DELETE FROM product WHERE id=$1"
)

func NewProductRepository(db *sqlx.DB) RepositoryProduct {
	return &repoProduct{
		DB: db,
	}
}

// var db, err = sqlx.Connect("postgres", "Product=postgres password=root dbname=db_golang sslmode=disable")

func (rp *repoProduct) GetProductall() ([]entity.Product, error) {
	product := []entity.Product{}
	err := rp.DB.Select(&product, getProductAll)

	return product, err
}

func (rp *repoProduct) GetProductById(id string) ([]entity.Product, error) {
	product := []entity.Product{}

	err := rp.DB.Select(&product, getProductById, id)
	return product, err
}

func (rp *repoProduct) InsertProduct(product *entity.Product) (*entity.Product, error) {
	var id string
	err := rp.DB.QueryRow(insertProduct, product.Product_name, product.Product_description, product.Price, product.Product_image).Scan(&id)

	return product, err
}

func (rp *repoProduct) UpdateProduct(id string, product *entity.Product) (int64, error) {
	res, err := rp.DB.Exec(updateProduct, product.Product_name, product.Product_description, product.Price, product.Product_image, id)

	rowsAfffected, err := res.RowsAffected()

	return rowsAfffected, err
}

func (rp *repoProduct) DeleteProduct(id string) (int64, error) {
	res, err := rp.DB.Exec(deleteProduct, id)

	RowsAffected, err := res.RowsAffected()

	return RowsAffected, err
}
