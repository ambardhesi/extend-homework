package transactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type TransactionsService interface {
	GetTransaction(accessToken string, transactionID string) (Transaction, error)
}

type ExtendTransactionsService struct {
	client *resty.Client
}

func NewExtendTransactionsService() *ExtendTransactionsService {
	return &ExtendTransactionsService{
		client: resty.New(),
	}
}

func (txSvc *ExtendTransactionsService) GetTransaction(accessToken string, txID string) (Transaction, error) {
	res, err := txSvc.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/vnd.paywithextend.v2021-03-12+json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Get(fmt.Sprintf("https://api.paywithextend.com/transactions/%s", txID))

	var transaction Transaction
	if err != nil {
		return transaction, err
	}

	if res.StatusCode() != http.StatusOK {
		return transaction, errors.New(string(res.Body()))
	}

	body := string(res.Body())

	err = json.Unmarshal([]byte(body), &transaction)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
