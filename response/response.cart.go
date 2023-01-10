package response

type Cart struct {
	Status     int         `json:"status"`
	Is_success bool        `json:"is_success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
