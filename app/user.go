package app

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"klbTest/models"
	"net/http"
	"path"
)

func GetLoginPageHandler(w http.ResponseWriter, r *http.Request) {

	filepath := path.Join("assets/pages", "login.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetRegisterPageHandler(w http.ResponseWriter, r *http.Request) {

	filepath := path.Join("assets/pages", "register.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RegisterPostHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		// In case of any error, we respond with an error to the user
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the information about the bird from the form info
		inputRegister := models.InputRegister{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		}
		err = models.Register(db, inputRegister)
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", 301)
	}

	return http.HandlerFunc(fn)
}

func LoginPostHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		// In case of any error, we respond with an error to the user
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the information about the bird from the form info
		inputLogin := models.InputLogin{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		}
		userData, err := models.Login(db, inputLogin)
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tokenVar := r.Form.Get("email") + ":" + r.Form.Get("password")
		token := base64.StdEncoding.EncodeToString([]byte(tokenVar))
		inputToken := models.InputCreateToken{
			Token:    token,
			UserId:   userData.ID,
			IsActive: true,
		}
		err = models.CreateToken(db, inputToken)
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(userData)
		http.Redirect(w, r, "/produk", 301)
	}

	return http.HandlerFunc(fn)
}

func GetUserTokenHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		// In case of any error, we respond with an error to the user
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the information about the bird from the form info
		inputToken := models.InputGetToken{
			Token: r.Form.Get("token"),
		}
		tokenData, err := models.GetToken(db, inputToken)
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(tokenData)
		// produkJson, err := json.Marshal(tokenData)
		w.WriteHeader(http.StatusOK)
		// w.Write(produkJson)
	}

	return http.HandlerFunc(fn)
}
