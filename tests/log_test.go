package lorekeeper

import (
	"log"
	"sync"
	"testing"

	"github.com/trviph/lorekeeper"
)

// This test demonstrates on how to use [lorekeeper.Keeper] with the std [log].
func TestLog(t *testing.T) {
	// Create a Keeper
	keeper, err := lorekeeper.NewKeeper(
		// Set the Keeper name, this will be used when generate log files.
		lorekeeper.WithName("Lorekeeper Test Log"),
		// Set the extension of archived logs.
		lorekeeper.WithExtension(".old"),
		// Set the time layout of archived logs.
		lorekeeper.WithTimeLayout("20060102150405.000"),
		// Specify the folder where the log files will be stored.
		lorekeeper.WithFolder("."),
		// Each log file hold a maximum of 50 Kibibyte before being rotated.
		lorekeeper.WithMaxsize(50*lorekeeper.Kb),
		// Set the name layout of archived logs.
		lorekeeper.WithArchiveNameLayout("{{ .name }}_{{ .time }}"),
	)
	if err != nil {
		t.Errorf("failed to create a new keeper, caused by %s", err)
	}

	// Create loggers
	debugLogger := log.New(keeper, "[DEBUG] ", log.Lmsgprefix|log.LstdFlags|log.Llongfile)
	infoLogger := log.New(keeper, "[INFO] ", log.Lmsgprefix|log.LstdFlags)
	warningLogger := log.New(keeper, "[WARN] ", log.Lmsgprefix|log.LstdFlags|log.Lshortfile)

	// Use loggers
	debugLogger.Printf("this is a debug information")
	infoLogger.Printf("this is an additional information")
	warningLogger.Printf("i am warning you")

	// You should see multiple log files being created
	var wg sync.WaitGroup

	n := 1000
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			debugLogger.Printf("[%d] flooding the log with debug information...", id)
			infoLogger.Printf("[%d] flooding the log with additional information...", id)
			warningLogger.Printf("[%d] flooding the log with warning...", id)
		}(i)
	}

	wg.Wait()
}
