# Contributing to assets-gen

Thank you for your interest in contributing to **assets-gen**! We welcome contributions of all kinds: feature enhancements, bug fixes, documentation improvements, and more.

---

## 📝 How to Contribute

### 1. Reporting Bugs

1. Check existing issues: search for similar reports in [Issues](https://github.com/Nidal-Bakir/assets-gen/issues) to avoid duplicates.
2. Open a new issue and include:

   - **Title**: concise summary of the problem.
   - **Description**: detailed steps to reproduce.
   - **Environment**: OS, Go version (`go version`), CLI version (`assets-gen --version`).
   - **Output**: error messages, logs, or screenshots.

### 2. Suggesting Enhancements

1. Search existing feature requests in [Issues](https://github.com/Nidal-Bakir/assets-gen/issues).

2. If none exist, open a new issue:

   - **Title**: short summary of the request.
   - **Motivation**: why this would be useful.
   - **Example**: CLI command usage or expected behavior.

### 3. Submitting a Pull Request (PR)

1. **Fork** the repository.

2. **Create a branch**: `git checkout -b feature/your-feature-name` or `bugfix/your-bug-description`.

3. **Make your changes**:

   - Follow Go best practices and project style.
   - Write clear, descriptive commit messages (see below).
   - Use the Makefile for formatting and linting:
     ```bash
     make fmt    # formats code with go fmt
     make chk    # runs go vet & staticcheck
     ```

4. **Commit and push** your branch to your fork.

5. **Open a pull request** against the `main` branch:
   - Provide a clear description of changes.
   - Link any related issues (use `#issue-number`).
   - Describe how to verify the changes. your branch to your fork.

---

## 📋 Code Style & Guidelines

- **Go version**: Support Go 1.18 and above.
- **Formatting**: Use `make fmt`
- **Linting**: Ensure no `go vet` or `staticcheck` errors (`make chk`).
- **Documentation**:

  - Update `README.md` for new flags or commands.
  - Document exported functions in GoDoc style.
  - Update the cli verion in `cmd/cli/main.go` to reflect the new version

---

## ✍️ Commit Message Convention

Use [Conventional Commits](https://www.conventionalcommits.org/) for consistency:

```
<type>(<scope>): <subject>

<body>

<footer>
```

- **type**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`.
- **scope**: optional, e.g., `cmd`, `assetsgen`, `cli`.
- **subject**: brief summary (max 72 chars).
- **body**: detailed explanation, if necessary.
- **footer**: reference issues (`Closes #123`).

**Example:**

```
feat(core): add support for web icons

This will allow the cli to generate web icons

Closes #45
```

---
