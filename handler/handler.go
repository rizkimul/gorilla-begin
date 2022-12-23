package handler

import (
	"encoding/json"
	"log"
	"net/http"

	// "strings"
	"github.com/google/uuid"
	"github.com/rizkimul/gorilla-begin/v2/helper"
	"github.com/rizkimul/gorilla-begin/v2/model"
	"github.com/rizkimul/gorilla-begin/v2/response"
	"golang.org/x/crypto/bcrypt"
)

var person = []model.Person{}

func Hashpassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	return string(bytes), err
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(person); i++ {
		hash, _ := Hashpassword(person[i].Password)
		person[i].Password = hash
	}
	result := []response.Response{}

	for _, v := range *&person {
		res := response.Response{
			Id:          v.Id,
			Name:        v.Name,
			Email:       v.Email,
			Phonenumber: v.Phonenumber,
		}
		*&result = append(*&result, res)
	}

	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	u := model.Person{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		return
	}

	validate := helper.Validation(u)
	if len(validate) > 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(validate); err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		u.Id = uuid.NewString()

		person = append(person, u)

		response, err := json.Marshal(&u)
		if err != nil {
			log.Println(err.Error())
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}

}

func GetUserbyId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	index := helper.IndexbyID(person, id)

	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(person[index]); err != nil {
		log.Println(err.Error())
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	index := helper.IndexbyID(person, id)
	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	u := model.Person{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err.Error())
	}

	validate := helper.Validation(u)
	if len(validate) > 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(validate); err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		person[index] = u
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&u); err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	index := helper.IndexbyID(person, id)

	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
	}

	person = append(person[:index], person[index+1:]...)

	w.WriteHeader(http.StatusOK)
}
