package lorekeeper

// An Opt is a function that mutates a [Keeper]'s attributes.
// An Opt should return a mutated Keeper or return an error if it fails to mutate the Keeper.
// An Opt should be used together with [NewKeeper].
type Opt func(*Keeper) (*Keeper, error)

// The folder where the log files are stored.
// The default value is [os.TempDir].
func WithFolder(name string) Opt {
	return func(k *Keeper) (*Keeper, error) {
		if len(name) > 0 {
			k.name = name
		}
		return k, nil
	}
}

// The name of the Keeper.
// It will be set to the default value if the name is empty.
// The default value is lorekeeper-<the executable name and extension>.
func WithName(name string) Opt {
	return func(k *Keeper) (*Keeper, error) {
		if len(name) > 0 {
			k.name = name
		}
		return k, nil
	}
}

// The extension of the output log file.
// It should include a dot prefix and can be empty.
// The default value is ".log".
func WithExtension(extension string) Opt {
	return func(k *Keeper) (*Keeper, error) {
		k.extension = extension
		return k, nil
	}
}

// Set the timestamp format for the backup log filename.
// The default value is "20060102T150405.999999999MST" derived from [time.RFC3339Nano].
//
// The layout must be of a valid Go time layout, since this package use [time.Time.Format]
// it will not return an error if the layout is invalid,
// instead it will use whatever default layout that method is using.
// It should include nanosecond in order to avoid name conflict.
// See more about Go time layout at [time package constants].
//
// [time package constants]: https://pkg.go.dev/time#pkg-constants
func WithTimeFormat(layout string) Opt {
	return func(k *Keeper) (*Keeper, error) {
		k.timeFormat = layout
		return k, nil
	}
}

// Maximum size in bytes per log file.
// Keeper will rotate the log file if its size exceeds this value.
// Set this value to zero or negative will disable this feature.
// The default value is one [MB].
func WithMaxsize(size int) Opt {
	return func(k *Keeper) (*Keeper, error) {
		k.maxsize = size
		return k, nil
	}
}
