package entity

type Cart struct {
	Id        int    `json:"id" validate:"isdefault"`
	Cart_name string `json:"cart_name" validate:"required"`
}
