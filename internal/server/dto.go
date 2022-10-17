package server

import (
	"github.com/ambardhesi/extend-homework/pkg/cards"
	"github.com/ambardhesi/extend-homework/pkg/transactions"
)

type GetVirtualCardsResponse struct {
	VirtualCards []cards.VirtualCard `json:"virtual_cards,omitempty"`
}

type GetVirtualCardTransactionsResponse struct {
	Transactions []transactions.Transaction `json:"transactions,omitempty"`
}

type GetTransactionResponse struct {
	Transaction transactions.Transaction `json:"transaction,omitempty"`
}
