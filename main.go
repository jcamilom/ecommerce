package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	port = 3000
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	usr := new(User)
	found, err := getItem("Email", vars["user"], "Users", usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else if found == false {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(usr)
	}
}

func main() {
	loadEnvVars()
	r := mux.NewRouter()
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
