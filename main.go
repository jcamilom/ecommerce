package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users = []User{
	{ID: "1", Name: "juan", Email: "juan@mail.com", Password: "1234"},
	{ID: "2", Name: "pedro", Email: "pedro@mail.com", Password: "1234"},
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", getUsersHandler).Methods("GET")
	http.ListenAndServe(":3000", r)
}
