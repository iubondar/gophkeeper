package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/models"

	"github.com/go-resty/resty/v2"
)

type APIClient struct {
	httpc        *resty.Client
	accessToken  string
	refreshToken string
}

func NewAPIClient(serverURL string) *APIClient {
	client := resty.New().SetBaseURL(serverURL)
	return &APIClient{httpc: client}
}

// handleAuthenticateResponse обрабатывает ответ аутентификации и обновляет токены в клиенте
func (c *APIClient) handleAuthenticateResponse(responseBody []byte) error {
	var out models.AuthenticateOut
	err := json.Unmarshal(responseBody, &out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal authenticate response: %w", err)
	}

	c.accessToken = out.AccessToken
	c.refreshToken = out.RefreshToken

	return nil
}

func (c *APIClient) Register(ctx context.Context, in models.RegisterIn) error {
	response, err := c.httpc.R().
		SetContext(ctx).
		SetBody(in).
		Post("/api/register")

	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("failed to register: %s", response.String())
	}

	return c.handleAuthenticateResponse(response.Body())
}

// Login выполняет запрос на вход в систему
func (c *APIClient) Login(ctx context.Context, in models.LoginIn) (salt string, err error) {
	response, err := c.httpc.R().
		SetContext(ctx).
		SetBody(in).
		Post("/api/login")

	if err != nil {
		return "", err
	}

	if response.IsError() {
		return "", fmt.Errorf("failed to login: %s", response.String())
	}

	var out models.LoginOut
	err = json.Unmarshal(response.Body(), &out)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal login response: %w", err)
	}

	return out.Salt, nil
}

func (c *APIClient) Authenticate(ctx context.Context, in models.AuthenticateIn) error {
	response, err := c.httpc.R().
		SetContext(ctx).
		SetBody(in).
		Post("/api/authenticate")

	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("failed to authenticate: %s", response.String())
	}

	return c.handleAuthenticateResponse(response.Body())
}

// UploadSecret загружает секрет на сервер
func (c *APIClient) UploadSecret(ctx context.Context, secret models.SecretData) error {
	// TODO: Implement

	return nil
}

// UpdateSecret обновляет секрет на сервере
func (c *APIClient) UpdateSecret(ctx context.Context, secret models.SecretData) error {
	// TODO: Implement

	return nil
}

// GetSecret получает секрет с сервера
func (c *APIClient) GetSecret(ctx context.Context, secretName string) (*models.SecretData, error) {
	// TODO: Implement

	return nil, nil
}

// DeleteSecret удаляет секрет с сервера
func (c *APIClient) DeleteSecret(ctx context.Context, secretName string) error {
	// TODO: Implement

	return nil
}
