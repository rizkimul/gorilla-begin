package entity

import "time"

type Product struct {
	Id                 string    `schema:"id" validate:"isdefault" db:"id"`
	ProductName        string    `schema:"product_name" validate:"required" db:"product_name"`
	ProductDescription string    `schema:"product_desc" validate:"required" db:"product_description"`
	Price              float64   `schema:"price" validate:"required,numeric" db:"price"`
	ProductImage       string    `schema:"image" validate:"required" db:"product_image"`
	CreatedAt          time.Time `schema:"created_at" db:"created_at"`
}
