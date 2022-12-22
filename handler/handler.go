package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "strings"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	Id          string `json:"id" validate:"isdefault"`
	Name        string `json:"firstname" validate:"required,alpha"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	Phonenumber string `json:"phone" validate:"required,gte=11,lt=12,numeric"`
}

var person = []Person{}

func Hashpassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	return string(bytes), err
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(person); i++ {
		hash, _ := Hashpassword(person[i].Password)
		person[i].Password = hash
	}
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(person)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error in encode response object", http.StatusBadRequest)
		return
	}
}

// func validateEmail(email string) bool {
// 	return !strings.Contains(email, "@")
// }

func CreateUser(w http.ResponseWriter, r *http.Request) {
	u := Person{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error in decode response object", http.StatusBadRequest)
		return
	}
	validate := validator.New()

	errs := []string{}
	err := validate.Struct(u)
	if err != nil {
		validationerrors := err.(validator.ValidationErrors)
		for _, err := range validationerrors {
			switch err.Tag() {
			case "required":
				errs = append(errs, fmt.Errorf("%s is required", err.Field()).Error())
			case "alpha":
				errs = append(errs, fmt.Errorf("%s is alphabet only", err.Field()).Error())
			case "email":
				errs = append(errs, fmt.Errorf("%s is not valid email format", err.Field()).Error())
			case "numeric":
				errs = append(errs, fmt.Errorf("%s is numeric only", err.Field()).Error())
			case "gte":
				errs = append(errs, fmt.Errorf("%s value must be greater than %s", err.Field(), err.Param()).Error())
			case "lte":
				errs = append(errs, fmt.Errorf("%s value must lower than %s", err.Field(), err.Param()).Error())
			}
		}
		if len(errs) > 0 {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(errs); err != nil {
				panic(err.Error())
			}
			return

		}
		// responsebody := map[string]string{"error": validationerrors.Error()}
	}

	// if u.Firstname == "" {
	// 	errs = append(errs, fmt.Errorf("firstname is required").Error())
	// }
	// if u.Lastname == "" {
	// 	errs = append(errs, fmt.Errorf("lastname is required").Error())
	// }
	// if u.Email == "" || validateEmail(u.Email) {
	// 	errs = append(errs, fmt.Errorf("valid email is required").Error())
	// }

	// if len(errs) > 0 {
	// 	w.Header().Add("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	if err := json.NewEncoder(w).Encode(errs); err != nil {
	// 	}
	// 	return
	// }
	u.Id = uuid.NewString()

	person = append(person, u)

	response, err := json.Marshal(&u)
	if err != nil {
		panic(err.Error())
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func indexbyID(persons []Person, id string) int {
	for i := 0; i < len(persons); i++ {
		if persons[i].Id == id {
			return i
		}
	}
	return -1
}

func GetUserbyId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	index := indexbyID(person, id)

	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(person[index]); err != nil {
		panic(err.Error())
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	index := indexbyID(person, id)

	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	u := Person{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		panic(err.Error())
	}

	person[index] = u
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&u); err != nil {
		panic(err.Error())
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	index := indexbyID(person, id)

	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
	}

	person = append(person[:index], person[index+1:]...)

	w.WriteHeader(http.StatusOK)
}
