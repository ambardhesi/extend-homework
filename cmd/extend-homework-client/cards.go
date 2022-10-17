package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdGetVirtualCards = cobra.Command{
		Use:   "cards",
		Short: "Gets all the virtual cards for a user",
		RunE:  getVirtualCards,
	}
)

var (
	cmdGetVirtualCardTransactions = cobra.Command{
		Use:   "card-transactions cardID status",
		Short: "Gets all the transactions for a card",
		RunE:  getVirtualCardTransactions,
		Args:  cobra.MinimumNArgs(2),
	}
)

func getVirtualCards(combraCmd *cobra.Command, args []string) error {
	resp, err := makeClient().GetVirtualCards()

	if err != nil {
		return err
	}

	var out bytes.Buffer
	json.Indent(&out, []byte(*resp), "", "\t")
	fmt.Println(out.String())

	return nil
}

func getVirtualCardTransactions(cobraCmd *cobra.Command, args []string) error {
	cardID := args[0]
	status := args[1]

	resp, err := makeClient().GetVirtualCardTransactions(cardID, status)

	if err != nil {
		return err
	}

	var out bytes.Buffer
	json.Indent(&out, []byte(*resp), "", "\t")
	fmt.Println(out.String())

	return nil
}
