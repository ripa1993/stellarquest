package utils

import (
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
)

func ExpandError(err error) {
	if err2, ok := err.(*horizonclient.Error); ok {
		fmt.Println("Error has additional info")
		fmt.Println(err2.ResultCodes())
		fmt.Println(err2.ResultString())
		fmt.Println(err2.Problem)
	}
}
