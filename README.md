# Lorekeeper

[![CI](https://github.com/trviph/lorekeeper/actions/workflows/ci.yaml/badge.svg)](https://github.com/trviph/lorekeeper/actions/workflows/ci.yaml)

Lorekeeper is a Go package that handles log files rotation. Lorekeeper should work well with the standard log library.

**Note:** Lorekeeper is not a logging package, it only manages the log files.

## Quick Guide

### With Standard Log Package

For more in-depth usage, see this [file](tests/log_test.go).

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
