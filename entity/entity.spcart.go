package entity

type SpCart struct {
	Id         int
	CartId     int     `json:"cart_id" validate:"required" db:"cart_id"`
	ProductId  int     `json:"product_id" validate:"required" db:"product_id"`
	QtyProduct int     `json:"qty" validate:"required" db:"qty_product"`
	TotalPrice float64 `json:"total" validate:"isdefault" db:"total_price"`
}
