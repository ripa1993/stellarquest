package main

import (
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"log"
)

func main() {
	kp, _ := keypair.Parse("SDDZZOXTWA3UK43F6EIAHH4VP3MUYPMYDWX3OAR5LQGX3KXOF4T5QS7W")
	issuerKp, _ := keypair.Parse("SDZO6S6V64GO7AQKM2QRSXG32SWTCTMJ4SUTGZT4HAS4RXX2POSONNCN")


	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(sourceAccount.AccountID)

	asset := txnbuild.CreditAsset{
		Code:   "RipazCoins",
		Issuer: issuerKp.Address(),
	}

	md := txnbuild.ChangeTrust{
		Line: asset,
		Limit: "1000",
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&md},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(100),
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
		if err2, ok := err.(*horizonclient.Error); ok {
			fmt.Println("Error has additional info")
			fmt.Println(err2.ResultCodes())
			fmt.Println(err2.ResultString())
			fmt.Println(err2.Problem)
		}
		return
		log.Fatalln(err)
	}

	log.Println(res)


	// send

	ar2 := horizonclient.AccountRequest{AccountID: issuerKp.Address()}
	sourceAccount2, err := client.AccountDetail(ar2)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(sourceAccount2.AccountID)


	md2 := txnbuild.Payment{
		Destination:   kp.Address(),
		Amount:        "10",
		Asset:         asset,
	}

	tx2, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount2,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&md2},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(100),
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tx2, err = tx2.Sign(network.TestNetworkPassphrase, issuerKp.(*keypair.Full))
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
		if err2, ok := err.(*horizonclient.Error); ok {
			fmt.Println("Error has additional info")
			fmt.Println(err2.ResultCodes())
			fmt.Println(err2.ResultString())
			fmt.Println(err2.Problem)
		}
		return
		log.Fatalln(err)
	}

	log.Println(res2)
}

