package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/response"
)

type SPCartRepository interface {
	Getspcartall() ([]response.ShoppingCart, error)
	GetspcartById(id string) ([]entity.ShoppingCart, error)
	Insertspcart(spcart *entity.ShoppingCart) (*entity.ShoppingCart, error)
	Updatespcart(id string, spcart *entity.ShoppingCart) (int64, error)
	Deletespcart(id string) (int64, error)
}

type spcartrepo struct {
	DB *sqlx.DB
}

const (
	getspcartAll  = "SELECT a.cart_name, b.product_name, b.price, b.product_image, c.qty_product, c.total_price FROM cart a, product b, shopping_cart c where c.cart_id = a.id and c.product_id = b.id"
	getspcartById = "SELECT * FROM shopping_cart WHERE id=$1"
	insertspcart  = "INSERT INTO shopping_cart (cart_id, product_id, qty_product, total_price) VALUES ($1, $2, $3, $4)"
	updatespcart  = "UPDATE shopping_cart SET (cart_id, product_id, qty_product, total_price) = ($1, $2, $3, $4) WHERE id=$5"
	deletespcart  = "DELETE FROM shopping_cart WHERE id=$1"
)

func NewSPCartRepository(db *sqlx.DB) SPCartRepository {
	return &spcartrepo{
		DB: db,
	}
}

func (r *spcartrepo) Getspcartall() ([]response.ShoppingCart, error) {
	spcart := []response.ShoppingCart{}
	err := r.DB.Select(&spcart, getspcartAll)

	return spcart, err
}

func (r *spcartrepo) GetspcartById(id string) ([]entity.ShoppingCart, error) {
	spcart := []entity.ShoppingCart{}

	err := r.DB.Select(&spcart, getspcartById, id)
	return spcart, err
}

func (r *spcartrepo) Insertspcart(spcart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	var id string
	err := r.DB.QueryRow(insertspcart, spcart.Cart_id, spcart.Product_id, spcart.Qty_product, spcart.Total_price).Scan(&id)

	return spcart, err
}

func (r *spcartrepo) Updatespcart(id string, spcart *entity.ShoppingCart) (int64, error) {
	res, err := r.DB.Exec(updatespcart, spcart.Cart_id, spcart.Product_id, spcart.Qty_product, spcart.Total_price, id)

	rowsAfffected, err := res.RowsAffected()

	return rowsAfffected, err
}

func (r *spcartrepo) Deletespcart(id string) (int64, error) {
	res, err := r.DB.Exec(deletespcart, id)

	RowsAffected, err := res.RowsAffected()

	return RowsAffected, err
}
