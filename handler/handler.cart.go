package handler

import (
	"encoding/json"
	"log"
	"net/http"

	// "strings"

	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/helper"
	"github.com/rizkimul/gorilla-begin/v2/repository"
	"github.com/rizkimul/gorilla-begin/v2/response"
	"github.com/rizkimul/gorilla-begin/v2/services"
)

type CartHandler interface {
	GetCarts(w http.ResponseWriter, r *http.Request)
	CreateCart(w http.ResponseWriter, r *http.Request)
	GetCartbyId(w http.ResponseWriter, r *http.Request)
	UpdateCart(w http.ResponseWriter, r *http.Request)
	DeleteCart(w http.ResponseWriter, r *http.Request)
}

type carthandler struct {
	srvc  services.CartServices
	repos repository.CartRepository
}

// var srvc services.Services = services.NewServices()

func NewCartHandler(srvc services.CartServices, repos repository.CartRepository) CartHandler {
	return &carthandler{
		srvc:  srvc,
		repos: repos,
	}
}

func (h *carthandler) GetCarts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	person, err := h.srvc.Getall()
	if err != nil {
		log.Println(err)
		return
	}
	result := []response.Cart{}

	for _, v := range person {
		res := response.Cart{
			Id:        v.Id,
			Cart_name: v.Cart_name,
		}
		result = append(result, res)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *carthandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var u entity.Cart
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
		h.srvc.Insert(&u)
	}

}

func (h *carthandler) GetCartbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	person, err := h.srvc.GetById(id)
	if err != nil {
		log.Println(err.Error())
		return
	}
	result := []response.Cart{}

	for _, v := range person {
		res := response.Cart{
			Id:        v.Id,
			Cart_name: v.Cart_name,
		}
		result = append(result, res)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *carthandler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var u entity.Cart
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

func (h *carthandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	h.srvc.Delete(id)

	w.WriteHeader(http.StatusOK)
}
