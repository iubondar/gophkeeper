package cmd

import (
	"context"
)

type HealthAPIClient interface {
	HealthCheck(ctx context.Context) error
}

type HealthCommand struct {
	apiClient HealthAPIClient
}

func NewHealthCommand(apiClient HealthAPIClient) *HealthCommand {
	return &HealthCommand{apiClient: apiClient}
}

func (c *HealthCommand) Execute(ctx context.Context, data any) (any, error) {
	err := c.apiClient.HealthCheck(ctx)
	return nil, err
}

func (c *HealthCommand) GetName() string {
	return "health"
}

func (c *HealthCommand) GetDescription() string {
	return "Проверить доступность сервера GophKeeper"
}
