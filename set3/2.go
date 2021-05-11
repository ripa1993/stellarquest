package main

import (
	"fmt"
	"github.com/ripa1993/stellarquest/utils"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"log"
)

func main() {
	kp, _ := keypair.Parse("SCWMLELHE2KOMH7CMSGVVQRSSIMOYAAOR34RQWNZUCIAIFEPNYOGFYBM")

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}

	fmt.Println(kp.Address())


	operations := make([]txnbuild.Operation, 0, 100)
	i := 0
	for i<100 {
		operations = append(operations, &txnbuild.ManageData{
			Name:          fmt.Sprintf("Test%d", i),
			Value:         []byte("test"),
		})
		i++
	}



	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           operations,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds: 		  txnbuild.NewInfiniteTimeout(),
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
		utils.ExpandError(err)
		log.Fatalln(err)
	}

	log.Println(res)
}
