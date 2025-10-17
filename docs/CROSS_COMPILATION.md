# Cross-Compilation Guide

ProxyBuild unterstützt Cross-Compilation für verschiedene Betriebssysteme und Architekturen.

## Lokale Verwendung

### Basis-Syntax

```bash
ProxyBuild -build <config.json> -os <target-os> -arch <target-arch>
```

### Unterstützte Plattformen

#### Betriebssysteme (GOOS)
- `linux` - Linux
- `darwin` - macOS
- `windows` - Windows
- `freebsd` - FreeBSD
- `openbsd` - OpenBSD

#### Architekturen (GOARCH)
- `amd64` - 64-bit x86
- `arm64` - 64-bit ARM (Apple Silicon, ARM servers)
- `386` - 32-bit x86
- `arm` - 32-bit ARM

### Beispiele

#### Für Linux bauen (auf macOS/Windows)
```bash
ProxyBuild -build config.json -os linux -arch amd64
```

#### Für Windows bauen (auf macOS/Linux)
```bash
ProxyBuild -build config.json -os windows -arch amd64 -output my-tool.exe
```

#### Für macOS Apple Silicon bauen
```bash
ProxyBuild -build config.json -os darwin -arch arm64
```

#### Für Raspberry Pi (ARM) bauen
```bash
ProxyBuild -build config.json -os linux -arch arm64
```

### Alle Plattformen auf einmal bauen

Bash-Script für Multi-Platform-Build:

```bash
#!/bin/bash

CONFIG="proxy-config.json"
NAME="my-proxy"

# Linux
ProxyBuild -build $CONFIG -os linux -arch amd64 -output "${NAME}-linux-amd64"
ProxyBuild -build $CONFIG -os linux -arch arm64 -output "${NAME}-linux-arm64"

# macOS
ProxyBuild -build $CONFIG -os darwin -arch amd64 -output "${NAME}-macos-intel"
ProxyBuild -build $CONFIG -os darwin -arch arm64 -output "${NAME}-macos-silicon"

# Windows
ProxyBuild -build $CONFIG -os windows -arch amd64 -output "${NAME}-windows.exe"

echo "✅ Alle Builds abgeschlossen!"
```

## GitHub Actions

### Einzelne Plattform

```yaml
- name: Build for Linux
  uses: YOUR_USERNAME/ProxyBuild@v1
  with:
    config-file: 'proxy-config.json'
    target-os: 'linux'
    target-arch: 'amd64'
```

### Mehrere Plattformen mit Matrix

```yaml
name: Multi-Platform Build

on: [push]

jobs:
  build:
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
            name: macos-intel
          - os: darwin
            arch: arm64
            name: macos-silicon
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
          
      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: ${{ matrix.name }}
          path: my-proxy-${{ matrix.name }}*
```

### Release mit allen Plattformen

```yaml
name: Release

on:
  push:
    tags: ['v*']

jobs:
  release:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - os: linux
            arch: amd64
          - os: darwin
            arch: arm64
          - os: windows
            arch: amd64
    steps:
      - uses: actions/checkout@v4
      
      - name: Build ${{ matrix.os }}-${{ matrix.arch }}
        id: build
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: 'proxy-config.json'
          target-os: ${{ matrix.os }}
          target-arch: ${{ matrix.arch }}
          upload-artifact: 'false'
          
      - name: Upload to release
        uses: softprops/action-gh-release@v1
        with:
          files: ${{ steps.build.outputs.executable-path }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Tipps

### Automatische Dateinamen
- Windows-Builds bekommen automatisch `.exe`-Extension
- Verwende `-output` für custom Namen

### Performance
- Cross-Compilation ist deutlich schneller als native Builds auf verschiedenen Runnern
- Ein einzelner Ubuntu-Runner kann für alle Plattformen bauen

### Testing
- Executables können nur auf der Zielplattform getestet werden
- Für umfassendes Testing: Verwende Matrix mit verschiedenen `runs-on`

### Größe optimieren

Für kleinere Executables:

```bash
# Mit zusätzlichen Build-Flags (erfordert Anpassung in main.go)
ProxyBuild -build config.json -os linux -arch amd64
```

Für noch kleinere Binaries kannst du in `main.go` die Build-Flags erweitern:

```go
buildCmd := exec.Command("go", "build", 
    "-ldflags", "-s -w",  // Strip debug info
    "-o", filepath.Join("..", outputName), 
    ".")
```

## Häufige Plattform-Kombinationen

### Server/Cloud
```bash
# AMD64 für die meisten Cloud-Provider
ProxyBuild -build config.json -os linux -arch amd64

# ARM64 für AWS Graviton, Oracle Cloud ARM
ProxyBuild -build config.json -os linux -arch arm64
```

### Desktop
```bash
# Windows
ProxyBuild -build config.json -os windows -arch amd64 -output tool.exe

# macOS Intel
ProxyBuild -build config.json -os darwin -arch amd64

# macOS Apple Silicon (M1/M2/M3)
ProxyBuild -build config.json -os darwin -arch arm64
```

### IoT/Embedded
```bash
# Raspberry Pi 3/4/5
ProxyBuild -build config.json -os linux -arch arm64

# Ältere ARM-Geräte
ProxyBuild -build config.json -os linux -arch arm
```

## Troubleshooting

### CGO-bezogene Fehler
Falls du CGO benötigst, ist Cross-Compilation komplexer. ProxyBuild verwendet reines Go und sollte keine Probleme haben.

### Fehlende Dependencies
Stelle sicher, dass alle Dependencies pure Go sind (keine C-Bibliotheken).

### Executable funktioniert nicht
- Prüfe die richtige OS/Arch-Kombination
- Teste auf der Zielplattform
- Für Windows: Vergiss nicht die `.exe`-Extension

