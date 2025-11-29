# Contributing to Volt 

Thanks for wanting to contribute to Volt! Whether you're fixing a typo, adding a feature, or just poking around the codebase, we appreciate you being here.

## Quick Start

1. **Clone the Repo**
   ```bash
   git clone https://github.com/owenHochwald/Volt.git
   cd volt
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Run Volt**
   ```bash
   go run main.go
   ```

4. **Run Tests**
   ```bash
   go test ./...
   ```


## What Can I Work On?

- **Good First Issues**: Check out issues labeled `good first issue` - these are beginner-friendly
- **Bug Fixes**: Found something broken? Fix it!
- **Features**: Have an idea? Open an issue first so we can discuss it
- **Documentation**: Improve README, add examples, write better comments
- **Performance**: Make Volt faster and more efficient

## Development Workflow

### 1. Create a Branch

Use descriptive branch names:
```bash
git checkout -b feat/request-history
git checkout -b fix/response-rendering
git checkout -b docs/update-readme
```

Prefixes we use:
- `feat/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation
- `refactor/` - Code refactoring
- `test/` - Adding tests
- `perf/` - Performance improvements

### 2. Write Code

**Code Style**
- Use meaningful variable names - no single letters unless it's an obvious loop counter
- Add comments for complex logic, especially in TUI event handling
- Keep functions focused - if it's doing too much, split it up

**Bubble Tea Best Practices**
- Keep your `Update()` functions pure - no side effects
- Use commands (`tea.Cmd`) for async operations
- Don't block the UI!!! Use goroutines with messages
- Test your components in isolation when possible

### 3. Write Tests

We're not aiming for 100% coverage, but core functionality should be tested!

**Testing TUI Components**
- Test model state changes directly
- Test complex utils or requester / request changes

**Run Tests**
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/http

# Verbose output
go test -v ./...
```

### 4. Commit Your Changes

Write clear commit messages.

### 5. Push & Create PR

```bash
git push origin your-branch-name
```

Then open a PR on GitHub!


## Performance Considerations

- **HTTP Requests**: Use connection pooling, don't create new clients for each request
- **Rendering**: Only update the view when state actually changes
- **Memory**: Be mindful of storing large responses in memory
- **Goroutines**: Always clean up goroutines when components unmount

## Documentation

- **Code Comments**: Explain *why*, not *what*
- **README**: Keep it up to date with new features
- **Examples**: Add usage examples for new features

## Questions?

- Open an issue for discussion
- Check existing issues and PRs - your question might be answered
- Be respectful and patient - we're all learning!

---

**Thanks for contributing!** Every PR, issue, and comment helps make this project better. 