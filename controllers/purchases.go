package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jcamilom/ecommerce/context"
	"github.com/jcamilom/ecommerce/models"
)

// NewPurchases is used to create a new Purchases controller
func NewPurchases(ps models.PurchaseService, prs models.ProductsService) *Purchases {
	return &Purchases{
		ps:  ps,
		prs: prs,
	}
}

type Purchases struct {
	ps  models.PurchaseService
	prs models.ProductsService
}

// Create registers a new purchase
//
// POST /purchases
func (p *Purchases) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := context.User(r.Context())
	if user == nil {
		log.Println("Error while fetching the user from the context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pr := new(createPurchaseRequest)
	err := json.NewDecoder(r.Body).Decode(pr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if pr.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := p.prs.ByID(pr.ID)
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
	purchase := &models.Purchase{
		ID:    "1000",
		Email: "juan@mail.com",
		Item: models.PurchaseItem{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price,
		},
	}
	err = p.ps.Create(purchase)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Println("Purchase created")
}

type createPurchaseRequest struct {
	ID string `json:"id"`
}
