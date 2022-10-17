package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ambardhesi/extend-homework/pkg/apperror"
	"github.com/go-resty/resty/v2"
)

type AuthService interface {
	GetAccessToken(emailID string) (string, error)
	SignIn(email string, password string) error
}

type InMemoryDbAuthService struct {
	tokenByEmailId        map[string]AccessToken
	refreshTokenByEmailId map[string]string
	client                *resty.Client
}

type AccessToken struct {
	Token     string
	ValidTill time.Time
}

func NewInMemoryDbAuthService() *InMemoryDbAuthService {
	return &InMemoryDbAuthService{
		tokenByEmailId: make(map[string]AccessToken),
		refreshTokenByEmailId: make(map[string]string),
		client:         resty.New(),
	}
}

func (authSvc *InMemoryDbAuthService) GetAccessToken(emailID string) (string, error) {
	token, ok := authSvc.tokenByEmailId[emailID]
	if !ok {
		return "", &apperror.Error{
			Code: apperror.ENOTFOUND,
		}
	}

	// Token has expired, need to refresh
	if time.Now().After(token.ValidTill) {
		refreshToken, ok := authSvc.refreshTokenByEmailId[emailID]
		if !ok {
			return "", &apperror.Error{
				Code: apperror.ENOTFOUND,
			}
		}

		req := RenewRequest{
			RefreshToken: refreshToken,
		}

		res, err := authSvc.client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/vnd.paywithextend.v2021-03-12+json").
			SetBody(req).
			Post("https://api.paywithextend.com/renewauth")

		if err != nil {
			return "", err
		}

		if res.StatusCode() != http.StatusOK {
			return "", errors.New(string(res.Body()))
		}

		body := string(res.Body())

		var renewResponse RenewResponse
		err = json.Unmarshal([]byte(body), &renewResponse)
		if err != nil {
			return "", err
		}

		authSvc.storeTokens(emailID, renewResponse.Token, renewResponse.RefreshToken)

		return renewResponse.Token, nil

	}

	return token.Token, nil
}

func (authSvc *InMemoryDbAuthService) SignIn(emailID string, password string) error {
	req := SignInRequest{
		Email:    emailID,
		Password: password,
	}

	res, err := authSvc.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/vnd.paywithextend.v2021-03-12+json").
		SetBody(req).
		Post("https://api.paywithextend.com/signin")

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New(string(res.Body()))
	}

	body := string(res.Body())

	var signInResponse SignInResponse
	err = json.Unmarshal([]byte(body), &signInResponse)
	if err != nil {
		return err
	}

	authSvc.storeTokens(emailID, signInResponse.Token, signInResponse.RefreshToken)

	return nil
}

func (authSvc *InMemoryDbAuthService) storeTokens(emailID string, token string, refreshToken string) {
	authSvc.tokenByEmailId[emailID] = AccessToken{
		Token:     token,
		ValidTill: time.Now().Add(time.Minute * 10),
	}

	authSvc.refreshTokenByEmailId[emailID] = refreshToken
}
