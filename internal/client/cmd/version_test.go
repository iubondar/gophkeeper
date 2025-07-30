package cmd

import (
	"context"
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	version := "1.0.0"
	buildTime := "2024-01-01T00:00:00Z"
	cmd := NewVersionCommand(version, buildTime)

	if cmd == nil {
		t.Fatal("NewVersionCommand returned nil")
	}

	if cmd.version != version {
		t.Errorf("Expected version '%s', got '%s'", version, cmd.version)
	}

	if cmd.buildTime != buildTime {
		t.Errorf("Expected buildTime '%s', got '%s'", buildTime, cmd.buildTime)
	}
}

func TestVersionCommand_Execute(t *testing.T) {
	ctx := context.Background()
	version := "1.0.0"
	buildTime := "2024-01-01T00:00:00Z"
	cmd := NewVersionCommand(version, buildTime)

	result, err := cmd.Execute(ctx, nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	versionResult, ok := result.(*VersionResult)
	if !ok {
		t.Fatal("Expected result to be *VersionResult")
	}

	if versionResult.Version != version {
		t.Errorf("Expected version '%s', got '%s'", version, versionResult.Version)
	}

	if versionResult.BuildTime != buildTime {
		t.Errorf("Expected buildTime '%s', got '%s'", buildTime, versionResult.BuildTime)
	}
}

func TestVersionCommand_GetName(t *testing.T) {
	cmd := NewVersionCommand("1.0.0", "2024-01-01T00:00:00Z")

	name := cmd.GetName()
	if name != "version" {
		t.Errorf("Expected name 'version', got '%s'", name)
	}
}

func TestVersionCommand_GetDescription(t *testing.T) {
	cmd := NewVersionCommand("1.0.0", "2024-01-01T00:00:00Z")

	description := cmd.GetDescription()
	expected := "Показать версию и дату сборки клиента"
	if description != expected {
		t.Errorf("Expected description '%s', got '%s'", expected, description)
	}
}
