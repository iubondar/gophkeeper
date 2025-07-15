package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/internal/models"
	"strings"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type APIClient struct {
	httpc        *resty.Client
	accessToken  string
	refreshToken string
}

func NewAPIClient(serverURL string) *APIClient {
	if !strings.HasPrefix(serverURL, "http") {
		serverURL = "http://" + serverURL
	}
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

// handleErrorResponse обрабатывает ошибки HTTP ответа и разбирает JSONError
func (c *APIClient) handleErrorResponse(response *resty.Response) error {
	if response.StatusCode() >= 400 {
		// Пытаемся разобрать JSONError
		jsonErr, err := models.ParseJSONError(response.Body())
		if err == nil {
			return errors.New(jsonErr.Message)
		} else {
			zap.L().Sugar().Debugln("Failed to parse JSON error", zap.Error(err))
		}
		// Если не удалось разобрать JSONError, возвращаем обычную ошибку
		return errors.New(response.String())
	}
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

	if err := c.handleErrorResponse(response); err != nil {
		return err
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

	if err := c.handleErrorResponse(response); err != nil {
		return "", err
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

	if err := c.handleErrorResponse(response); err != nil {
		return err
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

// HealthCheck проверяет доступность сервера по /api/health
func (c *APIClient) HealthCheck(ctx context.Context) error {
	response, err := c.httpc.R().
		SetContext(ctx).
		Get("/health")

	if err != nil {
		return err
	}

	if err := c.handleErrorResponse(response); err != nil {
		return err
	}

	return nil
}
