package entity

type Cart struct {
	Id            int            `json:"id" validate:"isdefault"`
	CartName      string         `json:"cart_name" validate:"required" db:"cart_name"`
	ShoppingCarts []ShoppingCart `json:"shopping_cart"`
}
