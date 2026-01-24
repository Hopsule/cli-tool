# Homebrew Core Submission Checklist

This document tracks readiness for Homebrew Core submission.

## Requirements Status

### ✅ Completed

- [x] **License**: MIT License added
- [x] **Open Source**: Repository is public on GitHub
- [x] **Stable Release**: v0.4.4 released with proper versioning
- [x] **CI/CD**: GitHub Actions for testing, linting, and building
- [x] **Tests**: Unit tests with good coverage
- [x] **Documentation**: README, INSTALL, CONTRIBUTING, and INTERACTIVE_GUIDE
- [x] **Formula**: Homebrew-compatible formula created
- [x] **Build from Source**: Formula builds from source (not pre-built binary)
- [x] **Cross-platform**: Builds on macOS (Darwin) and Linux

### 🔄 In Progress

- [ ] **Popularity**: Need more GitHub stars and usage
- [ ] **Stability**: Need more production usage and feedback
- [ ] **Test Coverage**: Aim for >80% coverage
- [ ] **More Tests**: Add integration tests
- [ ] **Shell Completions**: Add bash/zsh/fish completions

### 📋 Before Submission

1. **Verify Formula**:
   ```bash
   brew install --build-from-source ./homebrew-core-formula.rb
   brew test hopsule
   brew audit --strict --online hopsule
   ```

2. **Calculate SHA256** for release:
   ```bash
   curl -L https://github.com/Hopsule/cli-tool/archive/v0.4.4.tar.gz | shasum -a 256
   ```

3. **Update Formula** with correct SHA256

4. **Fork homebrew-core**:
   ```bash
   cd $(brew --repo homebrew/core)
   git checkout -b hopsule
   ```

5. **Copy Formula**:
   ```bash
   cp homebrew-core-formula.rb $(brew --repo homebrew/core)/Formula/h/hopsule.rb
   ```

6. **Test Thoroughly**:
   ```bash
   brew install --build-from-source hopsule
   brew test hopsule
   brew audit --strict --online hopsule
   ```

7. **Submit PR** to homebrew-core with description:
   ```
   hopsule 0.4.4 (new formula)
   
   Decision & Memory Layer CLI for AI teams
   
   - Build from source
   - MIT licensed
   - Stable release v0.4.4
   - CI/CD with GitHub Actions
   - Comprehensive tests and documentation
   ```

## Homebrew Core Standards

### Formula Requirements

- ✅ Build from source (no pre-compiled binaries)
- ✅ Compatible with current macOS versions
- ✅ Uses standard Go build process
- ✅ Includes test block
- ✅ Proper license specified
- ✅ Head formula for development builds

### Code Quality

- ✅ Clean Go code following best practices
- ✅ Comprehensive tests
- ✅ Linting with golangci-lint
- ✅ CI/CD pipeline
- ✅ Good documentation

### Project Maturity

- ⏳ Need more GitHub stars (current: ~0, target: 50+)
- ⏳ Need more production users
- ⏳ Need stable track record (few months)

## Timeline

1. **Now**: All technical requirements met ✅
2. **Next 1-2 months**: Build community and gather feedback
3. **After gaining traction**: Submit to Homebrew Core

## Alternative Distribution

While building popularity for Homebrew Core:

1. **Current Tap**: `brew install hopsule/tap/hopsule` ✅
2. **Direct Download**: GitHub releases ✅
3. **Go Install**: `go install github.com/Hopsule/cli-tool/cmd/decision@latest`
4. **Package Managers**: Consider apt, yum, pacman, etc.

## Resources

- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [How to Open a Homebrew Pull Request](https://docs.brew.sh/How-To-Open-a-Homebrew-Pull-Request)
- [Homebrew Core](https://github.com/Homebrew/homebrew-core)

## Notes

- Homebrew Core typically accepts software that is:
  - Well-known and widely used
  - Stable and maintained
  - Has a significant user base
  
- Our current approach:
  1. Perfect the formula and codebase ✅
  2. Build community and user base 🔄
  3. Submit when we meet popularity thresholds ⏳
