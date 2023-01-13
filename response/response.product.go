package response

type Product struct {
	Id                 string  `json:"id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"produdct_description"`
	Price              float64 `json:"price"`
	ProductImage       string  `json:"image"`
}
