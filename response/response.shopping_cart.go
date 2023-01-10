package response

type ShoppingCart struct {
	CartName     string  `db:"cart_name"`
	ProductName  string  `db:"product_name"`
	Price        float64 `db:"price"`
	ProductImage string  `db:"product_image"`
	QtyProduct   int     `db:"qty_product"`
	TotalPrice   float64 `db:"total_price"`
}
