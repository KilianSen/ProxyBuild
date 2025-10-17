# ProxyBuild - Command Proxy Tool

[![Tests](https://github.com/YOUR_USERNAME/ProxyBuild/workflows/Tests/badge.svg)](https://github.com/YOUR_USERNAME/ProxyBuild/actions)
[![codecov](https://codecov.io/gh/YOUR_USERNAME/ProxyBuild/branch/main/graph/badge.svg)](https://codecov.io/gh/YOUR_USERNAME/ProxyBuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/ProxyBuild)](https://goreportcard.com/report/github.com/YOUR_USERNAME/ProxyBuild)

Ein flexibles Tool, das als Proxy für beliebige Commands fungiert und konfigurierbare Hooks vor/nach Sub-Commands ausführt.

## Features

- **Proxy-Modus**: Leitet alle Aufrufe an ein Basis-Command weiter
- **Hooks**: Führt zusätzliche Commands vor oder nach bestimmten Sub-Commands aus
- **Build-Modus**: Erstellt ein eigenständiges Executable mit eingebetteter Konfiguration
- **GitHub Action**: Automatischer Build von Proxies in CI/CD Pipelines

## Installation

### Lokale Installation

```bash
go build -o ProxyBuild main.go
```

### Als GitHub Action

Füge diese Action zu deinem Repository hinzu:

```yaml
- uses: YOUR_USERNAME/ProxyBuild@v1
  with:
    config-file: 'proxy-config.json'
```

Siehe [docs/GITHUB_ACTION.md](docs/GITHUB_ACTION.md) für detaillierte Informationen.

## Verwendung

### 1. Proxy-Modus mit Konfigurationsdatei

```bash
./ProxyBuild -config config.json [sub-command] [args...]
```

Beispiel:
```bash
./ProxyBuild -config example-config.json up -d
```

### 2. Build-Modus

Erstellt ein eigenständiges Executable mit eingebetteter Konfiguration:

```bash
./ProxyBuild -build config.json
```

Dies erstellt z.B. `docker-compose-proxy`, das direkt verwendet werden kann:

```bash
./docker-compose-proxy up -d
```

### 3. Cross-Compilation

Baue Executables für andere Plattformen:

```bash
# Für Linux bauen
./ProxyBuild -build config.json -os linux -arch amd64

# Für Windows bauen
./ProxyBuild -build config.json -os windows -arch amd64

# Für macOS Apple Silicon bauen
./ProxyBuild -build config.json -os darwin -arch arm64 -output my-tool
```

**Unterstützte Plattformen:**
- **OS:** linux, darwin (macOS), windows, freebsd, openbsd
- **Arch:** amd64, arm64, 386, arm

Siehe [Cross-Compilation Guide](docs/CROSS_COMPILATION.md) für detaillierte Informationen.

## Konfiguration

Die Konfigurationsdatei ist eine JSON-Datei mit folgendem Format:

```json
{
  "base_command": "docker-compose",
  "hooks": {
    "up": [
      {
        "command": "echo",
        "args": ["Starting services..."],
        "when": "before"
      },
      {
        "command": "echo",
        "args": ["Services started!"],
        "when": "after"
      }
    ],
    "down": [
      {
        "command": "echo",
        "args": ["Cleaning up..."],
        "when": "after"
      }
    ]
  }
}
```

### Konfigurationsfelder

- **base_command**: Das Command, das als Proxy verwendet wird (z.B. `docker-compose`, `git`, `kubectl`)
- **hooks**: Map von Sub-Commands zu Hook-Arrays
  - **command**: Das auszuführende Command
  - **args**: Array von Argumenten für das Command
  - **when**: Wann der Hook ausgeführt wird (`"before"` oder `"after"`)
  - **conditions** (optional): Bedingungen, unter denen der Hook ausgeführt wird
    - **on_error**: `true` = nur bei Fehler ausführen, `false` = nur bei Erfolg ausführen, `null` = immer ausführen
    - **args_contain**: Array von Strings - Hook wird nur ausgeführt, wenn alle diese Strings in den Argumenten enthalten sind
    - **args_match**: Array von Strings - Hook wird nur ausgeführt, wenn alle diese Strings exakt in den Argumenten vorkommen

## Beispiele

### Docker Compose mit Hooks

```json
{
  "base_command": "docker-compose",
  "hooks": {
    "up": [
      {
        "command": "echo",
        "args": ["🚀 Starting Docker containers..."],
        "when": "before"
      },
      {
        "command": "notify-send",
        "args": ["Docker", "Containers sind gestartet!"],
        "when": "after"
      }
    ],
    "down": [
      {
        "command": "docker",
        "args": ["system", "prune", "-f"],
        "when": "after"
      }
    ]
  }
}
```

### Git mit Pre/Post Commit Hooks

```json
{
  "base_command": "git",
  "hooks": {
    "commit": [
      {
        "command": "npm",
        "args": ["run", "lint"],
        "when": "before"
      },
      {
        "command": "npm",
        "args": ["test"],
        "when": "before"
      }
    ],
    "push": [
      {
        "command": "echo",
        "args": ["Pushing to remote..."],
        "when": "before"
      }
    ]
  }
}
```

### Mit Bedingungen (Conditions)

```json
{
  "base_command": "docker-compose",
  "hooks": {
    "up": [
      {
        "command": "echo",
        "args": ["🚀 Starting services..."],
        "when": "before"
      },
      {
        "command": "echo",
        "args": ["✅ Services started successfully!"],
        "when": "after",
        "conditions": {
          "on_error": false
        }
      },
      {
        "command": "echo",
        "args": ["❌ Failed to start services! Rolling back..."],
        "when": "after",
        "conditions": {
          "on_error": true
        }
      },
      {
        "command": "echo",
        "args": ["Running in detached mode"],
        "when": "after",
        "conditions": {
          "args_contain": ["-d"]
        }
      }
    ],
    "down": [
      {
        "command": "docker",
        "args": ["system", "prune", "-f"],
        "when": "after",
        "conditions": {
          "on_error": false,
          "args_contain": ["--volumes"]
        }
      }
    ]
  }
}
```

## Use Cases

1. **Automatisches Cleanup**: Räume nach `docker-compose down` automatisch auf
2. **Notifications**: Sende Benachrichtigungen nach langen Commands
3. **Pre-Checks**: Führe Tests oder Linting vor Git-Commits aus
4. **Logging**: Protokolliere Command-Ausführungen
5. **Wrapper**: Erstelle custom Wrapper für beliebige CLI-Tools

## Development

### Tests ausführen

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

Siehe [tests/README.md](tests/README.md) für detaillierte Informationen zu den Tests.

## Lizenz

MIT
