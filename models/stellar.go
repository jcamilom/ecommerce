package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

// StellarService performs all the operation in the stellar network
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
	log.Println("Account created on stellar")
	log.Println("Seed:", kp.Seed())
	log.Println("Address:", kp.Address())
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

// ExecutePayment performs a payment operation in the stellar network
func (ss *StellarService) ExecutePayment(sourceSeed, destinationAddr, amount string) error {
	// Recover the keypair from the account seed
	kp, _ := keypair.Parse(sourceSeed)
	// Get information about the account
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := ss.client.AccountDetail(ar)
	if err != nil {
		log.Println("Unable to fetch account details")
		return err
	}

	// Construct the operation
	paymentOp := txnbuild.Payment{
		Destination: destinationAddr,
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}

	// Construct the transaction that will carry the operation
	tx := txnbuild.Transaction{
		SourceAccount: &sourceAccount,
		Operations:    []txnbuild.Operation{&paymentOp},
		Timebounds:    txnbuild.NewInfiniteTimeout(),
		Network:       network.TestNetworkPassphrase,
	}

	// Sign the transaction, serialise it to XDR, and base 64 encode it
	_, err = tx.BuildSignEncode(kp.(*keypair.Full))
	if err != nil {
		log.Println("Unable to encode the transaction")
		return err
	}

	// Submit the transaction
	_, err = ss.client.SubmitTransaction(tx)
	if err != nil {
		log.Println("Unable to submit the transaction")
		switch e := err.(type) {
		case *horizonclient.Error:
			fmt.Println("err type=" + e.Problem.Type)
			fmt.Println("err detailed=" + e.Problem.Detail)
			fmt.Print("err extras=")
			fmt.Printf("%+v\n", e.Problem.Extras["result_codes"])
		}
		return err
	}
	return nil
}
