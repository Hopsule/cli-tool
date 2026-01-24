# Contributing to Hopsule CLI

Thank you for your interest in contributing to Hopsule CLI! 🎉

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Setup

```bash
# Clone the repository
git clone https://github.com/Hopsule/cli-tool.git
cd cli-tool

# Install dependencies
go mod download

# Build
go build -o hopsule ./cmd/decision

# Run tests
go test ./...

# Run linter
golangci-lint run
```

## Code Standards

### Style Guide

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write tests for new features
- Keep functions small and focused

### Testing

- Write unit tests for all new code
- Maintain test coverage above 80%
- Run `go test -v -race ./...` before committing

### Commit Messages

Follow conventional commits:

```
feat: add new command
fix: resolve issue with config loading
docs: update README
test: add tests for interactive mode
chore: update dependencies
```

## Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Write** your code and tests
4. **Run** tests and linter
5. **Commit** your changes
6. **Push** to your fork
7. **Open** a Pull Request

### PR Checklist

- [ ] Tests pass (`go test ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Documentation updated
- [ ] Commit messages follow convention
- [ ] Code is formatted (`gofmt`)

## Architecture

```
cli-tool/
├── cmd/decision/        # Main entry point
├── internal/
│   ├── api/            # API client
│   ├── commands/       # CLI commands
│   ├── config/         # Configuration management
│   └── ui/             # Terminal UI components
├── .github/workflows/  # CI/CD
└── README.md
```

## Key Components

### Commands (`internal/commands/`)
Individual CLI command implementations.

### UI (`internal/ui/`)
Terminal user interface using bubbletea.

### Config (`internal/config/`)
Configuration loading and validation using viper.

## Testing Guidelines

### Unit Tests
```go
func TestNewFeature(t *testing.T) {
    // Arrange
    input := "test"
    
    // Act
    result := NewFeature(input)
    
    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Integration Tests
Place integration tests in `*_integration_test.go` files.

## Release Process

Releases are automated via GitHub Actions:

1. Tag a version: `git tag v1.0.0`
2. Push tag: `git push origin v1.0.0`
3. GitHub Actions builds and releases automatically

## Questions?

- Open an issue for bug reports
- Start a discussion for feature requests
- Reach out via GitHub Discussions

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
