package main

import (
	"github.com/ripa1993/stellarquest/utils"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"log"
)

func main() {
	kp, _ := keypair.Parse("SDPT5WKMILO2Z6EPNA2APFVIGEOGYQWOCPE4TEYK6QNW5G7RZCA3RJRR")

	newKey, _ := keypair.Parse("SCHKGEAIOSQUXFJZHCSWBKEOXP4P5ZMCMMR3KW2YEUT5SBHGU4X3HCM4")

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}
	log.Println(sourceAccount.AccountID)

	otherAccount, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: newKey.Address()})
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}
	log.Println(otherAccount.AccountID)

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{&txnbuild.ManageData{
				Name:  "Test",
				Value: []byte("Test"),
			}},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimeout(300),
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

	tx2, err := txnbuild.NewFeeBumpTransaction(
		txnbuild.FeeBumpTransactionParams{
			Inner:               tx,
			FeeAccount:          newKey.Address(),
			BaseFee:             txnbuild.MinBaseFee + 1,
			EnableMuxedAccounts: false,
		},
	)

	tx2, err = tx2.Sign(network.TestNetworkPassphrase, newKey.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	// Get the base 64 encoded transaction envelope
	txe, err := tx2.Base64()
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
