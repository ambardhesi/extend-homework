package transactions

import (
	"github.com/ambardhesi/extend-homework/pkg/customtime"
)

type Transaction struct {
	Id                     string          `json:"id,omitempty"`
	CardholderName         string          `json:"cardholderName,omitempty"`
	RecipientName          string          `json:"recipientName,omitempty"`
	Status                 string          `json:"status,omitempty"`
	AuthBillingAmountCents int             `json:"authBillingAmountCents,omitempty"`
	MerchantName           string          `json:"merchantName,omitempty"`
	AuthedAt               customtime.Time `json:"authedAt,omitempty"`
}
