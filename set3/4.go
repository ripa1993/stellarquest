package main

import (
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
	kp, _ := keypair.Parse("SDLT3H5C2FKJAWSIT7RF4F7X75BXRRAMJXAQO42HTTH4WM4Z5MPF3SS2")

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}


	fmt.Println(kp.Address())

	sn, _ := sourceAccount.GetSequenceNumber()

	currentAccount := txnbuild.NewSimpleAccount(kp.Address(), sn + 1)
	futureAccount := txnbuild.NewSimpleAccount(kp.Address(), sn + 2)

	manageData := txnbuild.ManageData{
		Name:          "Test",
		Value:         []byte("test"),
	}

	txFuture, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &futureAccount,
			IncrementSequenceNum: false,
			Operations:           []txnbuild.Operation{&manageData},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds: 		  txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	txFutureHash, err := txFuture.Hash(network.TestNetworkPassphrase)

	txFuture, err = txFuture.Sign(network.TestNetworkPassphrase)
	if err != nil {
		log.Fatalln(err)
	}

	// Get the base 64 encoded transaction envelope
	txFutureEnc, err := txFuture.Base64()
	if err != nil {
		log.Fatalln(err)
	}


	///



	preAuth, err := strkey.Encode(strkey.VersionByteHashTx, txFutureHash[:])

	setOptions := txnbuild.SetOptions{
		Signer:               &txnbuild.Signer{
			Address: preAuth,
			Weight:  1,
		},
	}

	txNow, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &currentAccount,
			IncrementSequenceNum: false,
			Operations:           []txnbuild.Operation{&setOptions},
			BaseFee:              500,
			Timebounds: 		  txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	txNow, err = txNow.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	// Get the base 64 encoded transaction envelope
	txNowEnc, err := txNow.Base64()
	if err != nil {
		log.Fatalln(err)
	}

	// Send the transaction to the network
	res2, err := client.SubmitTransactionXDR(txNowEnc)
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}

	log.Println(res2)

	fmt.Printf("Submitted TxNow with SN %d\n", txNow.SourceAccount().Sequence)

	sourceAccount2, err := client.AccountDetail(ar)
	fmt.Printf("Current SN %s\n", sourceAccount2.Sequence)

	fmt.Printf("\n%v\n", txFuture)

	fmt.Printf("Going to submit TxFuture with SN %d\n", txFuture.SourceAccount().Sequence)
	fmt.Println(txFutureEnc)

	// Send the transaction to the network
	res, err := client.SubmitTransactionXDR(txFutureEnc)
	if err != nil {
		utils.ExpandError(err)
		log.Fatalln(err)
	}

	log.Println(res)



}
