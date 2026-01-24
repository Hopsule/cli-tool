# Homebrew Core PR Template

When submitting to Homebrew Core, use this as a template for your PR.

---

## PR Title

```
hopsule 0.5.0 (new formula)
```

## PR Description

```markdown
### Description

Add **hopsule** - Decision & Memory Layer CLI for AI teams & coding tools.

Hopsule provides a decision-first workflow management system with:
- Structured decision capture and lifecycle
- Memory preservation for teams and AI agents
- Context-aware enforcement
- Portable decision packages

### Formula Details

- **License**: MIT
- **Version**: 0.5.0
- **Homepage**: https://github.com/Hopsule/cli-tool
- **Build**: From source using Go
- **Platform**: macOS and Linux
- **Dependencies**: Go (build-time only)

### Quality Checklist

- [x] Formula builds from source
- [x] No pre-compiled binaries
- [x] Includes test block
- [x] Follows Homebrew naming conventions
- [x] Uses `std_go_args` for standard Go builds
- [x] Proper license specified
- [x] Head formula for development
- [x] Works on Apple Silicon and Intel
- [x] Passes `brew audit --strict --online`
- [x] Passes `brew test`

### CI/CD

- **GitHub Actions**: https://github.com/Hopsule/cli-tool/actions
- **Tests**: Comprehensive unit tests
- **Linter**: golangci-lint
- **Coverage**: Growing (currently 15%)

### Documentation

- README: https://github.com/Hopsule/cli-tool#readme
- Contributing: https://github.com/Hopsule/cli-tool/blob/main/CONTRIBUTING.md
- Installation: https://github.com/Hopsule/cli-tool/blob/main/INSTALL.md

### Testing Done

```bash
# Install from formula
brew install --build-from-source ./Formula/h/hopsule.rb

# Verify install
which hopsule
hopsule --version

# Run tests
brew test hopsule

# Audit
brew audit --strict --online hopsule
```

### Community

- **GitHub Stars**: [Current count]
- **Production Users**: [Growing]
- **Active Development**: Regular releases
- **Responsive Maintainers**: Quick issue resolution

### Additional Notes

This formula enables installation via:
```bash
brew install hopsule
```

Currently available via tap: `brew install hopsule/tap/hopsule`

---

cc @Homebrew/maintainers
```

---

## Pre-Submission Checklist

Before opening the PR, ensure:

1. **Formula Location**: `Formula/h/hopsule.rb`

2. **Formula Audit**:
   ```bash
   cd $(brew --repo homebrew/core)
   brew audit --strict --online hopsule
   ```

3. **Formula Test**:
   ```bash
   brew install --build-from-source hopsule
   brew test hopsule
   ```

4. **Build on Multiple Platforms**:
   ```bash
   # macOS ARM64
   brew install --build-from-source hopsule
   
   # macOS Intel (if available)
   brew install --build-from-source hopsule
   
   # Linux (via GitHub Actions or container)
   ```

5. **Check Dependencies**:
   ```bash
   brew deps hopsule
   # Should only show build dependencies (go)
   ```

6. **Verify Installation**:
   ```bash
   hopsule --version
   hopsule --help
   hopsule config
   ```

7. **Clean Up**:
   ```bash
   brew uninstall hopsule
   brew cleanup
   ```

---

## Common Review Comments

### If reviewers ask about popularity:

> Hopsule is actively used by [X teams/users]. We have [Y stars] on GitHub and growing community engagement. The project solves a real problem in AI-assisted development workflows.

### If reviewers ask about stability:

> We follow semantic versioning with regular releases. The project has comprehensive tests, CI/CD, and active maintenance. Current version (0.5.0) is stable and production-ready.

### If reviewers suggest improvements:

> Thank you for the feedback! I'll update the formula accordingly.

---

## After Approval

1. **Merge**: Wait for maintainer approval and merge
2. **Test**: Install from Homebrew Core
   ```bash
   brew update
   brew install hopsule
   ```
3. **Update Docs**: Update README to show official Homebrew install
4. **Announce**: Share the news with community
5. **Archive Tap**: Keep tap for pre-release versions

---

## Resources

- [Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)
- [Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [How to Open a PR](https://docs.brew.sh/How-To-Open-a-Homebrew-Pull-Request)
- [Maintainer Guidelines](https://docs.brew.sh/Maintainer-Guidelines)
