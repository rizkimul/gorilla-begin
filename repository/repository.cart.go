package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rizkimul/gorilla-begin/v2/entity"
)

type CartRepository interface {
	Getcartall() ([]entity.Cart, error)
	GetcartById(id string) ([]entity.Cart, error)
	Insertcart(cart *entity.Cart) (*entity.Cart, error)
	Updatecart(id string, cart *entity.Cart) (int64, error)
	Deletecart(id string) (int64, error)
}

type cartrepo struct {
	DB *sqlx.DB
}

const (
	getcartAll  = "SELECT * FROM cart"
	getcartById = "SELECT * FROM cart WHERE id=$1"
	insertcart  = "INSERT INTO cart (cart_name) VALUES ($1)"
	updatecart  = "UPDATE cart SET cart_name = $1 WHERE id=$2"
	deletecart  = "DELETE FROM cart WHERE id=$1"
)

func NewCartRepository(db *sqlx.DB) CartRepository {
	return &cartrepo{
		DB: db,
	}
}

func (r *cartrepo) Getcartall() ([]entity.Cart, error) {
	cart := []entity.Cart{}
	err := r.DB.Select(&cart, getcartAll)

	return cart, err
}

func (r *cartrepo) GetcartById(id string) ([]entity.Cart, error) {
	cart := []entity.Cart{}

	err := r.DB.Select(&cart, getcartById, id)
	return cart, err
}

func (r *cartrepo) Insertcart(cart *entity.Cart) (*entity.Cart, error) {
	var id string
	err := r.DB.QueryRow(insertcart, cart.Cart_name).Scan(&id)

	return cart, err
}

func (r *cartrepo) Updatecart(id string, cart *entity.Cart) (int64, error) {
	res, err := r.DB.Exec(updatecart, cart.Cart_name, id)

	rowsAfffected, err := res.RowsAffected()

	return rowsAfffected, err
}

func (r *cartrepo) Deletecart(id string) (int64, error) {
	res, err := r.DB.Exec(deletecart, id)

	RowsAffected, err := res.RowsAffected()

	return RowsAffected, err
}
