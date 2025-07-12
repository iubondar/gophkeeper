package cmd

import (
	"context"
	"gophkeeper/internal/models"
)

// CommandProducer по сути является фабрикой команд
type CommandProducer interface {
	Register() RegisterHandler
}

// GophKeeperClient это интерфейс для работы с API сервера, который предоставляет все API методы
type GophKeeperClient interface {
	Register(ctx context.Context, in models.RegisterIn) error
}

// CommandFactory это фабрика команд
type CommandFactory struct {
	apiClient GophKeeperClient
}

// NewCommandFactory создает новый экземпляр CommandFactory
func NewCommandFactory(apiClient GophKeeperClient) *CommandFactory {
	return &CommandFactory{apiClient: apiClient}
}

// Register возвращает команду регистрации
func (f *CommandFactory) Register() RegisterHandler {
	return NewRegisterHandler(f.apiClient)
}
