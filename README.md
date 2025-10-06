# newestfiles

A command-line tool to find and list files by extension (or all files), sorted by modification time (newest first).

## Features

- Find files by extension(s) in current directory and subdirectories, or all files if no extensions specified
- Multiple sorting options:
  - By modification time (newest first - default)
  - By modification time (oldest first)
  - By file size (largest first)
  - By file size (smallest first)
- Support for multiple file extensions
- Case-insensitive extension matching
- Two output formats: plain text (default) and JSON
- Automatic extension normalization (adds dot if missing)

## Usage

```bash
# List all files (newest first)
newestfiles

# Plain text output with specific extensions (newest first)
newestfiles .go .txt .md

# Sort by oldest files first
newestfiles -o .go .txt .md

# Sort by largest files first
newestfiles -l .go .txt .md

# Sort by smallest files first
newestfiles -s .go .txt .md

# JSON output with size sorting
newestfiles -j -l .go .txt .md

# Extensions without dots are automatically normalized
newestfiles go txt md

# List all files with JSON output
newestfiles -j
```

### Command Line Options

- `-j`: Output in JSON format (default: plain text)
- `-o`: Sort by oldest files first (default: newest first)
- `-l`: Sort by largest files first
- `-s`: Sort by smallest files first

**Note:** Only one sorting option can be used at a time.

## Examples

### Plain Text Output
```bash
$ newestfiles .go .txt
main.go
utils.go
README.txt
```

### JSON Output
```bash
$ newestfiles -j .go .txt
["main.go","utils.go","README.txt"]
```

## Building

```bash
# Build the binary
go build

# Run tests
go test

# Install to GOPATH/bin
go install

```

## Testing

The project includes comprehensive tests covering:
- Plain text vs JSON output
- File filtering by extension
- Sorting by modification time (newest/oldest)
- Sorting by file size (largest/smallest)
- Conflicting sort flag validation
- Edge cases (no files found, no arguments)
- Case-insensitive matching
- Extension normalization

Run tests with:
```bash
go test
```

