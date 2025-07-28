package cmd

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// MockLoginAPIClient представляет мок для LoginAPIClient
type MockLoginAPIClient struct {
	loginFunc        func(ctx context.Context, in models.LoginIn) (salt string, err error)
	authenticateFunc func(ctx context.Context, in models.AuthenticateIn) error
}

func (m *MockLoginAPIClient) Login(ctx context.Context, in models.LoginIn) (salt string, err error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, in)
	}
	return "", nil
}

func (m *MockLoginAPIClient) Authenticate(ctx context.Context, in models.AuthenticateIn) error {
	if m.authenticateFunc != nil {
		return m.authenticateFunc(ctx, in)
	}
	return nil
}

func TestNewLoginCommand(t *testing.T) {
	mockAPI := &MockLoginAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewLoginCommand(mockAPI, crypto)

	if cmd == nil {
		t.Fatal("NewLoginCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != crypto {
		t.Error("crypto not set correctly")
	}
}

func TestLoginCommand_Execute(t *testing.T) {
	ctx := context.Background()
	credentials := models.UserCredentials{
		Login:    "testuser",
		Password: "testpass",
	}

	t.Run("successful login", func(t *testing.T) {
		mockAPI := &MockLoginAPIClient{
			loginFunc: func(ctx context.Context, in models.LoginIn) (salt string, err error) {
				if in.Login != "testuser" {
					t.Errorf("Expected login 'testuser', got '%s'", in.Login)
				}
				return "testsalt", nil
			},
			authenticateFunc: func(ctx context.Context, in models.AuthenticateIn) error {
				if in.Login != "testuser" {
					t.Errorf("Expected login 'testuser', got '%s'", in.Login)
				}
				return nil
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewLoginCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, credentials)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("login API error", func(t *testing.T) {
		expectedErr := errors.New("user not found")
		mockAPI := &MockLoginAPIClient{
			loginFunc: func(ctx context.Context, in models.LoginIn) (salt string, err error) {
				return "", expectedErr
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewLoginCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, credentials)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "ошибка при входе: "+expectedErr.Error() {
			t.Errorf("Expected error message containing 'ошибка при входе', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("authenticate API error", func(t *testing.T) {
		expectedErr := errors.New("invalid credentials")
		mockAPI := &MockLoginAPIClient{
			loginFunc: func(ctx context.Context, in models.LoginIn) (salt string, err error) {
				return "testsalt", nil
			},
			authenticateFunc: func(ctx context.Context, in models.AuthenticateIn) error {
				return expectedErr
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewLoginCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, credentials)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "ошибка при аутентификации: "+expectedErr.Error() {
			t.Errorf("Expected error message containing 'ошибка при аутентификации', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("invalid args type", func(t *testing.T) {
		mockAPI := &MockLoginAPIClient{}
		crypto := crypto.NewCrypto()
		cmd := NewLoginCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, "invalid args")

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "неверный тип аргументов для команды входа" {
			t.Errorf("Expected error message 'неверный тип аргументов для команды входа', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestLoginCommand_GetName(t *testing.T) {
	mockAPI := &MockLoginAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewLoginCommand(mockAPI, crypto)

	name := cmd.GetName()
	if name != "login" {
		t.Errorf("Expected name 'login', got '%s'", name)
	}
}

func TestLoginCommand_GetDescription(t *testing.T) {
	mockAPI := &MockLoginAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewLoginCommand(mockAPI, crypto)

	description := cmd.GetDescription()
	expected := "Вход в существующий аккаунт"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
