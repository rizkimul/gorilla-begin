package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	// "strings"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/schema"
	"github.com/joho/godotenv"
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
	} else {
		result := []response.Product{}

		for _, v := range product {
			res := response.Product{
				Id:                  v.Id,
				Product_name:        v.ProductName,
				Product_description: v.ProductDescription,
				Product_image:       v.ProductImage,
				Price:               v.Price,
			}
			result = append(result, res)
		}
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": result}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	_ = godotenv.Load(".env")
	var urlcloud = os.Getenv("CLOUD_SECRET_KEY")
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
	if err != nil {
		log.Println(err.Error())
		return
	}
	product.CreatedAt = time.Now()
	h.srvc.Insertproduct(&product)
	res := map[string]interface{}{"message": "Data Successfully Inserted", "is_success": true, "status": "201", "data": product}
	h.helper.ResponseJSON(w, http.StatusOK, res)

}

func (h *productHandler) GetProductbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	product, err := h.srvc.GetproductById(id)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "OK", "is_success": true, "status": "200", "data": product}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
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
	product.ProductImage = resp.SecureURL
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = h.srvc.Updateproduct(id, &product)
	if err != nil {
		res := map[string]interface{}{"message": "Bad Request", "is_success": false, "status": "400", "data": err.Error()}
		h.helper.ResponseJSON(w, http.StatusBadRequest, res)
		return
	} else {
		res := map[string]interface{}{"message": "Data Updated", "is_success": true, "status": "200", "data": product}
		h.helper.ResponseJSON(w, http.StatusOK, res)
		return
	}
}

func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, err := h.srvc.Deleteproduct(id)
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
