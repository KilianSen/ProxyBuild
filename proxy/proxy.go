package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Config definiert die Konfiguration für Command-Hooks
type Config struct {
	BaseCommand string            `json:"base_command"`
	Hooks       map[string][]Hook `json:"hooks"`
}

// Hook definiert einen Hook, der bei einem bestimmten Sub-Command ausgeführt wird
type Hook struct {
	Command    string     `json:"command"`
	Args       []string   `json:"args"`
	When       string     `json:"when"`       // "before" oder "after"
	Conditions Conditions `json:"conditions"` // Optionale Bedingungen
}

// Conditions definiert Bedingungen, unter denen ein Hook ausgeführt wird
type Conditions struct {
	OnError     *bool    `json:"on_error"`     // Nur bei Fehler (true) oder nur ohne Fehler (false)
	ArgsContain []string `json:"args_contain"` // Hook nur ausführen, wenn diese Strings in den Args enthalten sind
	ArgsMatch   []string `json:"args_match"`   // Hook nur ausführen, wenn Args exakt übereinstimmen
}

// Run führt den Proxy mit der gegebenen Konfiguration aus
func Run(config *Config, args []string) error {
	// Bestimme den Sub-Command (erstes Argument)
	subCommand := ""
	if len(args) > 0 {
		subCommand = args[0]
	}

	// Führe "before" Hooks aus
	if hooks, exists := config.Hooks[subCommand]; exists {
		for _, hook := range hooks {
			if hook.When == "before" && ShouldExecuteHook(hook, args, false) {
				if err := executeHook(hook); err != nil {
					return fmt.Errorf("Fehler beim Ausführen des before-Hooks: %w", err)
				}
			}
		}
	}

	// Führe den Basis-Command aus
	cmd := exec.Command(config.BaseCommand, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	baseCommandErr := cmd.Run()
	hadError := baseCommandErr != nil

	if hadError {
		// Fehler des Basis-Commands, aber trotzdem "after" Hooks ausführen
		_, err := fmt.Fprintf(os.Stderr, "Basis-Command fehlgeschlagen: %v\n", baseCommandErr)
		if err != nil {
			return err
		}
	}

	// Führe "after" Hooks aus
	if hooks, exists := config.Hooks[subCommand]; exists {
		for _, hook := range hooks {
			if hook.When == "after" && ShouldExecuteHook(hook, args, hadError) {
				if err := executeHook(hook); err != nil {
					return fmt.Errorf("Fehler beim Ausführen des after-Hooks: %w", err)
				}
			}
		}
	}

	return nil
}

// ShouldExecuteHook überprüft, ob ein Hook ausgeführt werden soll basierend auf den Bedingungen
func ShouldExecuteHook(hook Hook, args []string, hadError bool) bool {
	// Überprüfe OnError-Bedingung
	if hook.Conditions.OnError != nil {
		if *hook.Conditions.OnError != hadError {
			return false
		}
	}

	// Überprüfe ArgsContain-Bedingung
	if len(hook.Conditions.ArgsContain) > 0 {
		for _, substr := range hook.Conditions.ArgsContain {
			found := false
			for _, arg := range args {
				if strings.Contains(arg, substr) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	// Überprüfe ArgsMatch-Bedingung
	if len(hook.Conditions.ArgsMatch) > 0 {
		for _, matchArg := range hook.Conditions.ArgsMatch {
			found := false
			for _, arg := range args {
				if arg == matchArg {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

func executeHook(hook Hook) error {
	fmt.Printf("→ Hook: %s %s\n", hook.Command, strings.Join(hook.Args, " "))
	cmd := exec.Command(hook.Command, hook.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
