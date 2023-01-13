package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	// "strings"

	"github.com/rizkimul/gorilla-begin/v2/config"
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/helper"
	"github.com/rizkimul/gorilla-begin/v2/repository"
	"github.com/rizkimul/gorilla-begin/v2/response"
	"github.com/rizkimul/gorilla-begin/v2/services"
	"github.com/rizkimul/gorilla-begin/v2/utils"
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
	token  utils.Token
}

func NewHandler(srvc services.Services, repos repository.Repository, helper helper.Helper, token utils.Token) Handler {
	return &handler{
		srvc:   srvc,
		repos:  repos,
		helper: helper,
		token:  token,
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
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(validate); err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		loginRes, _ := h.srvc.Login(p.Name)
		if loginRes == (entity.Person{}) {
			res := map[string]interface{}{"message": "No Data Match", "is_success": false, "status": "400"}
			h.helper.ResponseJSON(w, http.StatusBadRequest, res)
			return
		} else {
			if err := h.helper.MatchPass(p.Password, loginRes.Password); err != nil {
				res := map[string]interface{}{"message": "Unauthorized", "is_success": false, "status": "401", "data": err.Error()}
				h.helper.ResponseJSON(w, http.StatusUnauthorized, res)
				return
			} else {
				conf, err := config.LoadConfig(".")
				if err != nil {
					res := map[string]interface{}{"message": err.Error(), "is_success": false, "status": "401"}
					h.helper.ResponseJSON(w, http.StatusInternalServerError, res)
					return
				}
				accessToken, err := h.token.CreateToken(p.Name, conf.AccessTokenExp)
				if err != nil {
					res := map[string]interface{}{"message": err.Error(), "is_success": false, "status": "401"}
					h.helper.ResponseJSON(w, http.StatusInternalServerError, res)
					return
				}
				refreshToken, err := h.token.CreateToken(p.Name, conf.RefreshTokenExp)
				if err != nil {
					res := map[string]interface{}{"message": err.Error(), "is_success": false, "status": "401"}
					h.helper.ResponseJSON(w, http.StatusInternalServerError, res)
					return
				}
				res := map[string]interface{}{"message": "Token Generated", "is_success": true, "data": map[string]string{"Access Token": accessToken, "Refresh Token": refreshToken}, "status": "200"}
				h.helper.ResponseJSON(w, http.StatusOK, res)
				return
			}
		}
	}

}

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	person, err := h.srvc.Getall()
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
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
		return
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
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
	} else {
		passhash, _ := h.helper.HashPass(u.Password)
		u.Password = passhash
		h.srvc.Insert(&u)
		res := map[string]interface{}{"message": "Data Inserted", "is_success": true, "status": "201"}
		h.helper.ResponseJSON(w, http.StatusCreated, res)
		return
	}

}

func (h *handler) GetUserbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	person, err := h.srvc.GetById(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		result := response.Response{
			Id:    person.Id,
			Name:  person.Name,
			Email: person.Email,
		}
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": result}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
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
	} else {
		u.UpdatedAt = time.Now()
		passhash, _ := h.helper.HashPass(u.Password)
		u.Password = passhash
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

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	pdfBytes, err := h.helper.PrintData(&result)
	if err != nil {
		res := map[string]interface{}{"message": "Internal Server Error", "is_success": false, "status": "500", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusInternalServerError, res)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=user_data.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)

}
