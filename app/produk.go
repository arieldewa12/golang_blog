package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"klbTest/models"
	"net/http"
	"path"
	"strconv"
)

func GetProdukPageHandler(w http.ResponseWriter, r *http.Request) {

	filepath := path.Join("assets/pages/produk", "produk.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetProdukAddPageHandler(w http.ResponseWriter, r *http.Request) {

	filepath := path.Join("assets/pages/produk", "add.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ListHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		// In case of any error, we respond with an error to the user
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the information about the bird from the form info
		inputList := models.InputList{
			Offset: 0,
			Limit:  50,
		}
		produkData, err := models.List(db, inputList)
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(produkData)
		produkJson, err := json.Marshal(produkData)
		w.WriteHeader(http.StatusOK)
		w.Write(produkJson)
	}

	return http.HandlerFunc(fn)
}

func GetBySkuHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		// In case of any error, we respond with an error to the user
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the information about the bird from the form info
		inputGetBySKU := models.InputGetBySKU{
			SKU: r.Form.Get("sku"),
		}
		produkData, err := models.GetBySKU(db, inputGetBySKU)
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(produkData)
		produkJson, err := json.Marshal(produkData)
		w.WriteHeader(http.StatusOK)
		w.Write(produkJson)
	}

	return http.HandlerFunc(fn)
}

func UpsertProdukHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("masuk")
		err := r.ParseForm()

		// In case of any error, we respond with an error to the user
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the information about the bird from the form info
		price, err := strconv.Atoi(r.Form.Get("price"))
		isActive, _ := strconv.ParseBool(r.Form.Get("isActive"))
		inputUpsert := models.InputUpsert{
			Produk: models.Produk{
				ID:       r.Form.Get("id"),
				SKU:      r.Form.Get("sku"),
				Name:     r.Form.Get("name"),
				Price:    price,
				Unit:     r.Form.Get("unit"),
				IsActive: isActive,
			},
		}
		fmt.Println(inputUpsert)
		err = models.Upsert(db, inputUpsert)
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	return http.HandlerFunc(fn)
}

func DeleteProdukHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		// In case of any error, we respond with an error to the user
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the information about the bird from the form info
		inputDelete := models.InputDelete{
			ID: r.Form.Get("id"),
		}
		err = models.Delete(db, inputDelete)
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	return http.HandlerFunc(fn)
}
