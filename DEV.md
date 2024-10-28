# Development

This section provides all the necessary information to start contributing to or building SVMKit. It is divided into several subsections for clarity.

# Requirements

Ensure the following tools are installed and available in your `$PATH`:

- [golang 1.22](https://golang.org/dl/)
- [`golangci-lint`](https://golangci-lint.run/install)

You can use [Homebrew](https://brew.sh) or your preferred package manager to get these tools.

#### Structure

The repository includes the following:

| Directory/File | Description                                                   |
| -------------- | ------------------------------------------------------------- |
| `pkg`          | SVMKit Go packages.                                           |
| `build`        | Scripts for building validators and distributing through APT. |
| `Makefile`     | Contains commands for testing and building Go packages.       |
