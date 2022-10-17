package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdGetTransaction = cobra.Command{
		Use:   "transaction txID",
		Short: "Get a transaction by its ID",
		RunE:  getTransaction,
		Args:  cobra.MinimumNArgs(1),
	}
)

func getTransaction(cobraCmd *cobra.Command, args []string) error {
	txID := args[0]
	resp, err := makeClient().GetTransaction(txID)

	if err != nil {
		return err 
	}

	var out bytes.Buffer
	json.Indent(&out, []byte(*resp), "", "\t")
	fmt.Println(out.String())

	return nil
}
