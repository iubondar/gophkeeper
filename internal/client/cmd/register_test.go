package cmd

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/internal/client/crypto"
	"gophkeeper/internal/models"
)

// MockRegisterAPIClient представляет мок для RegisterAPIClient
type MockRegisterAPIClient struct {
	registerFunc func(ctx context.Context, in models.RegisterIn) error
}

func (m *MockRegisterAPIClient) Register(ctx context.Context, in models.RegisterIn) error {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, in)
	}
	return nil
}

func TestNewRegisterCommand(t *testing.T) {
	mockAPI := &MockRegisterAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewRegisterCommand(mockAPI, crypto)

	if cmd == nil {
		t.Fatal("NewRegisterCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}

	if cmd.crypto != crypto {
		t.Error("crypto not set correctly")
	}
}

func TestRegisterCommand_Execute(t *testing.T) {
	ctx := context.Background()
	credentials := models.UserCredentials{
		Login:    "newuser",
		Password: "newpass",
	}

	t.Run("successful registration", func(t *testing.T) {
		mockAPI := &MockRegisterAPIClient{
			registerFunc: func(ctx context.Context, in models.RegisterIn) error {
				if in.Login != "newuser" {
					t.Errorf("Expected login 'newuser', got '%s'", in.Login)
				}
				if in.Salt == "" {
					t.Error("Expected salt to be set")
				}
				if in.PasswordHash == "" {
					t.Error("Expected password hash to be set")
				}
				return nil
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewRegisterCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, credentials)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("registration API error", func(t *testing.T) {
		expectedErr := errors.New("user already exists")
		mockAPI := &MockRegisterAPIClient{
			registerFunc: func(ctx context.Context, in models.RegisterIn) error {
				return expectedErr
			},
		}
		crypto := crypto.NewCrypto()
		cmd := NewRegisterCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, credentials)

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "ошибка при регистрации: "+expectedErr.Error() {
			t.Errorf("Expected error message containing 'ошибка при регистрации', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("invalid args type", func(t *testing.T) {
		mockAPI := &MockRegisterAPIClient{}
		crypto := crypto.NewCrypto()
		cmd := NewRegisterCommand(mockAPI, crypto)

		result, err := cmd.Execute(ctx, "invalid args")

		if err == nil {
			t.Fatal("Expected error")
		}

		if err.Error() != "неверный тип аргументов для команды регистрации" {
			t.Errorf("Expected error message 'неверный тип аргументов для команды регистрации', got '%s'", err.Error())
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestRegisterCommand_GetName(t *testing.T) {
	mockAPI := &MockRegisterAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewRegisterCommand(mockAPI, crypto)

	name := cmd.GetName()
	if name != "register" {
		t.Errorf("Expected name 'register', got '%s'", name)
	}
}

func TestRegisterCommand_GetDescription(t *testing.T) {
	mockAPI := &MockRegisterAPIClient{}
	crypto := crypto.NewCrypto()
	cmd := NewRegisterCommand(mockAPI, crypto)

	description := cmd.GetDescription()
	expected := "Регистрация нового пользователя"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
