package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcamilom/ecommerce/context"
	"github.com/jcamilom/ecommerce/models"
)

// NewProducts is used to create a new Products controller
func NewProducts(ps models.ProductsService) *Products {
	return &Products{
		ps: ps,
	}
}

type Products struct {
	ps models.ProductsService
}

func (p *Products) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := context.User(r.Context())
	if user == nil {
		log.Println("Error while fetching the user from the context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	product, err := p.ps.ByID(vars["id"])
	if err != nil {
		if err == models.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		json.NewEncoder(w).Encode(product)
	}
}
