package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdSignIn = cobra.Command{
		Use:   "signin email password",
		Short: "Sign in with Extend",
		RunE:  signIn,
		Args:  cobra.MinimumNArgs(2),
	}
)

func signIn(cobraCmd *cobra.Command, args []string) error {
	emailID := args[0]
	password := args[1]

	resp, err := makeClient().SignIn(emailID, password)

	if err != nil {
		return err 
	}

	var out bytes.Buffer
	json.Indent(&out, []byte(*resp), "", "\t")
	fmt.Println(out.String())

	return nil
}
