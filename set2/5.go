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

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}
	log.Println(sourceAccount.AccountID)

	md := txnbuild.ClaimClaimableBalance{
		BalanceID: "0000000055a88594744ce9355825244f1476dbd6af0e4780171e3b7653e473891d4d6b05",
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
}
