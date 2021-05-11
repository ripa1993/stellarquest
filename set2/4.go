package main

import (
	"github.com/ripa1993/stellarquest/utils"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
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

	before10secs, err := xdr.NewClaimPredicate(xdr.ClaimPredicateTypeClaimPredicateBeforeRelativeTime, xdr.Int64(10))
	if err != nil {
		log.Println("time")
		log.Fatalln(err)
		return
	}
	notBefore10secs, err := xdr.NewClaimPredicate(xdr.ClaimPredicateTypeClaimPredicateNot, &before10secs)
	if err != nil {
		log.Println("not")
		log.Fatalln(err)
		return
	}

	md := txnbuild.CreateClaimableBalance{
		Amount:        "100",
		Asset:         txnbuild.NativeAsset{},
		Destinations:  []txnbuild.Claimant{
			txnbuild.NewClaimant(kp.Address(), &notBefore10secs),
		},
		SourceAccount: kp.Address(),
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
