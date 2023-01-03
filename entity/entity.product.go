package entity

import "time"

type Product struct {
	Id                  string    `schema:"id" validate:"isdefault"`
	Product_name        string    `schema:"product_name" validate:"required"`
	Product_description string    `schema:"product_desc" validate:"required"`
	Price               float64   `schema:"price" validate:"required,numeric"`
	Product_image       string    `schema:"image" validate:"required"`
	Created_at          time.Time `schema:"created_at"`
}
