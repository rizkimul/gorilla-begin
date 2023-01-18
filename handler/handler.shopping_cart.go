package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	// "strings"

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
	srvc     services.SPCartServices
	repos    repository.SPCartRepository
	prodSrvc services.ProductServices
	cartSrvc services.CartServices
	helper   helper.Helper
}

func NewSPCartHandler(srvc services.SPCartServices, repos repository.SPCartRepository, helper helper.Helper, prodSrvc services.ProductServices, cartSrvc services.CartServices) SPCartHandler {
	return &spcarthandler{
		srvc:     srvc,
		repos:    repos,
		helper:   helper,
		prodSrvc: prodSrvc,
		cartSrvc: cartSrvc,
	}
}

func (h *spcarthandler) GetSPCarts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	spcart, err := h.srvc.Getall()
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "200", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": spcart}
	h.helper.ResponseJSON(w, http.StatusOK, res)
}

func (h *spcarthandler) CreateSPCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var u entity.ShoppingCart
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		return
	}
	err := h.srvc.Insert(&u)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200"}
	h.helper.ResponseJSON(w, http.StatusOK, res)

}

func (h *spcarthandler) GetSPCartbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	param := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(param)
	spcart, err := h.srvc.GetById(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusInternalServerError, res)
		return
	}
	res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": spcart}
	h.helper.ResponseJSON(w, http.StatusOK, res)
}

func (h *spcarthandler) UpdateSPCart(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(param)

	var u entity.ShoppingCart
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err.Error())
		return
	}
	err := h.srvc.Update(id, &u)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	res := map[string]interface{}{"message": "Data Updated", "is_success": true, "status": "200"}
	h.helper.ResponseJSON(w, http.StatusOK, res)
}

func (h *spcarthandler) DeleteSPCart(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(param)

	err := h.srvc.Delete(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	res := map[string]interface{}{"message": "Data Deleted", "is_success": true, "status": "200"}
	h.helper.ResponseJSON(w, http.StatusOK, res)
}
