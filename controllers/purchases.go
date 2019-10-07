package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jcamilom/ecommerce/context"
	"github.com/jcamilom/ecommerce/models"
)

// NewPurchases is used to create a new Purchases controller
func NewPurchases(pus models.PurchaseService, ps models.ProductsService, us models.UserService) *Purchases {
	return &Purchases{
		pus: pus,
		ps:  ps,
		us:  us,
	}
}

type Purchases struct {
	pus models.PurchaseService
	ps  models.ProductsService
	us  models.UserService
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
	product, err := p.ps.ByID(pr.ID)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&messageResponse{
				Message: "Product not found at the store's stock",
			})
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	balance, err := p.us.GetBalance(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if float64(product.Price) > balance {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&messageResponse{
			Message: "Balance is not enough to execute the payment",
		})
		return
	}
	err = p.us.ExecutePayment(user, product.Price)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	purchase := &models.Purchase{
		Email: user.Email,
		ItemP: models.PurchaseItem{
			ID:    product.ID,
			NameP: product.Name,
			Price: product.Price,
		},
	}
	err = p.pus.Create(purchase)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Println("Purchase created")
}

// Get fetchs the purchases for a specific user
//
// GET /purchases
func (p *Purchases) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := context.User(r.Context())
	if user == nil {
		log.Println("Error while fetching the user from the context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	purchases, err := p.pus.ByEmail(user.Email)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Println("Purchases fetched")
	json.NewEncoder(w).Encode(purchases)
}

type createPurchaseRequest struct {
	ID string `json:"id"`
}
