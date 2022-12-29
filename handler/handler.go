package handler

import (
	"encoding/json"
	"log"
	"net/http"

	// "strings"
	"github.com/google/uuid"
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/helper"
	"github.com/rizkimul/gorilla-begin/v2/repository"
	"github.com/rizkimul/gorilla-begin/v2/response"
	"github.com/rizkimul/gorilla-begin/v2/services"
)

type Handler interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserbyId(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	srvc  services.Services
	repos repository.Repository
}

// var srvc services.Services = services.NewServices()

func NewHandler(srvc services.Services, repos repository.Repository) Handler {
	return &handler{
		srvc:  srvc,
		repos: repos,
	}
}

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	person, err := h.srvc.Getall()
	if err != nil {
		log.Println(err)
		return
	}
	result := []response.Response{}

	for _, v := range person {
		res := response.Response{
			Id:          v.Id,
			Name:        v.Name,
			Email:       v.Email,
			Phonenumber: v.Phonenumber,
		}
		result = append(result, res)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var u entity.Person
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		return
	}

	validate := helper.Validation(u)
	if len(validate) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(validate); err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		u.Id = uuid.NewString()

		h.srvc.Insert(&u)
	}

}

func (h *handler) GetUserbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	person, err := h.srvc.GetById(id)
	if err != nil {
		log.Println(err.Error())
		return
	}
	result := []response.Response{}

	for _, v := range person {
		res := response.Response{
			Id:          v.Id,
			Name:        v.Name,
			Email:       v.Email,
			Phonenumber: v.Phonenumber,
		}
		result = append(result, res)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var u entity.Person
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
		h.srvc.Update(id, &u)
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&u); err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	h.srvc.Delete(id)

	w.WriteHeader(http.StatusOK)
}
