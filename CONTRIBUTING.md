# Contributing to xAI SDK for Go

Thank you for your interest in contributing! The project follows a phased roadmap described in [`docs/development-plan.md`](docs/development-plan.md). Please review the guidelines below before opening issues or pull requests.

## Getting Started

1. **Fork and clone** the repository.
2. Install Go 1.22 or newer, and install [Buf](https://buf.build) for protobuf tooling.
3. Run `go test ./...` to ensure the workspace is set up correctly.
4. Review open issues or create a new issue describing the change you would like to make.

## Development Workflow

- **Branch naming**: use prefixes such as `feature/`, `fix/`, or `chore/` (e.g., `feature/chat-stream`).
- **Commit style**: follow conventional-style prefixes (e.g., `feat(chat): ...`, `fix(files): ...`). Each phase in `docs/development-plan.md` should result in at least one commit capturing the described work.
- **Testing**: include unit tests for new functionality. Integration tests should be gated behind environment variables as described in the plan.
- **Code style**: run `gofmt` and lint the code (golangci-lint will be introduced in a later phase).

## Pull Requests

- Reference the related issue and the development phase the change corresponds to.
- Include a summary of the change, testing performed, and any follow-up work required.
- Ensure documentation and examples are updated when functionality changes.

## Reporting Issues

- Provide as much detail as possible: Go version, OS, steps to reproduce, expected vs. actual behavior.
- For security issues, please reach out through responsible disclosure channels (will be documented once available).

## Code of Conduct

This project adheres to a Code of Conduct to foster a welcoming community. Please read [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) before participating.

Thank you for helping build the xAI Go SDK!
