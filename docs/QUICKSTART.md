# Quick Start Guide for Using ProxyBuild in Your Repository

This guide shows you how to use ProxyBuild as a GitHub Action in your own repository.

## Step 1: Create Your Proxy Configuration

Create a file called `proxy-config.json` in your repository root:

```json
{
  "base_command": "docker-compose",
  "hooks": {
    "up": [
      {
        "command": "echo",
        "args": ["🚀 Starting containers..."],
        "when": "before"
      },
      {
        "command": "echo",
        "args": ["✅ Containers started!"],
        "when": "after",
        "conditions": {
          "on_error": false
        }
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

## Step 2: Create a GitHub Workflow

Create `.github/workflows/build-proxy.yml` in your repository:

```yaml
name: Build Proxy

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Build proxy executable
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: 'proxy-config.json'
          
      # Executable is automatically uploaded as an artifact
      # Download it from the Actions tab
```

## Step 3: Push and Run

```bash
git add proxy-config.json .github/workflows/build-proxy.yml
git commit -m "Add proxy configuration and build workflow"
git push
```

The workflow will run automatically, and you can download the built executable from the Actions tab.

## Advanced: Auto-Release on Tags

To automatically create releases with your proxy executable:

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
          
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: ${{ steps.build.outputs.executable-path }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Then create a release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Multi-Platform Builds

To build for multiple platforms:

```yaml
name: Multi-Platform Build

on: [push]

jobs:
  build:
    strategy:
      matrix:
        include:
          - os: linux
            arch: amd64
          - os: windows
            arch: amd64
          - os: darwin
            arch: arm64
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build proxy for ${{ matrix.os }}-${{ matrix.arch }}
        uses: YOUR_USERNAME/ProxyBuild@v1
        with:
          config-file: 'proxy-config.json'
          target-os: ${{ matrix.os }}
          target-arch: ${{ matrix.arch }}
          output-name: 'my-proxy-${{ matrix.os }}-${{ matrix.arch }}'
```

## That's It!

Your repository now automatically builds proxy executables. Download them from:
- The Actions tab (Artifacts section)
- Release page (if using auto-release)

## Need Help?

- See [full documentation](GITHUB_ACTION.md)
- Check [example configurations](../examples/example-config.json)
- Review [example workflows](../.github/workflows/example-usage.yml)
name: Test Action

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test-action:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Build ProxyBuild tool
        run: |
          go build -o ProxyBuild main.go
          
      - name: Test building a proxy
        run: |
          ./ProxyBuild -build example-config.json
          
      - name: Verify proxy was built
        run: |
          ls -la docker-compose-proxy
          chmod +x docker-compose-proxy
          
      - name: Test proxy execution
        run: |
          # The proxy should at least respond to help
          ./docker-compose-proxy --help || true
          
      - name: Upload test artifact
        uses: actions/upload-artifact@v4
        with:
          name: test-proxy-executable
          path: docker-compose-proxy
