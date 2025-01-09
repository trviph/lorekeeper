# Lorekeeper

[![Go Reference](https://pkg.go.dev/badge/github.com/trviph/lorekeeper.svg)](https://pkg.go.dev/github.com/trviph/lorekeeper) [![CI](https://github.com/trviph/lorekeeper/actions/workflows/ci.yaml/badge.svg)](https://github.com/trviph/lorekeeper/actions/workflows/ci.yaml) [![codecov](https://codecov.io/gh/trviph/lorekeeper/graph/badge.svg?token=7DDZ8QNJHW)](https://codecov.io/gh/trviph/lorekeeper)

Lorekeeper is a Go package that manages log rotation. It should work nicely with the Go standard log package.

**Note:** Lorekeeper is not a logging package, it only manages the log files, and should be use with other logging packages such as the standard [log](#with-standard-log-package), [slog](#with-standard-slog-package), or [logrus](#with-logrus).

## Quick Guide

### With Standard log Package

Package [log](https://pkg.go.dev/log) implements a simple logging package.

```go
package main

import (
    "log"
    "github.com/trviph/lorekeeper"
)

func main() {
    // Init lorekeeper, with the default configurations.
    defaultKeeper, err := lorekeeper.NewKeeper()
    if err != nil {
        panic(err)
    }
    defer defaultKeeper.Close()

    // Using lorekeeper with log
    logger := log.New(keeper, "[INFO] ", log.Lmsgprefix|log.LstdFlags)

    // Starting using the logger
    logger.Printf("that's it!")
}
```

### With Standard slog Package

Package [slog](https://pkg.go.dev/log/slog) provides structured logging, in which log records include a message, a severity level, and various other attributes expressed as key-value pairs.

```go
package main

import (
    "log/slog"
    "github.com/trviph/lorekeeper"
)

func main() {
    // Init lorekeeper, with the default configurations.
    defaultKeeper, err := lorekeeper.NewKeeper()
    if err != nil {
        panic(err)
    }
    defer defaultKeeper.Close()

    // Using lorekeeper with slog
    logger := slog.New(slog.NewJSONHandler(keeper, nil))

    // Starting using the logger
    logger.Info("this is info", "msg", "testing")
    logger.Error("this is error", "msg", "testing")
}
```

### With Logrus

[Logrus](https://github.com/sirupsen) is a structured, pluggable logging package for Go.  

```go
package main

import (
    log "github.com/sirupsen/logrus"
    "github.com/trviph/lorekeeper"
)

func main() {
    // Init lorekeeper, with the default configurations.
    defaultKeeper, err := lorekeeper.NewKeeper()
    if err != nil {
        panic(err)
    }
    defer defaultKeeper.Close()

    // Using lorekeeper with logrus
    logger := log.New()
    logger.SetLevel(log.InfoLevel)
    logger.SetOutput(defaultKeeper)

    // Starting using the logger
    logger.Info("this will go into the log file")
    logger.Debug("this will not")
}
```
