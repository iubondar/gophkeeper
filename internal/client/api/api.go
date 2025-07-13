package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/models"

	"github.com/go-resty/resty/v2"
)

type APIClient struct {
	httpc *resty.Client
}

func NewAPIClient(serverURL string) *APIClient {
	client := resty.New().SetBaseURL(serverURL)
	return &APIClient{httpc: client}
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

	return nil
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
