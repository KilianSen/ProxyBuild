package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"ProxyBuild/proxy"
)

type BuildOptions struct {
	ConfigFile string
	GOOS       string
	GOARCH     string
	OutputName string
}

func main() {
	buildCmd := flag.String("build", "", "Erstellt ein neues ausführbares Programm mit der angegebenen Konfigurationsdatei")
	configFile := flag.String("config", "", "Konfigurationsdatei für den Proxy-Modus")
	goos := flag.String("os", "", "Ziel-Betriebssystem für Cross-Compilation (z.B. linux, darwin, windows)")
	goarch := flag.String("arch", "", "Ziel-Architektur für Cross-Compilation (z.B. amd64, arm64)")
	outputName := flag.String("output", "", "Name des Output-Executables (optional)")
	flag.Parse()

	if *buildCmd != "" {
		// Build-Modus: Erstelle ein neues ausführbares Programm
		buildOpts := BuildOptions{
			ConfigFile: *buildCmd,
			GOOS:       *goos,
			GOARCH:     *goarch,
			OutputName: *outputName,
		}
		if err := buildExecutable(buildOpts); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Fehler beim Erstellen: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Println("Executable erfolgreich erstellt!")
		return
	}

	if *configFile != "" {
		// Proxy-Modus mit Konfigurationsdatei
		config, err := loadConfig(*configFile)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Fehler beim Laden der Konfiguration: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}

		if err := proxy.Run(config, flag.Args()); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		return
	}

	// Zeige Hilfe an
	fmt.Println("ProxyBuild - Command Proxy Tool")
	fmt.Println("\nVerwendung:")
	fmt.Println("  ProxyBuild -config <config.json> [args...]  - Führt Proxy mit Konfiguration aus")
	fmt.Println("  ProxyBuild -build <config.json>             - Erstellt ein neues Executable")
	fmt.Println("\nBuild-Optionen:")
	fmt.Println("  -os <os>       Ziel-Betriebssystem (linux, darwin, windows)")
	fmt.Println("  -arch <arch>   Ziel-Architektur (amd64, arm64, 386)")
	fmt.Println("  -output <name> Name des Output-Executables")
	fmt.Println("\nBeispiele:")
	fmt.Println("  ProxyBuild -build config.json")
	fmt.Println("  ProxyBuild -build config.json -os linux -arch amd64")
	fmt.Println("  ProxyBuild -build config.json -os windows -arch amd64 -output my-tool.exe")
	fmt.Println("  ProxyBuild -build config.json -os darwin -arch arm64 -output my-tool-mac")
}

func loadConfig(filename string) (*proxy.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config proxy.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func buildExecutable(opts BuildOptions) error {
	// Lade Konfiguration
	config, err := loadConfig(opts.ConfigFile)
	if err != nil {
		return fmt.Errorf("Fehler beim Laden der Konfiguration: %w", err)
	}

	outputName := opts.OutputName
	if outputName == "" {
		outputName = filepath.Base(config.BaseCommand) + "-proxy"
		// Füge .exe für Windows hinzu
		if opts.GOOS == "windows" {
			outputName += ".exe"
		}
	}

	// Erstelle temporäres Build-Verzeichnis
	buildDir := "tmp_build"
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return err
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Fehler beim Aufräumen des Build-Verzeichnisses: %v\n", err)
			if err != nil {
				return
			}
		}
	}(buildDir)

	// Kopiere die Konfigurationsdatei ins Build-Verzeichnis
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Resolve applicable build env vars to config, before embedding config in build step
	envVars := os.Environ()
	if len(envVars) >= 1 {
		for _, envVar := range envVars {
			key := strings.Split(envVar, "=")[0]
			val := strings.Split(envVar, "=")[1]
			// Windows Env Syntax: %key%
			windowsSyntax := fmt.Sprintf("%%%s%%", key)

			// Unix Env Syntax: $key
			unixSyntax := fmt.Sprintf("$%s", key)

			configData = bytes.Replace(configData, []byte(unixSyntax), []byte(val), -1)
			configData = bytes.Replace(configData, []byte(windowsSyntax), []byte(val), -1)
		}
	}

	if err := os.WriteFile(filepath.Join(buildDir, "config.json"), configData, 0644); err != nil {
		return err
	}

	// Kopiere das Template main.go
	templatePath := "template/main.go"
	if err := copyFile(templatePath, filepath.Join(buildDir, "main.go")); err != nil {
		return fmt.Errorf("Fehler beim Kopieren des Templates: %w", err)
	}

	// Erstelle go.mod im Build-Verzeichnis
	goModContent := fmt.Sprintf(`module generated-proxy

go 1.24

require ProxyBuild/proxy v0.0.0

replace ProxyBuild/proxy => %s
`, filepath.Join("..", "proxy"))

	if err := os.WriteFile(filepath.Join(buildDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		return err
	}

	// Kompiliere das Programm
	fmt.Printf("Kompiliere %s...\n", outputName)

	// Zeige Cross-Compilation Info
	if opts.GOOS != "" || opts.GOARCH != "" {
		targetOS := opts.GOOS
		if targetOS == "" {
			targetOS = "current"
		}
		targetArch := opts.GOARCH
		if targetArch == "" {
			targetArch = "current"
		}
		fmt.Printf("Cross-Compiling für OS=%s, ARCH=%s\n", targetOS, targetArch)
	}

	buildCmd := exec.Command("go", "build", "-o", filepath.Join("..", outputName), ".")
	buildCmd.Dir = buildDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	// Setze Cross-Compilation-Flags, falls angegeben
	buildCmd.Env = os.Environ()
	if opts.GOOS != "" {
		buildCmd.Env = append(buildCmd.Env, "GOOS="+opts.GOOS)
	}
	if opts.GOARCH != "" {
		buildCmd.Env = append(buildCmd.Env, "GOARCH="+opts.GOARCH)
	}

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("Kompilierung fehlgeschlagen: %w", err)
	}

	fmt.Printf("✓ Executable erstellt: %s\n", outputName)
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
