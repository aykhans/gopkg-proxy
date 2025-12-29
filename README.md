# gopkg-proxy

A lightweight vanity import path server for Go packages. Use your own domain for clean, branded package imports.

## Configuration

### Add Packages

Edit `main.go`:

```go
var packages = []Package{
    {
        Path: "/my-package",
        Repo: "https://github.com/username/my-package",
        VCS:  "git",
    },
}
```

### Environment Variables

- `PORT` - Server port (default: `8421`)
- `HOST` - Override the domain used in package paths (if set, takes priority over headers)
- `HOST_HEADER` - Header to read domain from (default: `X-Forwarded-Host`)

## Usage

Instead of:

```go
import "github.com/aykhans/go-utils"
```

Use:

```go
import "gopkg.yourdomain.com/go-utils"
```
