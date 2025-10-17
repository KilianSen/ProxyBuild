# GitHub Action Usage

ProxyBuild kann als GitHub Action in anderen Repositories verwendet werden, um automatisch Command-Proxy-Executables zu bauen.

## Quick Start

Erstelle eine `.github/workflows/build-proxy.yml` Datei in deinem Repository:

```yaml
name: Build Proxy

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build proxy executable
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: 'proxy-config.json'
```

## Inputs

| Input | Beschreibung | Erforderlich | Standard |
|-------|-------------|--------------|----------|
| `config-file` | Pfad zur Proxy-Konfigurationsdatei (JSON) | Ja | - |
| `output-name` | Name des Output-Executables | Nein | `<base_command>-proxy` |
| `go-version` | Zu verwendende Go-Version | Nein | `1.24` |
| `upload-artifact` | Executable als Artifact hochladen | Nein | `true` |
| `target-os` | Ziel-Betriebssystem (linux, darwin, windows) | Nein | - |
| `target-arch` | Ziel-Architektur (amd64, arm64, 386) | Nein | - |

## Outputs

| Output | Beschreibung |
|--------|-------------|
| `executable-path` | Vollständiger Pfad zum gebauten Executable |
| `executable-name` | Name des gebauten Executables |

## Beispiele

### Basis-Verwendung

```yaml
- name: Build Docker Compose Proxy
  uses: YOUR_USERNAME/ProxyBuild@v1
  with:
    config-file: 'docker-compose-proxy.json'
```

### Mit Custom Output-Name

```yaml
- name: Build Custom Proxy
  uses: YOUR_USERNAME/ProxyBuild@v1
  with:
    config-file: 'my-config.json'
    output-name: 'my-custom-tool'
```

### Mehrere Proxies bauen

```yaml
jobs:
  build-proxies:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        proxy:
          - config: 'docker-compose-config.json'
            name: 'docker-compose-proxy'
          - config: 'kubectl-config.json'
            name: 'kubectl-proxy'
    steps:
      - uses: actions/checkout@v4
      
      - name: Build ${{ matrix.proxy.name }}
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: ${{ matrix.proxy.config }}
          output-name: ${{ matrix.proxy.name }}
```

### Build und Release bei Tags

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build proxy
        id: build
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: 'proxy-config.json'
          output-name: 'my-proxy'
          
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: ${{ steps.build.outputs.executable-path }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Multi-Platform Builds

```yaml
jobs:
  build:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      
      - name: Build proxy
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: 'proxy-config.json'
          
      - name: Upload to release
        # ... upload logic
```

### Cross-Compilation für alle Plattformen (von einem Runner)

```yaml
jobs:
  build-all-platforms:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - os: linux
            arch: amd64
            name: linux-amd64
          - os: linux
            arch: arm64
            name: linux-arm64
          - os: darwin
            arch: amd64
            name: macos-amd64
          - os: darwin
            arch: arm64
            name: macos-arm64
          - os: windows
            arch: amd64
            name: windows-amd64
    steps:
      - uses: actions/checkout@v4
      
      - name: Build for ${{ matrix.name }}
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: 'proxy-config.json'
          target-os: ${{ matrix.os }}
          target-arch: ${{ matrix.arch }}
          output-name: 'my-proxy-${{ matrix.name }}'
```

### Testen des gebauten Executables

```yaml
- name: Build proxy
  id: build
  uses: YOUR_USERNAME/ProxyBuild@v1
  with:
    config-file: 'proxy-config.json'

- name: Test executable
  run: |
    chmod +x ${{ steps.build.outputs.executable-path }}
    ${{ steps.build.outputs.executable-path }} --version
```

### Ohne Artifact Upload

```yaml
- name: Build proxy
  uses: YOUR_USERNAME/ProxyBuild@v1
  with:
    config-file: 'proxy-config.json'
    upload-artifact: 'false'
```

## Setup für dein Repository

1. **Erstelle eine Konfigurationsdatei** in deinem Repository (z.B. `proxy-config.json`):
   ```json
   {
     "base_command": "docker-compose",
     "hooks": {
       "up": [
         {
           "command": "echo",
           "args": ["Starting containers..."],
           "when": "before"
         }
       ]
     }
   }
   ```

2. **Erstelle einen Workflow** (`.github/workflows/build-proxy.yml`):
   ```yaml
   name: Build Proxy
   
   on: [push]
   
   jobs:
     build:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v4
         - uses: YOUR_USERNAME/ProxyBuild@v1
           with:
             config-file: 'proxy-config.json'
   ```

3. **Commit und Push** - Die Action wird automatisch ausgeführt

4. **Download das Artifact** von der Actions-Seite oder verwende es in nachfolgenden Steps

## Tipps

- **Versionierung**: Verwende spezifische Versionen (`@v1.0.0`) statt `@v1` für reproduzierbare Builds
- **Caching**: Die Action cached automatisch Go-Dependencies für schnellere Builds
- **Matrix Builds**: Verwende Matrix-Builds für mehrere Proxies gleichzeitig
- **Security**: Konfigurationsdateien sollten keine sensitiven Daten enthalten

## Troubleshooting

### Fehler: "config file not found"
- Stelle sicher, dass der Pfad zur Config-Datei relativ zum Repository-Root ist
- Prüfe, ob die Datei im Repository committed ist

### Fehler: "invalid JSON"
- Validiere deine JSON-Konfiguration mit einem JSON-Validator
- Prüfe auf Syntaxfehler (fehlende Kommas, Anführungszeichen, etc.)

### Executable funktioniert nicht auf anderen Systemen
- Baue separate Executables für verschiedene Plattformen (siehe Multi-Platform Beispiel)
- Verwende die Matrix-Strategy mit verschiedenen Runner-OS

## Weitere Ressourcen

- [Vollständige Dokumentation](../README.md)
- [Beispiel-Konfigurationen](../examples/example-config.json)
- [GitHub Actions Dokumentation](https://docs.github.com/en/actions)
# Example workflow showing how to use ProxyBuild Action
# This file serves as documentation and can be copied to other repositories

name: Build Proxy Example

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  workflow_dispatch:

jobs:
  build-proxy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Build proxy executable
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: 'proxy-config.json'
          output-name: 'my-custom-proxy'
          
      # The executable is now available in the workspace
      - name: Test the proxy
        run: |
          ./my-custom-proxy --help
          
  build-multiple-proxies:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        config:
          - file: 'docker-compose-config.json'
            name: 'docker-compose-proxy'
          - file: 'git-config.json'
            name: 'git-proxy'
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Build ${{ matrix.config.name }}
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: ${{ matrix.config.file }}
          output-name: ${{ matrix.config.name }}
          
  build-and-release:
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/')
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Build proxy
        id: build
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: 'proxy-config.json'
          
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: ${{ steps.build.outputs.executable-path }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

