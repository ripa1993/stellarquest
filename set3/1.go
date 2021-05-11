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
	kp, _ := keypair.Parse("SC5OXINMTJVALOEJMTF43JX7PBLPPKMOHXA6ELPT4G3CCZJI5XGGB255")

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		expandError(err)
		log.Fatalln(err)
	}

	fmt.Println(kp.Address())

	bumpTo := int64(110101115104111) // "nesho" as a buffer

	md := txnbuild.BumpSequence{
		BumpTo:        bumpTo,
	}


	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&md},
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
		expandError(err)
		log.Fatalln(err)
	}

	log.Println(res)
}

func expandError(err error)  {
	if err2, ok := err.(*horizonclient.Error); ok {
		fmt.Println("Error has additional info")
		fmt.Println(err2.ResultCodes())
		fmt.Println(err2.ResultString())
		fmt.Println(err2.Problem)
	}
}