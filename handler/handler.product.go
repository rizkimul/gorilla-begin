package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	// "strings"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/schema"
	"github.com/rizkimul/gorilla-begin/v2/entity"
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
}

func NewProductHandler(srvc services.ProductServices, productRepo repository.RepositoryProduct) ProductHandler {
	return &productHandler{
		srvc:        srvc,
		productRepo: productRepo,
	}
}

func (h *productHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	product, err := h.srvc.Getallproduct()
	if err != nil {
		log.Println(err)
		return
	}
	result := []response.Product{}

	for _, v := range product {
		res := response.Product{
			Id:                  v.Id,
			Product_name:        v.Product_name,
			Product_description: v.Product_description,
			Product_image:       v.Product_image,
			Price:               v.Price,
		}
		result = append(result, res)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var urlcloud = "cloudinary://475673691162386:yrcOGG9UdcYjqzruVtSq4uLWRzU@db9xlikgu"
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
	product.Product_image = resp.SecureURL
	if err != nil {
		log.Println(err.Error())
		return
	}

	h.srvc.Insertproduct(&product)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)

}

func (h *productHandler) GetProductbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	product, err := h.srvc.GetproductById(id)
	if err != nil {
		log.Println(err.Error())
		return
	}
	result := []response.Product{}

	for _, v := range product {
		res := response.Product{
			Id:                  v.Id,
			Product_name:        v.Product_name,
			Product_description: v.Product_description,
			Product_image:       v.Product_image,
			Price:               v.Price,
		}
		result = append(result, res)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var urlcloud = "cloudinary://475673691162386:yrcOGG9UdcYjqzruVtSq4uLWRzU@db9xlikgu"
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
	product.Product_image = resp.SecureURL
	if err != nil {
		log.Println(err.Error())
		return
	}

	h.srvc.Updateproduct(id, &product)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)

}

func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	h.srvc.Deleteproduct(id)

	w.WriteHeader(http.StatusOK)
}
