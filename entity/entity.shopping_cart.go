package entity

type ShoppingCart struct {
	Cart_id     int     `json:"cart_id" validate:"required"`
	Product_id  int     `json:"product_id" validate:"required"`
	Qty_product int     `json:"qty" validate:"required"`
	Total_price float64 `json:"total" validate:"isdefault"`
}
