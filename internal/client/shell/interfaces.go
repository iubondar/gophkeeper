package shell

import (
	"context"
	"gophkeeper/internal/models"
)

// CommandRegistry интерфейс для реестра команд
type CommandRegistry interface {
	Execute(ctx context.Context, commandName string, args any) (any, error)
}

// InputHandler интерфейс для обработки пользовательского ввода
type InputHandler interface {
	GetUserChoice() string
	GetUserCredentials() (*models.UserCredentials, error)
	GetSecretName() (string, error)
	GetTextData() (*models.TextSecretData, error)
	GetLoginPasswordData() (*models.LoginPasswordData, error)
	GetCardData() (*models.CardData, error)
	GetFileData() (*models.FileData, error)
	GetUpdatedTextData(name string) (*models.TextSecretData, error)
	GetUpdatedLoginPasswordData(name string) (*models.LoginPasswordData, error)
	GetUpdatedCardData(name string) (*models.CardData, error)
	GetUpdatedFileData(name string) (*models.FileData, error)
	PromptEnterNewData()
}

// Display интерфейс для отображения информации
type Display interface {
	Welcome()
	ServerConnected()
	InvalidChoice()
	ErrorMsg(err error)
	SwitchToUserMenuNotice()
	Goodbye()
	Logout()
	BackNotAllowed()
	SuccessRegistration()
	SuccessLogin()
	SuccessUpload()
	SuccessUpdate()
	SuccessGet()
	SuccessDelete()
	SuccessGeneric(commandName string)
	DisplayTextSecret(secret *models.TextSecretData, metadata string)
	DisplayLoginPassword(secret *models.LoginPasswordData, metadata string)
	DisplayCardData(secret *models.CardData, metadata string)
	DisplayFileData(secret *models.FileData, metadata string)
	SuccessDownloadFile(filePath string)
	DisplaySecretInfo(secretName, secretType, metadata string, version int)
	DisplayVersion(version, buildTime string)
	MenuTitle(title string)
	MenuItem(id, title, description string)
	MenuChoice()
	MenuError(message string)
}
