package shell

import (
	"errors"
	"fmt"
	"gophkeeper/internal/client/cmd"
	"gophkeeper/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// newTestShell создает Shell для тестов с правильно настроенным MenuManager
func newTestShell(registry CommandRegistry, inputHandler InputHandler, display Display) *Shell {
	menuManager := NewMenuManager()
	menuManager.SetDisplay(display)

	return &Shell{
		commandRegistry: registry,
		menuManager:     menuManager,
		inputHandler:    inputHandler,
		display:         display,
		version:         "test",
		buildTime:       "test",
	}
}

// newTestShellWithMenuManager создает Shell для тестов с пользовательским MenuManager
func newTestShellWithMenuManager(registry CommandRegistry, inputHandler InputHandler, display Display, menuManager *MenuManager) *Shell {
	menuManager.SetDisplay(display)

	return &Shell{
		commandRegistry: registry,
		menuManager:     menuManager,
		inputHandler:    inputHandler,
		display:         display,
		version:         "test",
		buildTime:       "test",
	}
}

func TestNewShell(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	version := "1.0.0"
	buildTime := "2024-01-01"

	shell := NewShell(mockRegistry, mockInputHandler, mockDisplay, version, buildTime)

	if shell == nil {
		t.Fatal("NewShell returned nil")
	}

	if shell.version != version {
		t.Errorf("Expected version to be %s, got %s", version, shell.version)
	}

	if shell.buildTime != buildTime {
		t.Errorf("Expected buildTime to be %s, got %s", buildTime, shell.buildTime)
	}

	if shell.menuManager == nil {
		t.Error("Expected menuManager to be created")
	}

	if shell.inputHandler == nil {
		t.Error("Expected inputHandler to be created")
	}

	if shell.display == nil {
		t.Error("Expected display to be created")
	}
}

func TestNewShellWithDefaults(t *testing.T) {
	// Создаем реальный CommandRegistry для теста
	registry := &cmd.CommandRegistry{}
	version := "1.0.0"
	buildTime := "2024-01-01"

	shell := NewShellWithDefaults(registry, version, buildTime)

	if shell == nil {
		t.Fatal("NewShellWithDefaults returned nil")
	}

	if shell.version != version {
		t.Errorf("Expected version to be %s, got %s", version, shell.version)
	}

	if shell.buildTime != buildTime {
		t.Errorf("Expected buildTime to be %s, got %s", buildTime, shell.buildTime)
	}
}

func TestShell_executeCommand_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	credentials := &models.UserCredentials{
		Login:    "testuser",
		Password: "testpass",
	}

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetUserCredentials().Return(credentials, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandRegister, gomock.Any()).Return("success", nil)

	// Выполняем тест
	result, err := shell.executeCommand(CommandRegister)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %v", result)
	}
}

func TestShell_executeCommand_Register_InputError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	expectedErr := errors.New("input error")

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetUserCredentials().Return(nil, expectedErr)

	// Выполняем тест
	result, err := shell.executeCommand(CommandRegister)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestShell_executeCommand_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	credentials := &models.UserCredentials{
		Login:    "testuser",
		Password: "testpass",
	}

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetUserCredentials().Return(credentials, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandLogin, gomock.Any()).Return("success", nil)

	// Выполняем тест
	result, err := shell.executeCommand(CommandLogin)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %v", result)
	}
}

func TestShell_executeCommand_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	secretName := "test_secret"

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetSecretName().Return(secretName, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandGet, secretName).Return("success", nil)

	// Выполняем тест
	result, err := shell.executeCommand(CommandGet)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %v", result)
	}
}

func TestShell_executeCommand_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	secretName := "test_secret"

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetSecretName().Return(secretName, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandDelete, secretName).Return("success", nil)

	// Выполняем тест
	result, err := shell.executeCommand(CommandDelete)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %v", result)
	}
}

func TestShell_executeCommand_Upload_Text(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	// Устанавливаем состояние действия
	mockMenuManager.SetActionState(CommandUpload, CommandText)

	textData := &models.TextSecretData{
		Name:     "test_secret",
		Text:     "test text",
		Metadata: "test metadata",
	}

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetTextData().Return(textData, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandUpload, textData).Return("success", nil)

	// Выполняем тест
	result, err := shell.executeCommand(CommandUpload)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %v", result)
	}
}

func TestShell_executeCommand_Upload_LoginPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	// Устанавливаем состояние действия
	mockMenuManager.SetActionState(CommandUpload, CommandLoginPassword)

	loginPasswordData := &models.LoginPasswordData{
		Name:     "test_secret",
		Login:    "testuser",
		Password: "testpass",
		URL:      "https://example.com",
	}

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetLoginPasswordData().Return(loginPasswordData, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandUpload, loginPasswordData).Return("success", nil)

	// Выполняем тест
	result, err := shell.executeCommand(CommandUpload)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %v", result)
	}
}

func TestShell_executeCommand_Upload_Card(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	// Устанавливаем состояние действия
	mockMenuManager.SetActionState(CommandUpload, CommandCard)

	cardData := &models.CardData{
		Name:   "test_card",
		Number: "1234567890123456",
		Holder: "Test User",
		Expiry: "12/25",
		CVV:    "123",
	}

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetCardData().Return(cardData, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandUpload, cardData).Return("success", nil)

	// Выполняем тест
	result, err := shell.executeCommand(CommandUpload)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %v", result)
	}
}

func TestShell_executeCommand_Upload_File(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	// Устанавливаем состояние действия
	mockMenuManager.SetActionState(CommandUpload, CommandFile)

	fileData := &models.FileData{
		Name:     "test_file",
		FilePath: "/path/to/file.txt",
	}

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetFileData().Return(fileData, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandUpload, fileData).Return("success", nil)

	// Выполняем тест
	result, err := shell.executeCommand(CommandUpload)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected result 'success', got %v", result)
	}
}

func TestShell_executeCommand_Upload_UnknownDataType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(mockRegistry, nil, mockDisplay, mockMenuManager)

	// Устанавливаем состояние действия с неизвестным типом
	mockMenuManager.SetActionState(CommandUpload, "unknown_type")

	// Выполняем тест
	result, err := shell.executeCommand(CommandUpload)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestShell_executeCommand_Upload_NoDataType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(mockRegistry, nil, mockDisplay, mockMenuManager)

	// Не устанавливаем состояние действия

	// Выполняем тест
	result, err := shell.executeCommand(CommandUpload)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestShell_executeCommand_UnknownCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, nil, mockDisplay)

	// Выполняем тест
	result, err := shell.executeCommand("unknown_command")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestShell_handleCommandResult_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(nil, nil, mockDisplay)

	// Создаем результат команды get
	getResult := &cmd.GetSecretResult{
		Type:     models.SecretTypeText,
		Data:     &models.TextSecretData{Name: "test", Text: "test text"},
		Metadata: "test metadata",
	}

	// Настраиваем ожидания
	mockDisplay.EXPECT().DisplayTextSecret(gomock.Any(), "test metadata")

	// Выполняем тест
	shell.handleCommandResult(CommandGet, getResult)
}

func TestShell_handleCommandResult_Get_LoginPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(nil, nil, mockDisplay)

	// Создаем результат команды get для логин/пароль
	getResult := &cmd.GetSecretResult{
		Type: models.SecretTypeLoginPassword,
		Data: &models.LoginPasswordData{
			Name:     "test",
			Login:    "testuser",
			Password: "testpass",
			URL:      "https://example.com",
		},
		Metadata: "test metadata",
	}

	// Настраиваем ожидания
	mockDisplay.EXPECT().DisplayLoginPassword(gomock.Any(), "test metadata")

	// Выполняем тест
	shell.handleCommandResult(CommandGet, getResult)
}

func TestShell_handleCommandResult_Get_Card(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(nil, nil, mockDisplay)

	// Создаем результат команды get для карты
	getResult := &cmd.GetSecretResult{
		Type: models.SecretTypeCard,
		Data: &models.CardData{
			Name:   "test",
			Number: "1234567890123456",
			Holder: "Test User",
			Expiry: "12/25",
			CVV:    "123",
		},
		Metadata: "test metadata",
	}

	// Настраиваем ожидания
	mockDisplay.EXPECT().DisplayCardData(gomock.Any(), "test metadata")

	// Выполняем тест
	shell.handleCommandResult(CommandGet, getResult)
}

func TestShell_handleCommandResult_Get_File(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(nil, nil, mockDisplay)

	// Создаем результат команды get для файла
	getResult := &cmd.GetSecretResult{
		Type: models.SecretTypeFile,
		Data: &models.FileData{
			Name:     "test",
			FilePath: "/path/to/file.txt",
		},
		Metadata: "test metadata",
	}

	// Настраиваем ожидания
	mockDisplay.EXPECT().DisplayFileData(gomock.Any(), "test metadata")
	mockDisplay.EXPECT().SuccessDownloadFile("/path/to/file.txt")

	// Выполняем тест
	shell.handleCommandResult(CommandGet, getResult)
}

func TestShell_handleCommandResult_NonGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(nil, nil, mockDisplay)

	// Настраиваем ожидания
	mockDisplay.EXPECT().SuccessUpload()

	// Выполняем тест
	shell.handleCommandResult(CommandUpload, "success")
}

func TestShell_printSuccessMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(nil, nil, mockDisplay)

	// Настраиваем ожидания для различных команд
	mockDisplay.EXPECT().SuccessRegistration()
	mockDisplay.EXPECT().SuccessLogin()
	mockDisplay.EXPECT().SuccessUpload()
	mockDisplay.EXPECT().SuccessUpdate()
	mockDisplay.EXPECT().SuccessGet()
	mockDisplay.EXPECT().SuccessDelete()
	mockDisplay.EXPECT().SuccessGeneric("unknown")

	// Тестируем различные команды
	commands := []string{CommandRegister, CommandLogin, CommandUpload, CommandUpdate, CommandGet, CommandDelete, "unknown"}

	for _, cmd := range commands {
		shell.printSuccessMessage(cmd)
	}
}

func TestShell_checkServerHealth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)

	shell := newTestShell(mockRegistry, nil, nil)

	// Настраиваем ожидания для успешной проверки
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, nil)

	// Выполняем тест
	err := shell.checkServerHealth()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShell_checkServerHealth_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)

	shell := newTestShell(mockRegistry, nil, nil)

	expectedErr := errors.New("server error")

	// Настраиваем ожидания для неуспешной проверки
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, expectedErr)

	// Выполняем тест
	err := shell.checkServerHealth()

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestShell_handleExitCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(nil, nil, mockDisplay)

	// Настраиваем ожидания
	mockDisplay.EXPECT().Goodbye()

	// Выполняем тест
	err := shell.handleExitCommand()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShell_handleLogoutCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(nil, nil, mockDisplay, mockMenuManager)

	// Устанавливаем начальное состояние
	mockMenuManager.SwitchToState(MenuStateAuthenticated)
	mockMenuManager.SetActionState("upload", "text")

	// Настраиваем ожидания
	mockDisplay.EXPECT().Logout()

	// Выполняем тест
	shell.handleLogoutCommand()

	// Проверяем, что состояние изменилось
	if mockMenuManager.GetCurrentState() != MenuStateMain {
		t.Error("Expected state to be MenuStateMain")
	}

	if mockMenuManager.GetActionState() != nil {
		t.Error("Expected actionState to be cleared")
	}
}

func TestShell_handleBackCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(nil, nil, mockDisplay, mockMenuManager)

	// Устанавливаем состояние, из которого можно вернуться назад
	mockMenuManager.SwitchToState(MenuStateAuthenticated)

	// Выполняем тест
	shell.handleBackCommand()

	// Проверяем, что состояние изменилось
	if mockMenuManager.GetCurrentState() != MenuStateMain {
		t.Error("Expected state to be MenuStateMain")
	}
}

func TestShell_handleBackCommand_NoParent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(nil, nil, mockDisplay, mockMenuManager)

	// Устанавливаем состояние, из которого нельзя вернуться назад
	mockMenuManager.SwitchToState(MenuStateMain)

	// Настраиваем ожидания
	mockDisplay.EXPECT().BackNotAllowed()

	// Выполняем тест
	shell.handleBackCommand()

	// Проверяем, что состояние не изменилось
	if mockMenuManager.GetCurrentState() != MenuStateMain {
		t.Error("Expected state to remain MenuStateMain")
	}
}

func TestShell_handleActionCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMenuManager := NewMenuManager()

	shell := &Shell{
		menuManager: mockMenuManager,
	}

	// Выполняем тест
	shell.handleActionCommand(CommandUpload)

	// Проверяем, что состояние действия установлено
	actionState := mockMenuManager.GetActionState()
	if actionState == nil {
		t.Error("Expected actionState to be set")
	}
	if actionState.Action != CommandUpload {
		t.Errorf("Expected action to be %s, got %s", CommandUpload, actionState.Action)
	}

	// Проверяем, что переключились в меню выбора типа данных
	if mockMenuManager.GetCurrentState() != MenuStateDataType {
		t.Errorf("Expected state to be %s, got %s", MenuStateDataType, mockMenuManager.GetCurrentState())
	}
}

func TestShell_handleDataTypeCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	// Устанавливаем состояние действия
	mockMenuManager.SetActionState(CommandUpload, "")
	mockMenuManager.SwitchToState(MenuStateDataType)

	textData := &models.TextSecretData{
		Name:     "test_secret",
		Text:     "test text",
		Metadata: "test metadata",
	}

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetTextData().Return(textData, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandUpload, textData).Return("success", nil)
	mockDisplay.EXPECT().SuccessUpload()

	// Выполняем тест
	shell.handleDataTypeCommand(CommandText)

	// Проверяем, что состояние действия очищено
	if mockMenuManager.GetActionState() != nil {
		t.Error("Expected actionState to be cleared")
	}

	// Проверяем, что вернулись в меню авторизованного пользователя
	if mockMenuManager.GetCurrentState() != MenuStateAuthenticated {
		t.Errorf("Expected state to be %s, got %s", MenuStateAuthenticated, mockMenuManager.GetCurrentState())
	}
}

func TestShell_handleDataTypeCommand_NoActionState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()

	shell := &Shell{
		display:     mockDisplay,
		menuManager: mockMenuManager,
	}

	// Не устанавливаем состояние действия
	mockMenuManager.SwitchToState(MenuStateDataType)

	// Выполняем тест
	shell.handleDataTypeCommand(CommandText)

	// Проверяем, что состояние не изменилось
	if mockMenuManager.GetCurrentState() != MenuStateDataType {
		t.Error("Expected state to remain MenuStateDataType")
	}
}

func TestShell_handleGetDeleteCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	secretName := "test_secret"

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetSecretName().Return(secretName, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandGet, secretName).Return("success", nil)
	mockDisplay.EXPECT().SuccessGet()

	// Выполняем тест
	shell.handleGetDeleteCommand(CommandGet)
}

func TestShell_handleGetDeleteCommand_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	expectedErr := errors.New("command error")

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetSecretName().Return("", expectedErr)
	mockDisplay.EXPECT().ErrorMsg(gomock.Any())

	// Выполняем тест
	shell.handleGetDeleteCommand(CommandGet)
}

func TestShell_handleVersionCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, nil, mockDisplay)

	versionResult := &cmd.VersionResult{
		Version:   "1.0.0",
		BuildTime: "2024-01-01",
	}

	// Настраиваем ожидания
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandVersion, nil).Return(versionResult, nil)
	mockDisplay.EXPECT().DisplayVersion("1.0.0", "2024-01-01")

	// Выполняем тест
	shell.handleVersionCommand()
}

func TestShell_handleVersionCommand_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, nil, mockDisplay)

	expectedErr := errors.New("version error")

	// Настраиваем ожидания
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandVersion, nil).Return(nil, expectedErr)
	mockDisplay.EXPECT().ErrorMsg(expectedErr)

	// Выполняем тест
	shell.handleVersionCommand()
}

func TestShell_handleUpdateCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	secretName := "test_secret"
	showResult := &cmd.ShowSecretResult{
		Type:     models.SecretTypeText,
		Metadata: "test metadata",
		Version:  1,
	}
	updatedTextData := &models.TextSecretData{
		Name:     secretName,
		Text:     "updated text",
		Metadata: "updated metadata",
	}
	updateData := &cmd.UpdateData{
		SecretName: secretName,
		Version:    1,
		Type:       models.SecretTypeText,
		Data:       updatedTextData,
	}

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetSecretName().Return(secretName, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandShow, secretName).Return(showResult, nil)
	mockDisplay.EXPECT().DisplaySecretInfo(secretName, models.SecretTypeText, "test metadata", 1)
	mockInputHandler.EXPECT().PromptEnterNewData()
	mockInputHandler.EXPECT().GetUpdatedTextData(secretName).Return(updatedTextData, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandUpdate, updateData).Return("success", nil)
	mockDisplay.EXPECT().SuccessUpdate()

	// Выполняем тест
	shell.handleUpdateCommand()
}

func TestShell_handleUpdateCommand_GetSecretNameError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	expectedErr := errors.New("input error")

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetSecretName().Return("", expectedErr)
	mockDisplay.EXPECT().ErrorMsg(gomock.Any())

	// Выполняем тест
	shell.handleUpdateCommand()
}

func TestShell_handleUpdateCommand_ShowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	secretName := "test_secret"
	expectedErr := errors.New("show error")

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetSecretName().Return(secretName, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandShow, secretName).Return(nil, expectedErr)
	mockDisplay.EXPECT().ErrorMsg(gomock.Any())

	// Выполняем тест
	shell.handleUpdateCommand()
}

func TestShell_handleUpdateCommand_InvalidShowResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := newTestShell(mockRegistry, mockInputHandler, mockDisplay)

	secretName := "test_secret"

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetSecretName().Return(secretName, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandShow, secretName).Return("invalid_result", nil)
	mockDisplay.EXPECT().ErrorMsg(gomock.Any())

	// Выполняем тест
	shell.handleUpdateCommand()
}

func TestShell_handleUpdateCommand_UnsupportedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := &Shell{
		commandRegistry: mockRegistry,
		inputHandler:    mockInputHandler,
		display:         mockDisplay,
	}

	secretName := "test_secret"
	showResult := &cmd.ShowSecretResult{
		Type:     "unsupported_type",
		Metadata: "test metadata",
		Version:  1,
	}

	// Настраиваем ожидания
	mockInputHandler.EXPECT().GetSecretName().Return(secretName, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandShow, secretName).Return(showResult, nil)
	mockDisplay.EXPECT().DisplaySecretInfo(secretName, "unsupported_type", "test metadata", 1)
	mockInputHandler.EXPECT().PromptEnterNewData()
	mockDisplay.EXPECT().ErrorMsg(gomock.Any())

	// Выполняем тест
	shell.handleUpdateCommand()
}

func TestShell_Run_HealthCheckError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := &Shell{
		commandRegistry: mockRegistry,
		inputHandler:    mockInputHandler,
		display:         mockDisplay,
		menuManager:     NewMenuManager(),
	}

	expectedErr := errors.New("health check failed")

	// Настраиваем ожидания
	mockDisplay.EXPECT().Welcome()
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, expectedErr)

	// Выполняем тест
	err := shell.Run()

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestShell_Run_ExitCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := &Shell{
		commandRegistry: mockRegistry,
		inputHandler:    mockInputHandler,
		display:         mockDisplay,
		menuManager:     NewMenuManager(),
	}
	shell.menuManager.SetDisplay(mockDisplay)

	mockDisplay.EXPECT().Welcome()
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, nil)
	mockDisplay.EXPECT().ServerConnected()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("4")
	mockDisplay.EXPECT().Goodbye()

	err := shell.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShell_Run_InvalidChoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := &Shell{
		commandRegistry: mockRegistry,
		inputHandler:    mockInputHandler,
		display:         mockDisplay,
		menuManager:     NewMenuManager(),
	}
	shell.menuManager.SetDisplay(mockDisplay)

	mockDisplay.EXPECT().Welcome()
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, nil)
	mockDisplay.EXPECT().ServerConnected()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("99")
	mockDisplay.EXPECT().InvalidChoice()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("4")
	mockDisplay.EXPECT().Goodbye()

	err := shell.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShell_Run_LogoutCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := &Shell{
		commandRegistry: mockRegistry,
		inputHandler:    mockInputHandler,
		display:         mockDisplay,
		menuManager:     NewMenuManager(),
	}
	shell.menuManager.SetDisplay(mockDisplay)

	shell.menuManager.SwitchToState(MenuStateAuthenticated)

	mockDisplay.EXPECT().Welcome()
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, nil)
	mockDisplay.EXPECT().ServerConnected()
	expectUserMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("5")
	mockDisplay.EXPECT().Logout()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("4")
	mockDisplay.EXPECT().Goodbye()

	err := shell.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShell_Run_BackCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()
	mockMenuManager.SetDisplay(mockDisplay)

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	shell.menuManager.SwitchToState(MenuStateDataType)

	mockDisplay.EXPECT().Welcome()
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, nil)
	mockDisplay.EXPECT().ServerConnected()
	expectDataTypeMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("5")
	expectUserMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("5")
	mockDisplay.EXPECT().Logout()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("4")
	mockDisplay.EXPECT().Goodbye()

	err := shell.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShell_Run_BackCommand_NoParent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()
	mockMenuManager.SetDisplay(mockDisplay)

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	mockDisplay.EXPECT().Welcome()
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, nil)
	mockDisplay.EXPECT().ServerConnected()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("5")
	mockDisplay.EXPECT().InvalidChoice()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("4")
	mockDisplay.EXPECT().Goodbye()

	err := shell.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShell_Run_VersionCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()
	mockMenuManager.SetDisplay(mockDisplay)

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	versionResult := &cmd.VersionResult{
		Version:   "1.0.0",
		BuildTime: "2024-01-01",
	}

	mockDisplay.EXPECT().Welcome()
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, nil)
	mockDisplay.EXPECT().ServerConnected()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("3")
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandVersion, nil).Return(versionResult, nil)
	mockDisplay.EXPECT().DisplayVersion("1.0.0", "2024-01-01")
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("4")
	mockDisplay.EXPECT().Goodbye()

	err := shell.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShell_Run_RegisterCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()
	mockMenuManager.SetDisplay(mockDisplay)

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	credentials := &models.UserCredentials{
		Login:    "testuser",
		Password: "testpass",
	}

	mockDisplay.EXPECT().Welcome()
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, nil)
	mockDisplay.EXPECT().ServerConnected()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("1")
	mockInputHandler.EXPECT().GetUserCredentials().Return(credentials, nil)
	mockRegistry.EXPECT().Execute(gomock.Any(), CommandRegister, gomock.Any()).Return("success", nil)
	mockDisplay.EXPECT().SuccessRegistration()
	mockDisplay.EXPECT().SwitchToUserMenuNotice()
	expectUserMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("5")
	mockDisplay.EXPECT().Logout()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("4")
	mockDisplay.EXPECT().Goodbye()

	err := shell.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShell_Run_RegisterCommand_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)
	mockMenuManager := NewMenuManager()
	mockMenuManager.SetDisplay(mockDisplay)

	shell := newTestShellWithMenuManager(mockRegistry, mockInputHandler, mockDisplay, mockMenuManager)

	expectedErr := errors.New("registration error")

	mockDisplay.EXPECT().Welcome()
	mockRegistry.EXPECT().Execute(gomock.Any(), "health", nil).Return(nil, nil)
	mockDisplay.EXPECT().ServerConnected()
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("1")
	mockInputHandler.EXPECT().GetUserCredentials().Return(nil, expectedErr)
	mockDisplay.EXPECT().ErrorMsg(gomock.Any())
	expectMainMenu(mockDisplay)
	mockInputHandler.EXPECT().GetUserChoice().Return("4")
	mockDisplay.EXPECT().Goodbye()

	err := shell.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func expectMainMenu(mockDisplay *MockDisplay) {
	mockDisplay.EXPECT().MenuTitle("Главное меню")
	mockDisplay.EXPECT().MenuItem("1", "Регистрация", "Создать новый аккаунт")
	mockDisplay.EXPECT().MenuItem("2", "Вход", "Войти в существующий аккаунт")
	mockDisplay.EXPECT().MenuItem("3", "Версия", "Показать версию и дату сборки")
	mockDisplay.EXPECT().MenuItem("4", "Выход", "Завершить работу")
	mockDisplay.EXPECT().MenuChoice()
}

func expectUserMenu(mockDisplay *MockDisplay) {
	mockDisplay.EXPECT().MenuTitle("Меню пользователя")
	mockDisplay.EXPECT().MenuItem("1", "Загрузить", "Загрузить секрет")
	mockDisplay.EXPECT().MenuItem("2", "Обновить", "Обновить секрет")
	mockDisplay.EXPECT().MenuItem("3", "Получить", "Получить секрет")
	mockDisplay.EXPECT().MenuItem("4", "Удалить", "Удалить секрет")
	mockDisplay.EXPECT().MenuItem("5", "Выйти из аккаунта", "Выйти из аккаунта")
	mockDisplay.EXPECT().MenuChoice()
}

func expectDataTypeMenu(mockDisplay *MockDisplay) {
	mockDisplay.EXPECT().MenuTitle("Выберите тип данных")
	mockDisplay.EXPECT().MenuItem("1", "Текст", "Простой текстовый секрет")
	mockDisplay.EXPECT().MenuItem("2", "Логин-Пароль", "Данные для входа в систему")
	mockDisplay.EXPECT().MenuItem("3", "Данные карты", "Информация о банковской карте")
	mockDisplay.EXPECT().MenuItem("4", "Файл", "Зашифрованный файл")
	mockDisplay.EXPECT().MenuItem("5", "Назад", "Вернуться к предыдущему меню")
	mockDisplay.EXPECT().MenuChoice()
}

// TestIsAuthError проверяет функцию определения ошибок авторизации
func TestIsAuthError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "ErrUnauthorized",
			err:      models.ErrUnauthorized,
			expected: true,
		},
		{
			name:     "wrapped ErrUnauthorized",
			err:      fmt.Errorf("wrapped: %w", models.ErrUnauthorized),
			expected: true,
		},
		{
			name:     "regular error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "network error",
			err:      errors.New("network timeout"),
			expected: false,
		},
		{
			name:     "authentication required (old format)",
			err:      errors.New("authentication required"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAuthError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestShell_IntegrationAuthError проверяет интеграцию обработки ошибок авторизации
// с реальными сообщениями от API сервера
func TestShell_IntegrationAuthError(t *testing.T) {
	// Тестируем ErrUnauthorized
	t.Run("ErrUnauthorized", func(t *testing.T) {
		err := models.ErrUnauthorized
		assert.True(t, isAuthError(err), "ErrUnauthorized должна быть определена как ошибка авторизации")
	})

	// Тестируем обернутую ErrUnauthorized
	t.Run("wrapped_ErrUnauthorized", func(t *testing.T) {
		err := fmt.Errorf("wrapped: %w", models.ErrUnauthorized)
		assert.True(t, isAuthError(err), "Обернутая ErrUnauthorized должна быть определена как ошибка авторизации")
	})

	// Тестируем обычные ошибки, которые НЕ должны быть определены как ошибки авторизации
	regularErrors := []string{
		"network timeout",
		"connection refused",
		"file not found",
		"permission denied",
		"invalid input",
		"server error",
		"authentication required", // старый формат
		"unauthorized",            // старый формат
	}

	for _, errMsg := range regularErrors {
		t.Run(fmt.Sprintf("regular_error_%s", errMsg), func(t *testing.T) {
			err := errors.New(errMsg)
			assert.False(t, isAuthError(err), "Ошибка '%s' НЕ должна быть определена как ошибка авторизации", errMsg)
		})
	}
}

// TestShell_HandleAuthError проверяет обработку ошибок авторизации в shell
func TestShell_HandleAuthError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegistry := NewMockCommandRegistry(ctrl)
	mockInputHandler := NewMockInputHandler(ctrl)
	mockDisplay := NewMockDisplay(ctrl)

	shell := NewShell(mockRegistry, mockInputHandler, mockDisplay, "1.0.0", "2024-01-01")

	// Тест для handleGetDeleteCommand с ошибкой авторизации
	t.Run("handleGetDeleteCommand with auth error", func(t *testing.T) {
		authError := models.ErrUnauthorized

		mockInputHandler.EXPECT().GetSecretName().Return("test-secret", nil)
		mockRegistry.EXPECT().Execute(gomock.Any(), "get", "test-secret").Return(nil, authError)
		mockDisplay.EXPECT().AuthError()
		mockDisplay.EXPECT().Logout()

		shell.handleGetDeleteCommand("get")
	})

	// Тест для handleDataTypeCommand с ошибкой авторизации
	t.Run("handleDataTypeCommand with auth error", func(t *testing.T) {
		authError := models.ErrUnauthorized

		// Настраиваем меню для теста
		shell.menuManager.SetActionState("upload", "")
		shell.menuManager.SwitchToState(MenuStateDataType)

		mockInputHandler.EXPECT().GetTextData().Return(&models.TextSecretData{Name: "test", Text: "test"}, nil)
		mockRegistry.EXPECT().Execute(gomock.Any(), "upload", gomock.Any()).Return(nil, authError)
		mockDisplay.EXPECT().AuthError()
		mockDisplay.EXPECT().Logout()

		shell.handleDataTypeCommand("text")
	})

	// Тест для handleVersionCommand с ошибкой авторизации
	t.Run("handleVersionCommand with auth error", func(t *testing.T) {
		authError := models.ErrUnauthorized

		mockRegistry.EXPECT().Execute(gomock.Any(), "version", nil).Return(nil, authError)
		mockDisplay.EXPECT().AuthError()
		mockDisplay.EXPECT().Logout()

		shell.handleVersionCommand()
	})

	// Тест для handleUpdateCommand с ошибкой авторизации при получении данных секрета
	t.Run("handleUpdateCommand with auth error on show", func(t *testing.T) {
		authError := models.ErrUnauthorized

		mockInputHandler.EXPECT().GetSecretName().Return("test-secret", nil)
		mockRegistry.EXPECT().Execute(gomock.Any(), "show", "test-secret").Return(nil, authError)
		mockDisplay.EXPECT().AuthError()
		mockDisplay.EXPECT().Logout()

		shell.handleUpdateCommand()
	})

	// Тест для handleUpdateCommand с ошибкой авторизации при обновлении
	t.Run("handleUpdateCommand with auth error on update", func(t *testing.T) {
		authError := models.ErrUnauthorized

		showResult := &cmd.ShowSecretResult{
			Type:     models.SecretTypeText,
			Metadata: "test metadata",
			Version:  1,
		}

		mockInputHandler.EXPECT().GetSecretName().Return("test-secret", nil)
		mockRegistry.EXPECT().Execute(gomock.Any(), "show", "test-secret").Return(showResult, nil)
		mockDisplay.EXPECT().DisplaySecretInfo("test-secret", models.SecretTypeText, "test metadata", 1)
		mockInputHandler.EXPECT().PromptEnterNewData()
		mockInputHandler.EXPECT().GetUpdatedTextData("test-secret").Return(&models.TextSecretData{Name: "test", Text: "updated"}, nil)
		mockRegistry.EXPECT().Execute(gomock.Any(), "update", gomock.Any()).Return(nil, authError)
		mockDisplay.EXPECT().AuthError()
		mockDisplay.EXPECT().Logout()

		shell.handleUpdateCommand()
	})
}
