package entity

type Login struct {
	Name     string `json:"name" validate:"required,alpha"`
	Password string `json:"password" validate:"required"`
}
