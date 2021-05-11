package main

import (
	"crypto/sha256"
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

	asset := txnbuild.CreditAsset{
		Code:   "RipazSuper",
		Issuer: newKey.Address(),
	}

	md2 := txnbuild.ChangeTrust{
		Line:          asset,
		SourceAccount: kp.Address(),
	}
	md3 := txnbuild.Payment{
		Destination:   kp.Address(),
		Amount:        "1",
		Asset:         asset,
		SourceAccount: newKey.Address(),
	}


	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&md2, &md3},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
			Memo: 				  txnbuild.MemoHash(sha256.Sum256([]byte("Stellar Quest Series 2"))),
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full), newKey.(*keypair.Full))
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
