package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ambardhesi/extend-homework/pkg/auth"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	ServerAddress  string
	CertFilePath   string
	KeyFilePath    string
	CaCertFilePath string
}

type Client struct {
	Config     Config
	HttpClient *resty.Client
}

func NewClient(config Config) (*Client, error) {
	rClient := resty.New()
	tlsConfig, err := GetTLSConfig(config.CertFilePath,
		config.KeyFilePath, config.CaCertFilePath)

	if err != nil {
		fmt.Printf("Failed to create http client %v\n", err)
		return nil, err
	}
	rClient.SetTLSClientConfig(tlsConfig)

	return &Client{
		HttpClient: rClient,
		Config:     config,
	}, nil
}

func (c *Client) SignIn(emailID string, password string) (*string, error) {
	req := auth.SignInRequest{
		Email:    emailID,
		Password: password,
	}

	resp, err := c.HttpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(c.Config.ServerAddress + "/signin")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	body := string(resp.Body())
	return &body, nil
}

func (c *Client) GetVirtualCards() (*string, error) {
	resp, err := c.HttpClient.R().
		SetHeader("Content-Type", "application/json").
		Get(c.Config.ServerAddress + "/cards")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	body := string(resp.Body())
	return &body, nil
}

func (c *Client) GetVirtualCardTransactions(cardID string, status string) (*string, error) {
	resp, err := c.HttpClient.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("status", status).
		Get(c.Config.ServerAddress + "/cards/" + cardID + "/transactions")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	body := string(resp.Body())
	return &body, nil
}

func (c *Client) GetTransaction(txID string) (*string, error) {
	resp, err := c.HttpClient.R().
		SetHeader("Content-Type", "application/json").
		Get(c.Config.ServerAddress + "/transactions/" + txID)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	body := string(resp.Body())
	return &body, nil
}
