package models

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
)

type StellarService struct {
	client *horizonclient.Client
}

func NewStellarService() *StellarService {
	return &StellarService{
		client: horizonclient.DefaultTestNetClient,
	}
}

// CreateAccount creates an account on the Stellar network (testnet).
// Additinally, the account is filled with 10.000 lumens thanks to friendbot.
func (ss *StellarService) CreateAccount() (*keypair.Full, error) {
	kp, err := keypair.Random()
	if err != nil {
		log.Println("Unable to create a keypair for stellar network")
		return nil, err
	}
	// Create and fund the address on TestNet, using friendbot
	_, err = ss.client.Fund(kp.Address())
	if err != nil {
		log.Println("Unable to fund an account for stellar network")
		return nil, err
	}
	return kp, nil
}
