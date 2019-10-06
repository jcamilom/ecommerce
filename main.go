package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jcamilom/ecommerce/controllers"
	"github.com/jcamilom/ecommerce/middleware"
	"github.com/jcamilom/ecommerce/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	port = 3000
)

func main() {
	loadEnvVars()

	us := models.NewUserService()
	usersC := controllers.NewUsers(us)
	ps := models.NewProductsService()
	productsC := controllers.NewProducts(ps)

	requireUserMw := middleware.RequireUser{
		UserService: us,
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/register", usersC.Create).Methods("POST")
	r.HandleFunc("/products/{id}", requireUserMw.ApplyFn(productsC.GetProduct)).Methods("GET")
	fmt.Printf("Starting the server on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
