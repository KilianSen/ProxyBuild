package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Config definiert die Konfiguration für Command-Hooks
type Config struct {
	BaseCommand string            `json:"base_command"`
	Executor    Executor          `json:"executor"`
	Hooks       map[string][]Hook `json:"hooks"`
	EnvVars     map[string]string `json:"env_vars"`
}

type Executor string

const (
	ExecutorShell  Executor = "shell"
	ExecutorDirect Executor = "direct"
)

// Hook definiert einen Hook, der bei einem bestimmten Sub-Command ausgeführt wird
type Hook struct {
	Command    string     `json:"command"`
	Args       []string   `json:"args"`
	Executor   Executor   `json:"executor"`
	When       string     `json:"when"`       // "before" oder "after"
	Conditions Conditions `json:"conditions"` // Optionale Bedingungen
}

// Conditions definiert Bedingungen, unter denen ein Hook ausgeführt wird
type Conditions struct {
	OnError     *bool    `json:"on_error"`     // Nur bei Fehler (true) oder nur ohne Fehler (false)
	ArgsContain []string `json:"args_contain"` // Hook nur ausführen, wenn diese Strings in den Args enthalten sind
	ArgsMatch   []string `json:"args_match"`   // Hook nur ausführen, wenn Args exakt übereinstimmen
	OsMatch     []string `json:"os_match"`     // Hook nur ausführen, wenn OS partitive übereinstimmt
}

// Run führt den Proxy mit der gegebenen Konfiguration aus
func Run(config *Config, args []string) error {
	// Bestimme den Sub-Command (erstes Argument)
	subCommand := ""
	if len(args) > 0 {
		subCommand = args[0]
	}

	osString := runtime.GOOS

	// Führe "before" Hooks aus
	if hooks, exists := config.Hooks[subCommand]; exists {
		for _, hook := range hooks {
			if hook.When == "before" && ShouldExecuteHook(hook, args, false, osString) {
				if err := executeHook(hook); err != nil {
					return fmt.Errorf("fehler beim Ausführen des before-Hooks: %w", err)
				}
			}
		}
	}

	// Config EnvVars
	configEnv := config.EnvVars

	if configEnv == nil {
		configEnv = make(map[string]string)
	}

	for _, kV := range os.Environ() {
		key := strings.Split(kV, "=")[0]
		value := strings.Split(kV, "=")[1]
		configEnv[key] = value
	}

	var overloaded []string
	for key, value := range configEnv {
		overloaded = append(overloaded, fmt.Sprintf("%s=%s", key, value))
	}

	err := execute(config.BaseCommand, args, config.Executor, overloaded)

	// Führe "after" Hooks aus
	if hooks, exists := config.Hooks[subCommand]; exists {
		for _, hook := range hooks {
			if hook.When == "after" && ShouldExecuteHook(hook, args, err != nil, osString) {
				if err := executeHook(hook); err != nil {
					return fmt.Errorf("fehler beim Ausführen des after-Hooks: %w", err)
				}
			}
		}
	}

	return nil
}

// ShouldExecuteHook überprüft, ob ein Hook ausgeführt werden soll basierend auf den Bedingungen
func ShouldExecuteHook(hook Hook, args []string, hadError bool, os string) bool {
	// Check if one of the supplies OS's matches
	anyOsMatch := false
	if hook.Conditions.OsMatch != nil {
		for _, match := range hook.Conditions.OsMatch {
			if match == os {
				anyOsMatch = true
			}
		}
		if !anyOsMatch {
			return false
		}
	}

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
	return execute(hook.Command, hook.Args, hook.Executor, nil)
}

func execute(command string, args []string, executor Executor, env []string) error {

	if executor == "" {
		executor = ExecutorShell
	}

	if executor == ExecutorShell {
		shellCmd := ""
		var shellArgs []string
		if runtime.GOOS == "windows" {
			shellCmd = "cmd"
			shellArgs = []string{"/C", command}
			shellArgs = append(shellArgs, args...)
		} else {
			shellCmd = "/bin/sh"
			shellArgs = []string{"-c", command + " " + strings.Join(args, " ")}
		}
		cmd := exec.Command(shellCmd, shellArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if env != nil {
			cmd.Env = env
		}
		return cmd.Run()
	} else if executor == ExecutorDirect {
		cmd := exec.Command(command, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if env != nil {
			cmd.Env = env
		}
		return cmd.Run()
	}

	return fmt.Errorf("unbekannter Executor-Typ: %s", executor)
}
