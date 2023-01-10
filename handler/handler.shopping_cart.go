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
	srvc   services.SPCartServices
	repos  repository.SPCartRepository
	helper helper.Helper
}

func NewSPCartHandler(srvc services.SPCartServices, repos repository.SPCartRepository, helper helper.Helper) SPCartHandler {
	return &spcarthandler{
		srvc:   srvc,
		repos:  repos,
		helper: helper,
	}
}

func (h *spcarthandler) GetSPCarts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	spcart, err := h.srvc.Getall()
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "200", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": spcart}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
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
	id := strconv.Itoa(u.ProductId)
	p, _ := repository.NewProductRepository(db).GetProductById(id)
	u.TotalPrice = h.helper.CountTotal(u.QtyProduct, int(p.Price))
	h.srvc.Insert(&u)
	res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200"}
	h.helper.ResponseJSON(w, http.StatusOK, res)

}

func (h *spcarthandler) GetSPCartbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	spcart, err := h.srvc.GetById(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusInternalServerError, res)
		return
	} else {
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": spcart}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *spcarthandler) UpdateSPCart(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var u entity.ShoppingCart
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err.Error())
	}
	err := h.srvc.Update(id, &u)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": u}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *spcarthandler) DeleteSPCart(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	delete, err := h.srvc.Delete(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": delete}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}
