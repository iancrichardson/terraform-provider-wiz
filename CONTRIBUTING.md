# Contributing to Terraform Provider for Wiz

Thank you for your interest in contributing to the Terraform Provider for Wiz! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

## Code of Conduct

Please be respectful and considerate of others when contributing to this project. We aim to foster an inclusive and welcoming community.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```sh
   git clone https://github.com/YOUR_USERNAME/terraform-provider-wiz.git
   cd terraform-provider-wiz
   ```
3. Add the original repository as an upstream remote:
   ```sh
   git remote add upstream https://github.com/iancrichardson/terraform-provider-wiz.git
   ```
4. Create a new branch for your changes:
   ```sh
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

### Setting Up Your Development Environment

1. Install Go (version 1.22 or later)
2. Install Terraform (version 0.13 or later)
3. Install any IDE or editor of your choice (VSCode, GoLand, etc.)

### Building and Testing

1. Make your changes to the provider code
2. Build the provider:
   ```sh
   go build -o terraform-provider-wiz
   ```
3. Install the provider locally for testing:
   ```sh
   mkdir -p ~/.terraform.d/plugins/github.com/iancrichardson/wiz/0.4.0/$(go env GOOS)_$(go env GOARCH)/
   cp terraform-provider-wiz ~/.terraform.d/plugins/github.com/iancrichardson/wiz/0.4.0/$(go env GOOS)_$(go env GOARCH)/
   ```
4. Run tests:
   ```sh
   go test ./...
   ```

### Code Style and Guidelines

- Follow standard Go coding conventions
- Use meaningful variable and function names
- Write clear comments and documentation
- Include tests for new features or bug fixes

## Pull Request Process

1. Update your fork with the latest changes from the upstream repository:
   ```sh
   git fetch upstream
   git rebase upstream/main
   ```
2. Push your changes to your fork:
   ```sh
   git push origin feature/your-feature-name
   ```
3. Create a Pull Request from your fork to the main repository
4. Provide a clear description of the changes and any relevant issue numbers
5. Wait for review and address any feedback

## Release Process

Releases are managed by the maintainers of the project. If you're a maintainer, follow these steps:

1. Update the version in `main.go`
2. Update the version in example files and README
3. Update the CHANGELOG.md file with the changes in the new version
4. Commit these changes with a message like "Bump version to v0.x.0"
5. Create and push a new tag:
   ```sh
   git tag v0.x.0
   git push origin v0.x.0
   ```
6. The GitHub Actions workflow will automatically build the provider for all supported platforms and create a new release

Thank you for contributing to the Terraform Provider for Wiz!
