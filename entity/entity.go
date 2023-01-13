package entity

import "time"

type Person struct {
	Id        string    `json:"id" validate:"isdefault" db:"id"`
	Name      string    `json:"name" validate:"required" db:"name"`
	Email     string    `json:"email" validate:"required,email" db:"email"`
	Password  string    `json:"password" validate:"required" db:"password"`
	CreatedAt time.Time `json:"created_at" validate:"isdefault" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" validate:"isdefault" db:"updated_at"`
}
