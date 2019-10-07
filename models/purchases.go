package models

import (
	"time"

	"github.com/jcamilom/ecommerce/db"
)

var (
	// The DB table name for products
	dbPurchaseTableName = "Purchases"

	// The DB primary key for products
	dbPurchaseKeyName = "id"
)

type Purchase struct {
	ID    string    `json:"id"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
	Item  PurchaseItem
}

type PurchaseItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// PurchaseDB is used to interact with the purchases database.
type PurchaseDB interface {
	// Methods for querying for single purchase
	// ByID(id string) (*Purchase, error)
	// Methods for altering purchases
	Create(purchase *Purchase) error
}

// PurchaseService is a set of methods used to manipulate and
// work with the purchase model
type PurchaseService interface {
	PurchaseDB
}

func NewPurchaseService() PurchaseService {
	pdb := newPurchaseDB()
	pv := newPurchaseValidator(pdb)
	return &purchaseService{
		PurchaseDB: pv,
	}
}

var _ PurchaseService = &purchaseService{}

type purchaseService struct {
	PurchaseDB
}

type purchaseValFunc func(*Purchase) error

func runPurchaseValFuncs(purchase *Purchase, fns ...purchaseValFunc) error {
	for _, fn := range fns {
		if err := fn(purchase); err != nil {
			return err
		}
	}
	return nil
}

var _ PurchaseDB = &purchaseValidator{}

func newPurchaseValidator(pdb PurchaseDB) *purchaseValidator {
	return &purchaseValidator{
		PurchaseDB: pdb,
	}
}

type purchaseValidator struct {
	PurchaseDB
}

// Create will fill necessary data for the purchase
func (pv *purchaseValidator) Create(purchase *Purchase) error {
	err := runPurchaseValFuncs(purchase, pv.setCreationTime)
	if err != nil {
		return err
	}
	return pv.PurchaseDB.Create(purchase)
}

func (pv *purchaseValidator) setCreationTime(purchase *Purchase) error {
	purchase.Date = time.Now()
	return nil
}

var _ PurchaseDB = &purchaseDB{}

func newPurchaseDB() *purchaseDB {
	db := &db.DB{}
	return &purchaseDB{
		db: db,
	}
}

type purchaseDB struct {
	db *db.DB
}

func (pdb *purchaseDB) Create(purchase *Purchase) error {
	return pdb.db.PutItem(dbPurchaseTableName, purchase)
}
