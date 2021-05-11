package main

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"log"
)

func main() {
	kp, _ := keypair.Parse("SDDZZOXTWA3UK43F6EIAHH4VP3MUYPMYDWX3OAR5LQGX3KXOF4T5QS7W")
	otherKp, _ := keypair.Parse("SB4UUJRW4FHUNGSGMXJZWYJ7DUYZVKLSSSYCUMCTBCINAQ6NG3TWIHSB")

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatalln(err)
	}

	md := txnbuild.SetOptions{
		Signer:               &txnbuild.Signer{
			Address: otherKp.Address(),
			Weight:  1,
		},
		MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(2)),
		SourceAccount:        kp.Address(),
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&md},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	// Get the base 64 encoded transaction envelope
	txe, err := tx.Base64()
	if err != nil {
		log.Fatalln(err)
	}

	// Send the transaction to the network
	res, err := client.SubmitTransactionXDR(txe)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res)


	/////

	md2 := txnbuild.ManageData{
		Name:          "Test",
		Value:         []byte("test"),
		SourceAccount: kp.Address(),
	}

	tx2, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&md2},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tx2, err = tx2.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full), otherKp.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	// Get the base 64 encoded transaction envelope
	txe2, err := tx2.Base64()
	if err != nil {
		log.Fatalln(err)
	}

	// Send the transaction to the network
	res2, err := client.SubmitTransactionXDR(txe2)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res2)
}

