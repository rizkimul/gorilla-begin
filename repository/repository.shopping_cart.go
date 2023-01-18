package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rizkimul/gorilla-begin/v2/entity"
)

type SPCartRepository interface {
	Getspcartall() ([]entity.SpCart, error)
	GetspcartById(id int) (entity.ShoppingCart, error)
	Insertspcart(spcart *entity.ShoppingCart) error
	Updatespcart(id int, spcart *entity.ShoppingCart) error
	Deletespcart(id int) error
}

type spcartrepo struct {
	DB   *sqlx.DB
	prod RepositoryProduct
}

const (
	// getspcartAll       = "SELECT a.cart_name, b.product_name, b.price, b.product_image, c.qty_product, c.total_price FROM cart a, product b, shopping_cart c where c.cart_id = a.id and c.product_id = b.id"
	getspcartAll       = "SELECT * FROM shopping_cart"
	getspcartById      = "SELECT * FROM shopping_cart WHERE id=$1"
	insertspcart       = "INSERT INTO shopping_cart (cart_id, product_id, qty_product, total_price) VALUES ($1, $2, $3, $4)"
	updatespcart       = "UPDATE shopping_cart SET (cart_id, product_id, qty_product, total_price) = ($1, $2, $3, $4) WHERE id=$5"
	deletespcart       = "DELETE FROM shopping_cart WHERE id=$1"
	getProductToSpCart = "SELECT id, product_name, product_description, price, product_image, created_at FROM product WHERE id=$1"
)

func NewSPCartRepository(db *sqlx.DB, prod RepositoryProduct) SPCartRepository {
	return &spcartrepo{
		DB:   db,
		prod: prod,
	}
}

func (r *spcartrepo) Getspcartall() ([]entity.SpCart, error) {
	spcart := []entity.SpCart{}
	err := r.DB.Select(&spcart, getspcartAll)

	return spcart, err
}

func (r *spcartrepo) GetspcartById(id int) (entity.ShoppingCart, error) {
	var err error

	spcart := entity.ShoppingCart{}

	err = r.DB.Get(&spcart, getspcartById, id)

	spcart.Product, _ = r.prod.GetProductById(spcart.ProductId)

	return spcart, err
}

func (r *spcartrepo) Insertspcart(spcart *entity.ShoppingCart) error {
	_, err := r.DB.Exec(insertspcart, spcart.CartId, spcart.ProductId, spcart.QtyProduct, spcart.TotalPrice)

	return err
}

func (r *spcartrepo) Updatespcart(id int, spcart *entity.ShoppingCart) error {
	_, err := r.DB.Exec(updatespcart, spcart.CartId, spcart.ProductId, spcart.QtyProduct, spcart.TotalPrice, id)

	return err
}

func (r *spcartrepo) Deletespcart(id int) error {
	_, err := r.DB.Exec(deletespcart, id)

	return err
}
