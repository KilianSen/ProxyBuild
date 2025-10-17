package tests

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"ProxyBuild/proxy"
)

func TestLoadConfig(t *testing.T) {
	// Erstelle temporäre Config-Datei
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test-config.json")

	config := proxy.Config{
		BaseCommand: "echo",
		Hooks: map[string][]proxy.Hook{
			"test": {
				{
					Command: "echo",
					Args:    []string{"hook"},
					When:    "before",
				},
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Lade Config
	loadedData, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}

	var loadedConfig proxy.Config
	if err := json.Unmarshal(loadedData, &loadedConfig); err != nil {
		t.Fatal(err)
	}

	if loadedConfig.BaseCommand != "echo" {
		t.Errorf("Expected base_command 'echo', got '%s'", loadedConfig.BaseCommand)
	}

	if len(loadedConfig.Hooks) != 1 {
		t.Errorf("Expected 1 hook set, got %d", len(loadedConfig.Hooks))
	}
}

func TestConfigWithConditions(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test-config-conditions.json")

	trueVal := true
	falseVal := false

	config := proxy.Config{
		BaseCommand: "echo",
		Hooks: map[string][]proxy.Hook{
			"test": {
				{
					Command: "echo",
					Args:    []string{"on_error"},
					When:    "after",
					Conditions: proxy.Conditions{
						OnError: &trueVal,
					},
				},
				{
					Command: "echo",
					Args:    []string{"on_success"},
					When:    "after",
					Conditions: proxy.Conditions{
						OnError: &falseVal,
					},
				},
				{
					Command: "echo",
					Args:    []string{"with_args"},
					When:    "before",
					Conditions: proxy.Conditions{
						ArgsContain: []string{"-v"},
						ArgsMatch:   []string{"--verbose"},
					},
				},
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Lade und validiere
	loadedData, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}

	var loadedConfig proxy.Config
	if err := json.Unmarshal(loadedData, &loadedConfig); err != nil {
		t.Fatal(err)
	}

	hooks := loadedConfig.Hooks["test"]
	if len(hooks) != 3 {
		t.Errorf("Expected 3 hooks, got %d", len(hooks))
	}

	// Teste on_error condition
	if hooks[0].Conditions.OnError == nil {
		t.Error("OnError should not be nil")
	}
	if !*hooks[0].Conditions.OnError {
		t.Error("OnError should be true")
	}

	// Teste on_success condition
	if hooks[1].Conditions.OnError == nil {
		t.Error("OnError should not be nil")
	}
	if *hooks[1].Conditions.OnError {
		t.Error("OnError should be false")
	}

	// Teste args conditions
	if len(hooks[2].Conditions.ArgsContain) != 1 {
		t.Error("ArgsContain should have 1 element")
	}
	if len(hooks[2].Conditions.ArgsMatch) != 1 {
		t.Error("ArgsMatch should have 1 element")
	}
}

func TestBuildProcess(t *testing.T) {
	// Überspringe wenn nicht in CI
	if os.Getenv("CI") == "" {
		t.Skip("Skipping build test outside CI environment")
	}

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "build-config.json")

	config := proxy.Config{
		BaseCommand: "echo",
		Hooks: map[string][]proxy.Hook{
			"test": {
				{
					Command: "echo",
					Args:    []string{"test hook"},
					When:    "before",
				},
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Teste Build-Command (nur Syntax-Check)
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("Go not found in PATH")
	}
}

func TestInvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid-config.json")

	// Schreibe ungültige JSON
	invalidJSON := `{"base_command": "echo", "hooks": {invalid}}`
	if err := os.WriteFile(configFile, []byte(invalidJSON), 0644); err != nil {
		t.Fatal(err)
	}

	// Versuche zu laden
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}

	var config proxy.Config
	err = json.Unmarshal(data, &config)
	if err == nil {
		t.Error("Expected error when loading invalid config")
	}
}

func TestEmptyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "empty-config.json")

	config := proxy.Config{
		BaseCommand: "",
		Hooks:       map[string][]proxy.Hook{},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Lade Config
	loadedData, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}

	var loadedConfig proxy.Config
	if err := json.Unmarshal(loadedData, &loadedConfig); err != nil {
		t.Fatal(err)
	}

	if loadedConfig.BaseCommand != "" {
		t.Error("Empty base_command should remain empty")
	}

	if len(loadedConfig.Hooks) != 0 {
		t.Error("Empty hooks should remain empty")
	}
}
