package main

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"log"
)

const PublicKey3 = "GA64TIIF7BWISGMCV5IIMSPA4OUHQOHRLMGTJIICLCL7QYPDTEMISCGK"
const SecretKey3 = "SDDZZOXTWA3UK43F6EIAHH4VP3MUYPMYDWX3OAR5LQGX3KXOF4T5QS7W"

func main() {
	kp, _ := keypair.Parse(SecretKey3)

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatalln(err)
	}

	md := txnbuild.ManageData{
		Name:          "Hello",
		Value:         []byte("World"),
		SourceAccount: PublicKey3,
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
