package cmd

import (
	"context"
)

// VersionCommand представляет команду для отображения версии клиента.
// Возвращает информацию о версии и дате сборки приложения.
type VersionCommand struct {
	version   string
	buildTime string
}

// NewVersionCommand создает новую команду версии.
//
// Параметры:
//   - version: версия приложения
//   - buildTime: дата и время сборки приложения
//
// Возвращает:
//   - *VersionCommand: новый экземпляр команды версии
func NewVersionCommand(version, buildTime string) *VersionCommand {
	return &VersionCommand{
		version:   version,
		buildTime: buildTime,
	}
}

// Execute выполняет команду версии.
// Возвращает информацию о версии и дате сборки клиента.
//
// Параметры:
//   - ctx: контекст выполнения (не используется)
//   - args: не используется (может быть nil)
//
// Возвращает:
//   - any: VersionResult с информацией о версии
//   - error: всегда nil
func (c *VersionCommand) Execute(ctx context.Context, args any) (any, error) {
	return &VersionResult{
		Version:   c.version,
		BuildTime: c.buildTime,
	}, nil
}

// GetName возвращает имя команды.
//
// Возвращает:
//   - string: "version"
func (c *VersionCommand) GetName() string {
	return "version"
}

// GetDescription возвращает описание команды.
//
// Возвращает:
//   - string: описание команды версии
func (c *VersionCommand) GetDescription() string {
	return "Показать версию и дату сборки клиента"
}

// VersionResult представляет результат выполнения команды версии.
// Содержит информацию о версии и дате сборки приложения.
type VersionResult struct {
	Version   string // Версия приложения
	BuildTime string // Дата и время сборки
}
