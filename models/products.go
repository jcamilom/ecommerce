package models

import (
	"github.com/jcamilom/ecommerce/db"
)

var (
	// The DB table name for products
	dbProductsTableName = "Products"

	// The DB primary key for products
	dbProductsKeyName = "id"
)

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

// ProductDB is used to interact with the products database.
type ProductDB interface {
	// Methods for querying for single products
	ByID(id string) (*Product, error)
}

// ProductsService is a set of methods used to manipulate and
// work with the product model
type ProductsService interface {
	ProductDB
}

func NewProductsService() ProductsService {
	pdb := newProductDB()
	return &productsService{
		ProductDB: pdb,
	}
}

var _ ProductsService = &productsService{}

type productsService struct {
	ProductDB
}

var _ ProductDB = &productDB{}

func newProductDB() *productDB {
	db := &db.DB{}
	return &productDB{
		db: db,
	}
}

type productDB struct {
	db *db.DB
}

// ByID will look up a product with the provided ID.
func (pdb *productDB) ByID(id string) (*Product, error) {
	p := new(Product)
	key := struct {
		ID string `json:"id"`
	}{
		ID: id,
	}
	found, err := pdb.db.GetItem(key, dbProductsTableName, p)
	if err != nil {
		return nil, err
	} else if found == false {
		return nil, ErrNotFound
	} else {
		return p, nil
	}
}
