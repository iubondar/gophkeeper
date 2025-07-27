package cmd

import (
	"context"
)

// VersionCommand представляет команду для отображения версии
type VersionCommand struct {
	version   string
	buildTime string
}

// NewVersionCommand создает новую команду версии
func NewVersionCommand(version, buildTime string) *VersionCommand {
	return &VersionCommand{
		version:   version,
		buildTime: buildTime,
	}
}

// Execute выполняет команду версии
func (c *VersionCommand) Execute(ctx context.Context, args any) (any, error) {
	return &VersionResult{
		Version:   c.version,
		BuildTime: c.buildTime,
	}, nil
}

// GetName возвращает имя команды
func (c *VersionCommand) GetName() string {
	return "version"
}

// GetDescription возвращает описание команды
func (c *VersionCommand) GetDescription() string {
	return "Показать версию и дату сборки клиента"
}

// VersionResult представляет результат выполнения команды версии
type VersionResult struct {
	Version   string
	BuildTime string
}
