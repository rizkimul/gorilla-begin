package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	// "strings"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/schema"
	"github.com/rizkimul/gorilla-begin/v2/config"
	"github.com/rizkimul/gorilla-begin/v2/entity"
	"github.com/rizkimul/gorilla-begin/v2/helper"
	"github.com/rizkimul/gorilla-begin/v2/repository"
	"github.com/rizkimul/gorilla-begin/v2/response"
	"github.com/rizkimul/gorilla-begin/v2/services"
)

type ProductHandler interface {
	GetProducts(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProductbyId(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
}

type productHandler struct {
	srvc        services.ProductServices
	productRepo repository.RepositoryProduct
	helper      helper.Helper
}

func NewProductHandler(srvc services.ProductServices, productRepo repository.RepositoryProduct, helper helper.Helper) ProductHandler {
	return &productHandler{
		srvc:        srvc,
		productRepo: productRepo,
		helper:      helper,
	}
}

func (h *productHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	product, err := h.srvc.Getallproduct()
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		result := []response.Product{}

		for _, v := range product {
			res := response.Product{
				Id:                 v.Id,
				ProductName:        v.ProductName,
				ProductDescription: v.ProductDescription,
				ProductImage:       v.ProductImage,
				Price:              v.Price,
			}
			result = append(result, res)
		}
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": result}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	conf, _ := config.LoadConfig(".")
	var decoder = schema.NewDecoder()
	var product entity.Product
	cld, _ := cloudinary.NewFromURL(conf.CloudSecretKey)
	err := r.ParseMultipartForm(1 << 2)
	if err != nil {
		log.Println(err.Error())
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Println(err.Error())
		return
	}
	ctx := context.Background()
	resp, _ := cld.Upload.Upload(ctx, file, uploader.UploadParams{})
	err = decoder.Decode(&product, r.PostForm)
	product.ProductImage = resp.SecureURL
	if err != nil {
		log.Println(err.Error())
		return
	}
	product.CreatedAt = time.Now()
	h.srvc.Insertproduct(&product)
	respond := response.Product{
		Id:                 product.Id,
		ProductName:        product.ProductName,
		ProductDescription: product.ProductDescription,
		Price:              product.Price,
		ProductImage:       product.ProductImage,
	}
	res := map[string]interface{}{"message": "Data Successfully Inserted", "is_success": true, "status": "201", "data": respond}
	h.helper.ResponseJSON(w, http.StatusOK, res)

}

func (h *productHandler) GetProductbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	product, err := h.srvc.GetproductById(id)
	respond := response.Product{
		Id:                 product.Id,
		ProductName:        product.ProductName,
		ProductDescription: product.ProductDescription,
		Price:              product.Price,
		ProductImage:       product.ProductImage,
	}
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": respond}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	conf, _ := config.LoadConfig(".")
	var urlcloud = conf.CloudSecretKey
	var decoder = schema.NewDecoder()
	var product entity.Product
	cld, _ := cloudinary.NewFromURL(urlcloud)
	err := r.ParseMultipartForm(1 << 2)
	if err != nil {
		log.Println(err.Error())
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Println(err.Error())
		return
	}
	ctx := context.Background()
	resp, _ := cld.Upload.Upload(ctx, file, uploader.UploadParams{})
	err = decoder.Decode(&product, r.PostForm)
	product.ProductImage = resp.SecureURL
	product.UpdatedAt = time.Now()
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = h.srvc.Updateproduct(id, &product)
	respond := response.Product{
		Id:                 product.Id,
		ProductName:        product.ProductName,
		ProductDescription: product.ProductDescription,
		Price:              product.Price,
		ProductImage:       product.ProductImage,
	}
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "Data Updated", "is_success": true, "status": "200", "data": respond}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := h.srvc.Deleteproduct(id)
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
