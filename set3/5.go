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
	kp, _ := keypair.Parse("SC5WO2DCAEHPAR5UBGVEVBJJ24IRRI7UNC2YDM5AB4ZQ3UBJWSLVKGUG")
	kp2, _ := keypair.Parse("SDJSML5XKZYLPBDHTCJ43UU4GIWYNTLEAQBB7FIVVUVMOHEIP4EILY7K")

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}

	fmt.Println(kp.Address())

	asset := txnbuild.CreditAsset{
		Code:   "MariusLenk",
		Issuer: kp.Address(),
	}

	changeTrust := txnbuild.ChangeTrust{
		Line:   asset,
		Limit: "1000",
		SourceAccount: kp2.Address(),
	}

	setOptions := txnbuild.SetOptions{
		SetFlags:             []txnbuild.AccountFlag{txnbuild.AuthClawbackEnabled, txnbuild.AuthRevocable},
		ClearFlags: []txnbuild.AccountFlag{txnbuild.AuthRequired},
	}

	payment := txnbuild.Payment{
		Destination:   kp2.Address(),
		Amount:        "1",
		Asset:         asset,
	}

	revoke := txnbuild.Clawback{
		From:          kp2.Address(),
		Amount:        "1",
		Asset:         &asset,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&setOptions, &changeTrust, &payment, &revoke},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds: 		  txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full), kp2.(*keypair.Full))
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
