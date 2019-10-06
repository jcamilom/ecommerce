package models

import (
	"errors"
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

// GetBalance gets the balance of the provided address on the Stellar network
func (ss *StellarService) GetBalance(address string) (string, error) {
	accountRequest := horizonclient.AccountRequest{AccountID: address}
	hAccount0, err := ss.client.AccountDetail(accountRequest)
	if err != nil {
		return "", err
	}
	for _, balance := range hAccount0.Balances {
		if balance.Type == "native" {
			return balance.Balance, nil
		}
	}
	return "", errors.New("Can't find native balance")
}
