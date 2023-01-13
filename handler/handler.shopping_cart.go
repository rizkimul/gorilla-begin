package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
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
	} else {
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": spcart}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *spcarthandler) CreateSPCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var u entity.ShoppingCart
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		return
	}
	productId := strconv.Itoa(u.ProductId)
	cartid := strconv.Itoa(u.CartId)
	p, _ := h.prodSrvc.GetproductById(productId)
	c, _ := h.cartSrvc.GetById(cartid)
	v := reflect.ValueOf(c)
	if p == (entity.Product{}) {
		res := map[string]interface{}{"message": "Can't find product", "is_success": false, "status": "400"}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else if v.IsZero() {
		res := map[string]interface{}{"message": "Can't find cart", "is_success": false, "status": "400"}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
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
	productId := strconv.Itoa(u.ProductId)
	p, _ := h.prodSrvc.GetproductById(productId)
	u.TotalPrice = h.helper.CountTotal(u.QtyProduct, int(p.Price))
	log.Println(u.TotalPrice)
	err := h.srvc.Update(id, &u)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "Data Updated", "is_success": true, "status": "200"}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *spcarthandler) DeleteSPCart(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := h.srvc.Delete(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "Data Deleted", "is_success": true, "status": "200"}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}
