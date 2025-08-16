# Release Notes v1.0.6

## 🔧 Fixes

### Go Module Checksum Issue Resolution

- **Fixed:** Resolved Go module checksum mismatch error that prevented installation via `go install`
- **Issue:** The v1.0.5 tag was modified after publication, causing Go's checksum database to reject the module
- **Solution:** Created new immutable v1.0.6 tag to restore proper module installation

### What was the problem?

Users encountered this error when trying to install:

```
go: github.com/rocketpowerinc/go-pwr@v1.0.5: verifying module: checksum mismatch
        downloaded: h1:6usfd4+VfVqEcU6ZbGRIROb1KveHsvyzh0YB++O6eCw=
        sum.golang.org: h1:WTyArmww6Ylbx+z4y1hvfc6xH4GZQTK1psLGUQ6YcFM=

SECURITY ERROR
This download does NOT match the one reported by the checksum server.
```

### How it's fixed

- Created v1.0.6 as a clean, immutable tag
- Go's checksum database will now properly verify this version
- Installation via `go install` works correctly again

## 📥 Installation

```powershell
# Install latest version (recommended)
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@latest

# Or install specific version
go install -v github.com/rocketpowerinc/go-pwr/cmd/go-pwr@v1.0.6
```

## 📝 Notes

- No functional changes from v1.0.5
- This is purely a maintenance release to restore module installation
- Updated Windows-Bootstrap.md documentation to reflect the fix

## 🔗 Links

- [Repository](https://github.com/rocketpowerinc/go-pwr)
- [Installation Guide](https://github.com/rocketpowerinc/go-pwr/blob/main/Windows-Bootstrap.md)
