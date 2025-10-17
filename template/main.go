package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"ProxyBuild/proxy"
)

//go:embed config.json
var embeddedConfig []byte

func main() {
	var config proxy.Config
	if err := json.Unmarshal(embeddedConfig, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Fehler beim Laden der Konfiguration: %v\n", err)
		os.Exit(1)
	}

	if err := proxy.Run(&config, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
		os.Exit(1)
	}
}
