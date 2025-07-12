package api

import (
	"context"
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
