package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	// "strings"

	"github.com/jmoiron/sqlx"
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/helper"
	"github.com/rizkimul/gorilla-begin/v2/repository"
	"github.com/rizkimul/gorilla-begin/v2/services"
)

type SPCartHandler interface {
	GetSPCarts(w http.ResponseWriter, r *http.Request)
	CreateSPCart(w http.ResponseWriter, r *http.Request)
	GetSPCartbyId(w http.ResponseWriter, r *http.Request)
	UpdateSPCart(w http.ResponseWriter, r *http.Request)
	DeleteSPCart(w http.ResponseWriter, r *http.Request)
}

type spcarthandler struct {
	srvc  services.SPCartServices
	repos repository.SPCartRepository
}

func NewSPCartHandler(srvc services.SPCartServices, repos repository.SPCartRepository) SPCartHandler {
	return &spcarthandler{
		srvc:  srvc,
		repos: repos,
	}
}

func (h *spcarthandler) GetSPCarts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	spcart, err := h.srvc.Getall()
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(spcart)
}

func (h *spcarthandler) CreateSPCart(w http.ResponseWriter, r *http.Request) {
	dsn := "user=postgres password=root dbname=db_golang sslmode=disable"
	w.Header().Add("Content-Type", "application/json")
	db, _ := sqlx.Connect("postgres", dsn)
	var u entity.ShoppingCart
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		return
	}
	id := strconv.Itoa(u.Product_id)
	p, _ := repository.NewProductRepository(db).GetProductById(id)

	for _, v := range p {
		total := u.Qty_product * int(v.Price)
		u.Total_price = float64(total)
	}

	validate := helper.Validation(u)
	if len(validate) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(validate); err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		insert, err := h.srvc.Insert(&u)
		if err != nil {
			log.Println(err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(insert)
	}

}

func (h *spcarthandler) GetSPCartbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	spcart, err := h.srvc.GetById(id)
	if err != nil {
		log.Println(err.Error())
		return
	}
	json.NewEncoder(w).Encode(spcart)
}

func (h *spcarthandler) UpdateSPCart(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var u entity.ShoppingCart
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
		update, err := h.srvc.Update(id, &u)
		if err != nil {
			log.Println(err.Error())
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(update)
	}
}

func (h *spcarthandler) DeleteSPCart(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	delete, err := h.srvc.Delete(id)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(delete)
}
