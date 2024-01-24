# LeoLog
This is a logging handler for the Go [structured logger](https://pkg.go.dev/log/slog)

It is originally based of [this guide](https://dusted.codes/creating-a-pretty-console-logger-using-gos-slog-package)

It is used to print pretty logs, albeit not very fast.

Do not use this for performance critical applications.

## Usage

The following is a simple example of how to use this package.

```go
package main

import (
    "log/slog"

    "github.com/theleeeo/leolog"
)

func main() {
    logger := slog.New(prettylog.NewHandler(nil))
    slog.SetDefault(logger)

    slog.Into("ahhh there are cowboys in my bedroom!", "yee", "haw")
}
```