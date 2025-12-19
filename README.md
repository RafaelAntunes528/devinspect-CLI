# DevInspect

**DevInspect** is a cross-platform CLI tool written in Go that analyzes your project structure, counts files and lines, detects programming languages used, and generates a project quality score. Designed for developers, teams, and CI pipelines, it works on Windows, Linux, and macOS.

---

## Features

- Recursively scans all files and folders, including deeply nested directories
- Counts total files and lines of code
- Detects languages based on file extensions (Go, JavaScript, TypeScript, Python)
- Checks for presence of `README.md` and `.gitignore`
- Computes a project score based on structure and quality
- Supports `.devinspectignore` to skip folders
- Colored CLI output with tables
- JSON output for automation and CI
- Fully concurrent scanning for speed on large projects
- Single-binary cross-platform builds

---

## Installation

### Windows

1. Download the latest binary from [Releases](#) or build locally:

```powershell
go build -o devinspect.exe
````

2. Move the binary to a folder in your PATH (e.g., `C:\Tools`):

```powershell
move devinspect.exe C:\Tools\
```

3. Open a new terminal and test:

```powershell
devinspect scan .
```

---

### Linux / macOS

1. Build the binary:

```bash
go build -o devinspect
```

2. Move it to a folder in your PATH, e.g., `~/.local/bin`:

```bash
mv devinspect ~/.local/bin/
```

3. Test:

```bash
devinspect scan .
```

---

## Usage

```bash
# Scan the current directory
devinspect scan .

# Scan a specific directory
devinspect scan /path/to/project

# Output JSON report
devinspect scan /path/to/project --json

```

---

## `.devinspectignore`

Create a `.devinspectignore` file in your project root to skip folders:

```
node_modules
dist
build
vendor
```
