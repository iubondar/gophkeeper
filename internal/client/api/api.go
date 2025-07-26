package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/internal/auth"
	"gophkeeper/internal/models"
	"net/http"
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

// setAuthCookie устанавливает cookie аутентификации для запроса
func (c *APIClient) setAuthCookie(request *resty.Request) error {
	if c.accessToken == "" {
		return errors.New("access token is required")
	}
	request.SetCookie(&http.Cookie{
		Name:  auth.AuthCookieName,
		Value: c.accessToken,
	})
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
func (c *APIClient) UploadSecret(ctx context.Context, in models.UploadSecretIn) error {
	request := c.httpc.R().
		SetContext(ctx).
		SetBody(in)

	if err := c.setAuthCookie(request); err != nil {
		return err
	}

	response, err := request.Post("/api/upload")
	if err != nil {
		return err
	}

	if err := c.handleErrorResponse(response); err != nil {
		return err
	}

	return nil
}

// GetSecretVersion получает версию секрета с сервера
func (c *APIClient) GetSecretVersion(ctx context.Context, secretName string) (int, error) {
	request := c.httpc.R().
		SetContext(ctx).
		SetQueryParam("name", secretName)

	if err := c.setAuthCookie(request); err != nil {
		return 0, err
	}

	response, err := request.Get("/api/version")
	if err != nil {
		return 0, err
	}

	if err := c.handleErrorResponse(response); err != nil {
		return 0, err
	}

	var versionResult struct {
		Version int `json:"version"`
	}
	err = json.Unmarshal(response.Body(), &versionResult)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal version response: %w", err)
	}

	return versionResult.Version, nil
}

// UpdateSecret обновляет секрет на сервере
func (c *APIClient) UpdateSecret(ctx context.Context, secret models.SecretData, version int) error {
	// Подготавливаем данные для обновления
	updateData := models.UpdateSecretIn{
		Label:         secret.Name,
		Type:          secret.Type,
		Metadata:      secret.Metadata,
		EncryptedData: []byte(secret.Data),
		FileKey:       "", // TODO: добавить поддержку файлов
		Version:       version,
	}

	// Выполняем обновление
	request := c.httpc.R().
		SetContext(ctx).
		SetBody(updateData)

	if err := c.setAuthCookie(request); err != nil {
		return err
	}

	response, err := request.Put("/api/update")
	if err != nil {
		return err
	}

	if err := c.handleErrorResponse(response); err != nil {
		return err
	}

	return nil
}

// GetSecret получает секрет с сервера
func (c *APIClient) GetSecret(ctx context.Context, secretName string) (*models.GetSecretOut, error) {
	request := c.httpc.R().
		SetContext(ctx).
		SetQueryParam("name", secretName)

	if err := c.setAuthCookie(request); err != nil {
		return nil, err
	}

	response, err := request.Get("/api/get")
	if err != nil {
		return nil, err
	}

	if err := c.handleErrorResponse(response); err != nil {
		return nil, err
	}

	var out models.GetSecretOut
	err = json.Unmarshal(response.Body(), &out)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal get secret response: %w", err)
	}

	return &out, nil
}

// DeleteSecret удаляет секрет с сервера
func (c *APIClient) DeleteSecret(ctx context.Context, secretName string) error {
	request := c.httpc.R().
		SetContext(ctx).
		SetQueryParam("name", secretName)

	if err := c.setAuthCookie(request); err != nil {
		return err
	}

	response, err := request.Delete("/api/delete")
	if err != nil {
		return err
	}

	if err := c.handleErrorResponse(response); err != nil {
		return err
	}

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
