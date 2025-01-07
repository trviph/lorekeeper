package lorekeeper

import "os"

// A [Keeper] is a log file manager that handles writing to log files and rotates them.
// Use [DefaultKeeper] or [NewKeeper] to create a new Keeper.
type Keeper struct {
	// See [WithFolder] for documentation.
	folder string
	// See [WithName] for documentation.
	name string
	// See [WithExtension] for documentation.
	extension string
	// See [WithTimeFormat] for documentation
	timeFormat string
	// See [WithMaxsize] for documentation
	maxsize int
}

// Create a new [Keeper] with the default configurations.
// See [NewKeeper] if you want to configure the Keeper yourself.
func DefaultKeeper() *Keeper {
	return &Keeper{
		folder:     os.TempDir(),
		name:       defaultKeeperName(),
		extension:  ".log",
		timeFormat: "20060102T150405.999999999MST",
		maxsize:    MB,
	}
}

// Create a new [Keeper] with the provided options.
// This will create a [DefaultKeeper] if no option is provided.
// If at least one option is provided, this may also return an error if the option is invalid.
// See [Opt] for all available options.
//
// Example usage:
//
//		import "github.com/trviph/loremaster"
//
//		func main() {
//			keeper, err := lorekeeper.NewKeeper(
//				lorekeeper.WithName("Lorekeeper Example"),
//				lorekeeper.WithMaxByte(12 * lorekeeper.Kb),
//	 	)
//		}
func NewKeeper(opts ...Opt) (*Keeper, error) {
	keeper := DefaultKeeper()
	var err error
	for _, opt := range opts {
		keeper, err = opt(keeper)
		if err != nil {
			break
		}
	}
	return keeper, err
}
