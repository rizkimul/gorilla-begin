package response

type Product struct {
	Id                  string  `json:"id"`
	Product_name        string  `json:"product_name"`
	Product_description string  `json:"produdct_description"`
	Price               float64 `json:"price"`
	Product_image       string  `json:"image"`
}
