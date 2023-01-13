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

type CartHandler interface {
	GetCarts(w http.ResponseWriter, r *http.Request)
	CreateCart(w http.ResponseWriter, r *http.Request)
	GetCartbyId(w http.ResponseWriter, r *http.Request)
	UpdateCart(w http.ResponseWriter, r *http.Request)
	DeleteCart(w http.ResponseWriter, r *http.Request)
}

type carthandler struct {
	srvc     services.CartServices
	repos    repository.CartRepository
	helper   helper.Helper
	prodSrvc services.ProductServices
}

// var srvc services.Services = services.NewServices()

func NewCartHandler(srvc services.CartServices, repos repository.CartRepository, helper helper.Helper, prodSrvc services.ProductServices) CartHandler {
	return &carthandler{
		srvc:     srvc,
		repos:    repos,
		helper:   helper,
		prodSrvc: prodSrvc,
	}
}

func (h *carthandler) GetCarts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	cart, err := h.srvc.Getall()
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": cart}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *carthandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var u entity.Cart
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		return
	}

	validate := h.helper.Validation(u)
	if len(validate) > 0 {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": validate}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		h.srvc.Insert(&u)
		res := map[string]interface{}{"message": "Data Successfully Inserted", "is_success": true, "status": "201"}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}

}

func (h *carthandler) GetCartbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	cart, err := h.srvc.GetById(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	spCart := []entity.ShoppingCart{}
	for _, v := range cart.ShoppingCarts {
		productId := strconv.Itoa(v.ProductId)
		v.Product, _ = h.prodSrvc.GetproductById(productId)
		resp := entity.ShoppingCart{
			Id:         v.Id,
			CartId:     v.CartId,
			ProductId:  v.ProductId,
			Product:    v.Product,
			QtyProduct: v.QtyProduct,
			TotalPrice: v.TotalPrice,
		}
		spCart = append(spCart, resp)
	}
	cart.ShoppingCarts = spCart
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": cart}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *carthandler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	var u entity.Cart
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err.Error())
		return
	}

	validate := h.helper.Validation(u)
	if len(validate) > 0 {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": validate}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		_, err := h.srvc.Update(id, &u)
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
}

func (h *carthandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := h.srvc.Delete(id)
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
