package response

type Cart struct {
	Id           int            `json:"id" db:"id"`
	CartName     string         `json:"cartName" db:"cart_name"`
	ShoppingCart []ShoppingCart `json:"shoppingCart"`
}
