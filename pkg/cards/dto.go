package cards

import (
	"github.com/ambardhesi/extend-homework/pkg/customtime"
	"github.com/ambardhesi/extend-homework/pkg/transactions"
)

type VirtualCard struct {
	Id           string    `json:"id,omitempty"`
	Status       string    `json:"status,omitempty"`
	BalanceCents int    `json:"balanceCents,omitempty"`
	ValidFrom    customtime.Time `json:"validFrom,omitempty"`
	ValidTo      customtime.Time `json:"validTo,omitempty"`
	CreatedAt    customtime.Time `json:"createdAt,omitempty"`
	Issuer       Issuer    `json:"issuer,omitempty"`
}

type Issuer struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type GetExtendVirtualCardsResponse struct {
	VirtualCards []VirtualCard `json:"virtualCards,omitempty"`
}

type GetExtendVirtualCardTransactionsResponse struct {
	Transactions []transactions.Transaction `json:"transactions,omitempty"`
}
