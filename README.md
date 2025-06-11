# newestfiles

A command-line tool to find and list files by extension, sorted by modification time (newest first).

## Features

- Find files by extension(s) in current directory and subdirectories
- Sort files by modification time (newest first)
- Support for multiple file extensions
- Case-insensitive extension matching
- Two output formats: plain text (default) and JSON
- Automatic extension normalization (adds dot if missing)

## Usage

```bash
# Plain text output (default)
newestfiles .go .txt .md

# JSON output
newestfiles -j .go .txt .md

# Extensions without dots are automatically normalized
newestfiles go txt md
```

### Command Line Options

- `-j`: Output in JSON format (default: plain text)

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
make build

# Run tests
make test

# Run all checks (format, vet, test)
make check

# Install to GOPATH/bin
make install

# Clean build artifacts
make clean
```

## Testing

The project includes comprehensive tests covering:
- Plain text vs JSON output
- File filtering by extension
- Sorting by modification time
- Edge cases (no files found, no arguments)
- Case-insensitive matching
- Extension normalization

Run tests with:
```bash
make test
```

## Development

- Format code: `make fmt`
- Run linter: `make vet`
- Generate coverage report: `make test-coverage`
- Run example: `make run-example`
