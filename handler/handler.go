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

type Handler interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserbyId(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	PrinData(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	srvc   services.Services
	repos  repository.Repository
	helper helper.Helper
}

func NewHandler(srvc services.Services, repos repository.Repository, helper helper.Helper) Handler {
	return &handler{
		srvc:   srvc,
		repos:  repos,
		helper: helper,
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var p entity.Login
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println(err.Error())
		return
	}
	validate := h.helper.Validation(p)
	if len(validate) > 0 {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": validate}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}

	loginRes, err := h.srvc.Login(p.Email, p.Password)
	if err != nil {
		res := map[string]interface{}{"message": "No Data Match", "is_success": false, "status": "400"}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	res := map[string]interface{}{"message": "Token Generated", "is_success": true, "data": loginRes, "status": "200"}
	h.helper.ResponseJSON(w, http.StatusOK, res)
}

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	person, err := h.srvc.Getall()
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}

	result := []response.Response{}

	for _, v := range person {
		res := response.Response{
			Id:    v.Id,
			Name:  v.Name,
			Email: v.Email,
		}
		result = append(result, res)
	}
	res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": result}
	h.helper.ResponseJSON(w, http.StatusOK, res)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u entity.Person
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		return
	}

	validate := h.helper.Validation(u)
	if len(validate) > 0 {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": validate}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	err := h.srvc.Insert(&u)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	res := map[string]interface{}{"message": "Data Inserted", "is_success": true, "status": "201"}
	h.helper.ResponseJSON(w, http.StatusCreated, res)
}

func (h *handler) GetUserbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	person, err := h.srvc.GetById(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	result := response.Response{
		Id:    person.Id,
		Name:  person.Name,
		Email: person.Email,
	}
	res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": result}
	h.helper.ResponseJSON(w, http.StatusOK, res)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var u entity.Person
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err.Error())
	}

	validate := h.helper.Validation(u)
	if len(validate) > 0 {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": validate}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
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

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := h.srvc.Delete(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}
	res := map[string]interface{}{"message": "Data Deleted", "is_success": true, "status": "200"}
	h.helper.ResponseJSON(w, http.StatusOK, res)
}

func (h *handler) PrinData(w http.ResponseWriter, r *http.Request) {

	person, err := h.srvc.Getall()
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	}

	result := []response.Response{}

	for _, v := range person {
		res := response.Response{
			Id:    v.Id,
			Name:  v.Name,
			Email: v.Email,
		}
		result = append(result, res)
	}

	pdf, err := h.srvc.Print(result)
	if err != nil {
		res := map[string]interface{}{"message": "Internal Server Error", "is_success": false, "status": "500", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusInternalServerError, res)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=user_data.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdf)

}
