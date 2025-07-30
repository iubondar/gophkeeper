package cmd

import (
	"context"
	"errors"
	"testing"
)

// MockHealthAPIClient представляет мок для HealthAPIClient
type MockHealthAPIClient struct {
	healthCheckFunc func(ctx context.Context) error
}

func (m *MockHealthAPIClient) HealthCheck(ctx context.Context) error {
	if m.healthCheckFunc != nil {
		return m.healthCheckFunc(ctx)
	}
	return nil
}

func TestNewHealthCommand(t *testing.T) {
	mockAPI := &MockHealthAPIClient{}
	cmd := NewHealthCommand(mockAPI)

	if cmd == nil {
		t.Fatal("NewHealthCommand returned nil")
	}

	if cmd.apiClient != mockAPI {
		t.Error("apiClient not set correctly")
	}
}

func TestHealthCommand_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("successful health check", func(t *testing.T) {
		mockAPI := &MockHealthAPIClient{
			healthCheckFunc: func(ctx context.Context) error {
				return nil
			},
		}

		cmd := NewHealthCommand(mockAPI)
		result, err := cmd.Execute(ctx, nil)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("health check with error", func(t *testing.T) {
		expectedErr := errors.New("server unavailable")
		mockAPI := &MockHealthAPIClient{
			healthCheckFunc: func(ctx context.Context) error {
				return expectedErr
			},
		}

		cmd := NewHealthCommand(mockAPI)
		result, err := cmd.Execute(ctx, nil)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}

		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestHealthCommand_GetName(t *testing.T) {
	mockAPI := &MockHealthAPIClient{}
	cmd := NewHealthCommand(mockAPI)

	name := cmd.GetName()
	if name != "health" {
		t.Errorf("Expected name 'health', got '%s'", name)
	}
}

func TestHealthCommand_GetDescription(t *testing.T) {
	mockAPI := &MockHealthAPIClient{}
	cmd := NewHealthCommand(mockAPI)

	description := cmd.GetDescription()
	expected := "Проверить доступность сервера GophKeeper"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
