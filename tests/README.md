# Tests

Dieses Verzeichnis enthält alle Tests für ProxyBuild.

## Struktur

- `proxy_test.go` - Unit-Tests für die Proxy-Logik und Bedingungen
- `integration_test.go` - Integrationstests für Config-Laden und Build-Prozess

## Tests lokal ausführen

### Alle Tests ausführen

```bash
cd tests
go test -v ./...
```

### Tests mit Coverage

```bash
cd tests
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Tests mit Race Detector

```bash
cd tests
go test -race -v ./...
```

### Spezifische Tests ausführen

```bash
cd tests
go test -v -run TestShouldExecuteHook_OnErrorTrue
```

## Test-Kategorien

### Unit Tests (proxy_test.go)

- **Bedingungslogik**: Testet `shouldExecuteHook` mit verschiedenen Bedingungen
  - `TestShouldExecuteHook_NoConditions` - Hooks ohne Bedingungen
  - `TestShouldExecuteHook_OnErrorTrue` - Error-Bedingung (true)
  - `TestShouldExecuteHook_OnErrorFalse` - Success-Bedingung (false)
  - `TestShouldExecuteHook_ArgsContain` - Argument-Contains-Bedingung
  - `TestShouldExecuteHook_ArgsMatch` - Argument-Match-Bedingung
  - `TestShouldExecuteHook_CombinedConditions` - Kombinierte Bedingungen

### Integration Tests (integration_test.go)

- **Config-Handling**: Testet JSON-Serialisierung/Deserialisierung
  - `TestLoadConfig` - Laden einer Konfigurationsdatei
  - `TestConfigWithConditions` - Config mit Bedingungen
  - `TestInvalidConfig` - Fehlerbehandlung bei ungültiger Config
  - `TestEmptyConfig` - Leere Konfiguration

- **Build-Prozess**: Testet die Build-Funktionalität
  - `TestBuildProcess` - Build eines Proxy-Executables (nur in CI)

## CI/CD

Die Tests werden automatisch bei jedem Push und Pull Request über GitHub Actions ausgeführt. Siehe `.github/workflows/test.yml` für Details.

### Test-Matrix

Die Tests werden auf folgenden Plattformen ausgeführt:
- **Betriebssysteme**: Ubuntu, macOS, Windows
- **Go-Versionen**: 1.21, 1.22, 1.23

### Coverage

Coverage-Reports werden automatisch zu Codecov hochgeladen.

## Test-Daten

Tests verwenden temporäre Verzeichnisse (`t.TempDir()`) für Test-Dateien, die automatisch nach jedem Test aufgeräumt werden.

## Neue Tests hinzufügen

1. Erstelle eine neue Test-Funktion mit Prefix `Test`
2. Verwende `t.TempDir()` für temporäre Dateien
3. Nutze Table-Driven Tests für ähnliche Testfälle
4. Dokumentiere komplexe Tests mit Kommentaren

Beispiel:

```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := NewFeature(tt.input)
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

