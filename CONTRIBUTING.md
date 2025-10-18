# Contributing to lazydndplayer

Thank you for your interest in contributing to lazydndplayer!

## Development Setup

### Prerequisites
- Go 1.x or later
- Make (for using Makefile commands)
- Git

### Getting Started

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd lazydndplayer
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Build the application**
   ```bash
   make build
   ```

4. **Run tests**
   ```bash
   make test
   ```

## Development Workflow

### Before Making Changes

1. Create a new branch for your feature/fix
2. Ensure all tests pass: `make test`
3. Ensure code is properly formatted: `make fmt`

### While Developing

1. **Write tests** for new functionality in `tests/` directory
2. **Follow Go conventions** - use `gofmt` and `go vet`
3. **Run checks frequently**: `make check`
4. **Test your changes**: `make test`

### Before Committing

Run the full check suite:
```bash
make check
```

This will:
- Format your code (`make fmt`)
- Run static analysis (`make vet`)
- Run all tests (`make test`)

### Commit Messages

Use clear, descriptive commit messages:
- Start with a verb (Add, Fix, Update, Remove, etc.)
- Keep the first line under 50 characters
- Provide details in the body if needed

Good examples:
```
Add ability choice selector for feats
Fix crash when removing feats with nil BenefitTracker
Update feats.json with explicit choice structure
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific package tests
make test-models

# Run with coverage
make test-coverage

# View coverage in browser
make test-coverage-html
```

### Writing Tests

1. Place tests in `tests/` directory, organized by package
2. Use table-driven tests where appropriate
3. Test both success and failure cases
4. Include edge cases

Example test structure:
```go
func TestFeatureName(t *testing.T) {
    // Setup
    char := models.NewCharacter()
    
    // Execute
    result := SomeFunction(char)
    
    // Verify
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

## Code Style

### Go Style Guide

Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines:
- Use `gofmt` for formatting
- Use meaningful variable names
- Write clear comments for exported functions
- Keep functions focused and small

### Project Conventions

- **Models** go in `internal/models/`
- **UI Components** go in `internal/ui/components/`
- **UI Panels** go in `internal/ui/panels/`
- **Tests** go in `tests/<package>/`
- **Data files** go in `data/`

## Project Structure

```
lazydndplayer/
├── data/              # JSON data files (feats, spells, species)
├── internal/
│   ├── models/        # Core data models and logic
│   ├── storage/       # Save/load functionality
│   └── ui/            # Terminal UI components
│       ├── components/  # Reusable UI components
│       └── panels/      # Main UI panels
├── tests/
│   └── models/        # Test files
├── main.go            # Application entry point
└── Makefile           # Build and test commands
```

## Adding New Features

### Feature Development Checklist

- [ ] Create feature branch
- [ ] Write tests first (TDD approach recommended)
- [ ] Implement feature
- [ ] Update documentation
- [ ] Run `make check` to ensure quality
- [ ] Test manually in the application
- [ ] Update CHANGELOG.md
- [ ] Create pull request with clear description

### Adding New Models

1. Create model struct in `internal/models/`
2. Add JSON tags for serialization
3. Write constructor function (e.g., `NewFeature()`)
4. Add tests in `tests/models/`
5. Update character model if needed

### Adding UI Components

1. Create component in `internal/ui/components/`
2. Implement required methods: `Update()`, `View()`, `IsVisible()`, etc.
3. Add keyboard bindings
4. Update help text
5. Test in application

## Debugging

### Debug Logging

The application uses standard Go logging. To add debug logs:
```go
import "log"

log.Printf("Debug: value is %v", someValue)
```

### Common Issues

**Tests failing after model changes**
- Update test mocks/fixtures
- Ensure BenefitTracker is initialized
- Check JSON serialization

**UI not responding**
- Check keyboard handler priority
- Verify component visibility state
- Look for blocking operations

**Build failing**
- Run `make clean`
- Run `make deps`
- Check Go version: `go version`

## Making a Pull Request

1. Ensure all tests pass: `make test`
2. Ensure code is formatted: `make fmt`
3. Ensure no vet warnings: `make vet`
4. Update documentation if needed
5. Create PR with:
   - Clear title
   - Description of changes
   - Related issue numbers
   - Test results
   - Screenshots (if UI changes)

## Questions?

If you have questions or need help:
- Check existing issues and PRs
- Read the documentation in `*.md` files
- Look at similar features for examples

## Code of Conduct

- Be respectful and constructive
- Welcome newcomers
- Focus on the code, not the person
- Help others learn and grow

Thank you for contributing! 🎉

