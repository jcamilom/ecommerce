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

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	usr, err := getUser(vars["user"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else if usr != nil {
		json.NewEncoder(w).Encode(usr)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users/{user}", getUserHandler).Methods("GET")
	http.ListenAndServe(":3000", r)
}
