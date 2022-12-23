package model

type Person struct {
	Id          string `json:"id" validate:"isdefault"`
	Name        string `json:"name" validate:"required,alpha"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	Phonenumber string `json:"phone" validate:"required,gte=11,lt=12,numeric"`
}

// var Persons make([]Person{})
