package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rizkimul/gorilla-begin/v2/entity"
)

type CartRepository interface {
	Getcartall() ([]entity.Cart, error)
	GetcartById(id int) (entity.Cart, error)
	Insertcart(cart *entity.Cart) error
	Updatecart(id int, cart *entity.Cart) error
	Deletecart(id int) error
}

type cartrepo struct {
	DB     *sqlx.DB
	spRepo SPCartRepository
	prod   RepositoryProduct
}

const (
	getcartAll  = "SELECT * FROM cart"
	getcartById = "SELECT * FROM cart WHERE id=$1"
	insertcart  = "INSERT INTO cart (cart_name) VALUES ($1)"
	updatecart  = "UPDATE cart SET cart_name = $1 WHERE id=$2"
	deletecart  = "DELETE FROM cart WHERE id=$1"
	getspcart   = "SELECT * FROM shopping_cart WHERE cart_id=$1"
)

func NewCartRepository(db *sqlx.DB, spRepo SPCartRepository, prod RepositoryProduct) CartRepository {
	return &cartrepo{
		DB:     db,
		spRepo: spRepo,
		prod:   prod,
	}
}

func (r *cartrepo) Getcartall() ([]entity.Cart, error) {
	cart := []entity.Cart{}
	err := r.DB.Select(&cart, getcartAll)

	return cart, err
}

func (r *cartrepo) GetcartById(id int) (entity.Cart, error) {
	var err error
	cart := entity.Cart{}

	spCart := []entity.ShoppingCart{}

	err = r.DB.Get(&cart, getcartById, id)

	err = r.DB.Select(&spCart, getspcart, cart.Id)

	cart.ShoppingCarts = append(cart.ShoppingCarts, spCart...)

	return cart, err
}

func (r *cartrepo) Insertcart(cart *entity.Cart) error {
	_, err := r.DB.Exec(insertcart, cart.CartName)

	return err
}

func (r *cartrepo) Updatecart(id int, cart *entity.Cart) error {
	_, err := r.DB.Exec(updatecart, cart.CartName, id)

	return err
}

func (r *cartrepo) Deletecart(id int) error {
	_, err := r.DB.Exec(deletecart, id)

	return err
}
