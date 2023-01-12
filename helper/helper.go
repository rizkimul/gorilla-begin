package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/go-playground/validator/v10"
	"github.com/rizkimul/gorilla-begin/v2/response"
	"golang.org/x/crypto/bcrypt"
)

type Helper interface {
	Validation(u interface{}) []string
	CountTotal(qty int, price int) float64
	HashPass(pass string) (string, error)
	MatchPass(pass string, hash string) error
	ResponseJSON(w http.ResponseWriter, code int, payload interface{})
	PrintData(data *[]response.Response) ([]byte, error)
}

type helper struct{}

func NewHelper() Helper {
	return &helper{}
}

func (*helper) Validation(u interface{}) []string {
	errs := []string{}
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		validationerrors := err.(validator.ValidationErrors)
		for _, err := range validationerrors {
			switch err.Tag() {
			case "required":
				errs = append(errs, fmt.Errorf("%s is required", err.Field()).Error())
			case "alpha":
				errs = append(errs, fmt.Errorf("%s is alphabet only", err.Field()).Error())
			case "email":
				errs = append(errs, fmt.Errorf("%s is not valid email format", err.Field()).Error())
			case "numeric":
				errs = append(errs, fmt.Errorf("%s is numeric only", err.Field()).Error())
			case "gte":
				errs = append(errs, fmt.Errorf("%s value must be greater than %s", err.Field(), err.Param()).Error())
			case "lte":
				errs = append(errs, fmt.Errorf("%s value must lower than %s", err.Field(), err.Param()).Error())
			}
		}
	}
	return errs
}

func (*helper) CountTotal(qty int, price int) float64 {
	return float64(qty) * float64(price)
}

func (*helper) HashPass(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	return string(bytes), err
}

func (*helper) MatchPass(pass string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err
}

func (*helper) ResponseJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (*helper) PrintData(data *[]response.Response) ([]byte, error) {
	var tmpl *template.Template
	var err error

	var filepath = path.Join("public", "pdf-template", "print.html")
	tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		log.Println(err.Error())
	}

	var body bytes.Buffer

	if err = tmpl.Execute(&body, data); err != nil {
		return nil, err
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	page.EnableLocalFileAccess.Set(true)

	pdfg.AddPage(page)

	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)

	pdfg.Dpi.Set(300)

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	err = pdfg.Create()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return pdfg.Bytes(), nil
}
