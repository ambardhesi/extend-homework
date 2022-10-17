package cards

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ambardhesi/extend-homework/pkg/transactions"
	"github.com/go-resty/resty/v2"
)

type CardsService interface {
	GetVirtualCards(accessToken string) ([]VirtualCard, error)
	GetVirtualCardTransactions(accessToken string, cardID string, status string) ([]transactions.Transaction, error)
}

type ExtendCardsService struct {
	client *resty.Client
}

func NewExtendCardService() *ExtendCardsService {
	return &ExtendCardsService{
		client: resty.New(),
	}
}

func (extendSvc *ExtendCardsService) GetVirtualCards(accessToken string) ([]VirtualCard, error) {
	res, err := extendSvc.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/vnd.paywithextend.v2021-03-12+json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Get("https://api.paywithextend.com/virtualcards")

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, errors.New(string(res.Body()))
	}

	body := string(res.Body())

	var resp GetExtendVirtualCardsResponse
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return nil, err
	}

	return resp.VirtualCards, nil
}

func (extendSvc *ExtendCardsService) GetVirtualCardTransactions(accessToken string, cardID string, status string) ([]transactions.Transaction, error) {
	res, err := extendSvc.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/vnd.paywithextend.v2021-03-12+json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetQueryParam("status", status).
		Get(fmt.Sprintf("https://api.paywithextend.com/virtualcards/%s/transactions", cardID))

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, errors.New(string(res.Body()))
	}

	body := string(res.Body())

	var resp GetExtendVirtualCardTransactionsResponse
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return nil, err
	}

	return resp.Transactions, nil
}
