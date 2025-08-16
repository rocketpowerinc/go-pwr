# Release Notes v1.0.7

## 🔧 Critical Fix

### Go Module Proxy Caching Issue Resolution

- **Fixed:** Resolved Go module proxy caching issue that caused v1.0.6 to install with incorrect version number
- **Issue:** The v1.0.6 tag was moved after initial publication, causing Go's proxy to cache the wrong commit
- **Solution:** Created new immutable v1.0.7 tag with proper version string in code

### What was the problem?

Even after installing `@v1.0.6`, users saw:

```
go-pwr v1.0.5
Git commit: unknown
Build date: unknown
```

This occurred because Go's module proxy cached the original v1.0.6 tag before it was updated with the version bump.

### How it's fixed

- Created v1.0.7 as a clean, immutable tag with correct version string
- Updated `const version = "1.0.7"` in main.go
- Go's proxy will cache this version correctly on first access

## 📥 Installation

```powershell
# Install latest version (recommended)
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest

# Or install specific version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.7
```

Now correctly shows:
```
go-pwr v1.0.7
Git commit: unknown
Build date: unknown
Built with Go go1.25.0 for windows/amd64
Repository: https://github.com/rocketpowerinc/go-pwr
```

## 📝 Notes

- No functional changes from v1.0.6
- This is purely a maintenance release to fix version reporting
- Updated Windows-Bootstrap.md documentation to reflect v1.0.7

## 🎓 Lesson Learned

**Never move/retag Git tags after publication** - Go's module proxy caches them immutably. Always create new versions instead.

## 🔗 Links

- [Repository](https://github.com/rocketpowerinc/go-pwr)
- [Installation Guide](https://github.com/rocketpowerinc/go-pwr/blob/main/Windows-Bootstrap.md)
