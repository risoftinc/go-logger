# Contributing to Logger Package

Thank you for your interest in contributing to the Risoftinc. Logger package! This document provides guidelines for contributing to this project.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/). By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Issues

1. Check if the issue already exists in the [Issues](https://github.com/risoftinc/go-logger/issues) section
2. If not, create a new issue with:
   - Clear description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Go version and OS information

### Suggesting Enhancements

1. Check if the enhancement is already requested
2. Create an issue with the "enhancement" label
3. Provide detailed description of the proposed feature
4. Explain the use case and benefits

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes following the coding standards
4. Add tests for new functionality
5. Update documentation if needed
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

## Development Guidelines

### Code Style

- Follow Go standard formatting: `go fmt`
- Use meaningful variable and function names
- Add comments for public functions and types
- Keep functions small and focused

### Testing

- Write tests for new functionality
- Maintain test coverage above 90%
- Use table-driven tests where appropriate
- Test both success and error cases

### Documentation

- Update README.md for new features
- Add examples in the `example/` directory
- Update CHANGELOG.md for significant changes
- Document public APIs with Go doc comments

### Commit Messages

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Maintenance tasks

Examples:
```
feat(logger): add custom request ID key configuration
fix(context): handle nil context gracefully
docs(readme): update installation instructions
```

## Project Structure

```
logger/
├── logger.go              # Main package implementation
├── logger_test.go         # Unit tests
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── LICENSE                # MIT License
├── README.md              # Project documentation
├── CHANGELOG.md           # Change log
├── CONTRIBUTING.md        # This file
├── .gitignore            # Git ignore rules
├── example/              # Example programs
│   ├── main.go
│   ├── http_example.go
│   ├── context_behavior.go
│   └── custom_request_id.go
└── tools/                # Development tools
    └── test_runner.go
```

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/risoftinc/go-logger.git
   cd logger
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run tests:
   ```bash
   go test -v ./...
   ```

4. Run benchmarks:
   ```bash
   go test -bench=. -benchmem ./...
   ```

5. Check code coverage:
   ```bash
   go test -cover ./...
   ```

## Release Process

1. Update version in CHANGELOG.md
2. Update version in go.mod if needed
3. Create a release tag: `git tag v1.0.0`
4. Push the tag: `git push origin v1.0.0`
5. Create a GitHub release

## Questions?

If you have questions about contributing, please:
- Open an issue with the "question" label
- Contact the maintainers
- Check existing discussions

Thank you for contributing to the Logger package! 🚀
