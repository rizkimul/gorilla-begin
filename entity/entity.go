package entity

type Person struct {
	Id          string `json:"id" validate:"isdefault" db:"id"`
	Name        string `json:"name" validate:"required" db:"name"`
	Email       string `json:"email" validate:"required,email" db:"email"`
	Password    string `json:"password" validate:"required" db:"password"`
	PhoneNumber string `json:"phone" validate:"required,gte=11,lt=12,numeric" db:"phonenumber"`
	CreatedAt   string `json:"created_at" validate:"isdefault" db:"created_at"`
}
