// Lorekeeper is a Go package that manages log rotation.
//
// Note that Lorekeeper is not a full-blown logging package. It only receives logging messages and manages them in files and should be used as such.
// It should play nicely together with the standard [log], [log/slog] packages, or any packages that follow the standard structure, like [Logrus].
//
// The core of Lorekeeper is the [Keeper] struct, which implemented the [io.WriteCloser] interface.
// A [Keeper] with the same name should be safe to use in multiple goroutines in the same process,
// but not safe when using on multiple processes.
//
// # The Keeper Struct
//
// Keeper is a log manager that writes logs to files and rotates them.
// To create a Keeper, use the [NewKeeper] function. For example, the following code will create a Keeper with the default configuration:
//
//	import (
//		"log"
//
//		"github.com/trviph/lorekeeper"
//	)
//
//	func main() {
//		keeper, err := lorekeeper.NewKeeper()
//		if err != nil {
//			// Handle error
//		}
//		// Instrument standard log with Lorekeeper
//		logger := log.New(keeper, "[INFO] ", log.Lmsgprefix|log.LstdFlags)
//
//		// Every time we write to the logger, Lorekeeper will write the message to a file and handle it as per configuration
//		logger.Printf("this is a log message")
//	}
//
// # Configure the Keeper
//
// This package provides some configurations for the [Keeper].
// These configurations come in the form of WithXxx functions that follow the Go Options pattern.
// You should take a look at [Opt] and the WithXxx functions in the Go package reference for documentation on these configurations.
// An example of how to use these functions:
//
//	import (
//		"log"
//
//		"github.com/trviph/lorekeeper"
//	)
//
//	func main() {
//		keeper, err := lorekeeper.NewKeeper(
//			// Setting the name for the Keeper, this may affect how files will be named
//			lorekeeper.WithName("Example"),
//			// Setting the maximum size of 100 MegaBytes, if a log file exceeds this it will be rotated
//			lorekeeper.WithMaxSize(100*lorekeeper.MB),
//		)
//		if err != nil {
//			// Handle error
//		}
//		// Instrument standard log with Lorekeeper
//		logger := log.New(keeper, "[INFO] ", log.Lmsgprefix|log.LstdFlags)
//
//		// Every time we write to the logger, Lorekeeper will write the message to a file and handle it as per configuration
//		logger.Printf("this is a log message")
//	}
//
// # How Does This Work
//
// When creating a Keeper, it first looks into the folder containing logs, defined using the [WithFolder] option. It will throw an error if the folder does not yet exist, so make sure the folder exists, and has appropriate permissions.
//
// The Keeper will then scan the folder for any related logs in the folder.
// There are two kinds of log stored here, the first kind is the current log, which the Keeper is writing to, the name of the current log depends on [WithName] and [WithExtension] options.
// If the Keeper finds an existing current log, it will reuse that log, if not it will create a new one.
// The second kind is the archived log, which the keeper is keeping track of, the name of archives depends on [WithName], [WithExtension], [WithTimeLayout], and [WithArchiveNameLayout].
// Since the Keeper depends on the file name to determine which file to manage, be aware that changing any of the mentioned options will cause the logs from the previous execution to become orphaned and not managed by the Keeper.
//
// Every time the [Keeper.Write] is invoked, the Keeper will first check if the current log should be rotated before writing the message to the log.
//
// A rotation will happen if the current log size exceeds the max size, configured by using [WithMaxSize] option.
// During a rotation, the Keeper archives the current log by closing and renaming it based on the name template configured by using [WithArchiveNameLayout] and then opens a new log to replace the archived one.
// Afterwards if the number of archives exceeds the maximum number of allowed files, configured by [WithMaxFiles], the Keeper will keep deleting the oldest archives based on its last modified time until the number of archives is smaller than the configured value.
//
// You can also rotate the current log manually by using [Keeper.Rotate] or [Keeper.Close], the difference between these two functions is that after [Keeper.Rotate], you can continue to use the Keeper as it rotates the current log but keeps it open for further writing.
// [Keeper.Close] will rotate and close the current log preventing any subsequent call from writing more messages into it.
//
// The Keeper can rotate log based on cron schedule, see [WithCron] option for more info.
//
// [Logrus]: https://github.com/sirupsen/logrus
package lorekeeper
