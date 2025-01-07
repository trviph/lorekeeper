# Lorekeeper

[![CI](https://github.com/trviph/lorekeeper/actions/workflows/ci.yaml/badge.svg)](https://github.com/trviph/lorekeeper/actions/workflows/ci.yaml)

Lorekeeper is a Go package that handles log files rotation. Lorekeeper should work well with the standard log library.

**Note:** Lorekeeper is not a logging package, it only manages the log files.

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

    // Using lorekeeper with logrus
    logger := log.New()
    logger.SetLevel(log.InfoLevel)
    logger.SetOutput(defaultKeeper)

    // Starting using the logger
    logger.Info("this will go into the log file")
    logger.Debug("this will not")
}
```
