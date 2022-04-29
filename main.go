package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"klbTest/app"
	"klbTest/models"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func newRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")

	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	r.HandleFunc("/login", app.GetLoginPageHandler).Methods("GET")
	r.HandleFunc("/login", app.LoginPostHandler(db)).Methods("POST")
	r.HandleFunc("/register", app.GetRegisterPageHandler).Methods("GET")
	r.HandleFunc("/register", app.RegisterPostHandler(db)).Methods("POST")

	r.HandleFunc("/produk", app.GetProdukPageHandler).Methods("GET")
	r.HandleFunc("/produk/add", app.GetProdukAddPageHandler).Methods("GET")

	r.HandleFunc("/api/produk/list", app.ListHandler(db)).Methods("GET")
	r.HandleFunc("/api/produk/getBySku", app.GetBySkuHandler(db)).Methods("GET")
	r.HandleFunc("/api/produk/upsert", app.UpsertProdukHandler(db)).Methods("POST")
	r.HandleFunc("/api/produk/delete", app.DeleteProdukHandler(db)).Methods("POST")

	// r.HandleFunc("/bird/list", getBirdHandler).Methods("GET")
	// r.HandleFunc("/bird", createBirdHandler).Methods("POST")
	return r
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// getting env variables SITE_TITLE and DB_HOST
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", dbUsername, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(1 * time.Second)
	db.SetConnMaxLifetime(30 * time.Second)

	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	fmt.Println("database is reachable")

	// The router is now formed by calling the `newRouter` constructor function
	// that we defined above. The rest of the code stays the same
	r := newRouter(db)
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func basicAuth(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))

			auth := r.Header.Get("Authorization")
			if auth == "" {
				return
			}
			authToken := strings.Replace(auth, "Basic ", "", 1)
			inputToken := models.InputGetToken{
				Token: authToken,
			}
			tokenData, err := models.GetToken(db, inputToken)
			if err != nil {
				fmt.Println(fmt.Errorf("Error: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			expectedUsernameHash := sha256.Sum256([]byte(tokenData.Email))
			expectedPasswordHash := sha256.Sum256([]byte(tokenData.Password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
