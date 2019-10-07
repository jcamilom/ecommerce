package models

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jcamilom/ecommerce/db"
	"github.com/mitchellh/hashstructure"
)

var (
	// The DB table name for purchases
	dbPurchaseTableName = "Purchases"

	// The DB partition key for purchases
	dbPurchasePartitionKeyName = "email"

	// The DB sort key for purchases
	dbPurchaseSortKeyName = "id"
)

type Purchase struct {
	ID    string       `json:"id"`
	Email string       `json:"email"`
	Date  time.Time    `json:"date"`
	ItemP PurchaseItem `json:"item_p"`
}

type PurchaseItem struct {
	ID    string `json:"id"`
	NameP string `json:"name_p"`
	Price int    `json:"price"`
}

// PurchaseDB is used to interact with the purchases database.
type PurchaseDB interface {
	// Methods for querying several purchases
	ByEmail(email string) ([]Purchase, error)
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
	err := runPurchaseValFuncs(purchase,
		pv.setCreationTime,
		pv.setID,
	)
	if err != nil {
		return err
	}
	return pv.PurchaseDB.Create(purchase)
}

func (pv *purchaseValidator) setCreationTime(purchase *Purchase) error {
	purchase.Date = time.Now()
	return nil
}

func (pv *purchaseValidator) setID(purchase *Purchase) error {
	hash, err := hashstructure.Hash(purchase, nil)
	if err != nil {
		return err
	}

	hashStr := strconv.FormatUint(hash, 10)
	purchase.ID = hashStr
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

func (pdb *purchaseDB) ByEmail(email string) ([]Purchase, error) {
	purchases := []Purchase{}
	key := struct {
		Email string `json:":e"`
	}{
		Email: email,
	}
	keyCondExp := "email = :e"
	projectionExp := "id, email, item_p.id, item_p.price, item_p.name_p, #dt"
	expressionAttributeNames := map[string]*string{
		"#dt": aws.String("date"),
	}
	err := pdb.db.GetItems(dbPurchaseTableName, key, keyCondExp, projectionExp, expressionAttributeNames, &purchases)
	if err != nil {
		return nil, err
	}
	return purchases, nil
}
