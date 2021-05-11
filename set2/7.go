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

	other := "GDMUY6V2PJV35VN6FYBO3HN6SUSUCFM7J6VUYMYF2URFGKP3M4RUXI6O"

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}
	log.Println(sourceAccount.AccountID)


	pay := txnbuild.Payment{
		Destination:   other,
		Amount:        "1000",
		Asset:         txnbuild.NativeAsset{},
		SourceAccount: kp.Address(),
	}
	md := txnbuild.RevokeSponsorship{
		SourceAccount:    kp.Address(),
		Account:          &other,
		SponsorshipType:  txnbuild.RevokeSponsorshipTypeAccount,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&pay, &md},
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
		utils.ExpandError(err)
		log.Fatalln(err)
	}

	log.Println(res)
}
