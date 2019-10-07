package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcamilom/ecommerce/context"
	"github.com/jcamilom/ecommerce/models"
)

// NewProducts is used to create a new Products controller
func NewProducts(ps models.ProductsService, us models.UserService) *Products {
	return &Products{
		ps: ps,
		us: us,
	}
}

type Products struct {
	ps models.ProductsService
	us models.UserService
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

func (p *Products) AddFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := context.User(r.Context())
	if user == nil {
		log.Println("Error while fetching the user from the context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fr := new(addFavoriteRequest)
	err := json.NewDecoder(r.Body).Decode(fr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if fr.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := p.ps.ByID(fr.ID)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&messageResponse{
				Message: "Product not found at the store's stock",
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	newFavorite := models.Favorite{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}
	err = p.us.AddFavorite(user, newFavorite)
	if err != nil {
		switch err {
		case models.ErrIsFavorite:
			w.WriteHeader(http.StatusNotModified)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&messageResponse{
		Message: fmt.Sprintf("Product with id '%v' added to the favorites list!", product.ID),
	})
}

type addFavoriteRequest struct {
	ID string `json:"id"`
}
