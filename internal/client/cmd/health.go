package cmd

import (
	"context"
)

// HealthAPIClient интерфейс для проверки доступности сервера.
type HealthAPIClient interface {
	HealthCheck(ctx context.Context) error
}

// HealthCommand представляет команду проверки доступности сервера.
// Выполняет ping запрос к серверу для проверки его работоспособности.
type HealthCommand struct {
	apiClient HealthAPIClient
}

// NewHealthCommand создает новую команду проверки здоровья сервера.
//
// Параметры:
//   - apiClient: API клиент для проверки доступности сервера
//
// Возвращает:
//   - *HealthCommand: новый экземпляр команды проверки здоровья
func NewHealthCommand(apiClient HealthAPIClient) *HealthCommand {
	return &HealthCommand{apiClient: apiClient}
}

// Execute выполняет команду проверки доступности сервера.
// Отправляет запрос к эндпоинту /health для проверки работоспособности сервера.
//
// Параметры:
//   - ctx: контекст выполнения
//   - data: не используется (может быть nil)
//
// Возвращает:
//   - any: nil при успешной проверке
//   - error: ошибка в случае недоступности сервера
func (c *HealthCommand) Execute(ctx context.Context, data any) (any, error) {
	err := c.apiClient.HealthCheck(ctx)
	return nil, err
}

// GetName возвращает имя команды.
//
// Возвращает:
//   - string: "health"
func (c *HealthCommand) GetName() string {
	return "health"
}

// GetDescription возвращает описание команды.
//
// Возвращает:
//   - string: описание команды проверки здоровья
func (c *HealthCommand) GetDescription() string {
	return "Проверить доступность сервера GophKeeper"
}
