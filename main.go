package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jcamilom/ecommerce/controllers"
	"github.com/jcamilom/ecommerce/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	port = 3000
)

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	us := models.NewUserService()
	user, err := us.ByEmail(vars["user"])
	if err != nil {
		if err == models.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		json.NewEncoder(w).Encode(user)
	}
}

func main() {
	loadEnvVars()

	us := models.NewUserService()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/users", usersC.Create).Methods("POST")
	r.HandleFunc("/users/{user}", getUserHandler).Methods("GET")
	fmt.Printf("Starting the server on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
