package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/ripa1993/stellarquest/utils"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"log"
)

func main() {
	kp, _ := keypair.Parse("SA6Y4RMN25AHNDAXBGDA3FG7CKIIS2CU4YWTSFZVCENAJB34P3ESSQIT")
	preimage := "KanayeNet"
	preimageHash := sha256.Sum256([]byte(preimage))
	//encode preimageHash as a stellar HashX signer key
	hashx, err := strkey.Encode(strkey.VersionByteHashX, preimageHash[:])

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}

	fmt.Println(kp.Address())

	md := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: hashx,
			Weight:  txnbuild.Threshold(1),
		},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&md},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
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

	///

	md2 := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: hashx,
			Weight:  0,
		},
	}

	tx2, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&md2},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	tx2, err = tx2.SignHashX([]byte(preimage))
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
		utils.ExpandError(err)
		log.Fatalln(err)
	}

	log.Println(res2)
}
