package lorekeeper

import (
	"fmt"
	"strings"
	"text/template"
)

// An Opt is a function that mutates a [Keeper]'s attributes.
// An Opt should return a mutated Keeper or return an error if it fails to mutate the Keeper.
// An Opt should be used together with [NewKeeper].
type Opt func(*Keeper) (*Keeper, error)

// The folder where the log files are stored.
// The default value is [os.TempDir].
func WithFolder(path string) Opt {
	return func(k *Keeper) (*Keeper, error) {
		if len(path) > 0 {
			k.folder = path
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
			k.name = strings.ReplaceAll(strings.ToLower(name), " ", "-")
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

// Set the timestamp layout for the backup log filename.
// The default value is "2006-01-02-15-04-05.000000000-0700".
//
// The layout must be of a valid Go time layout, since this package use [time.Time.Format]
// it will not return an error if the layout is invalid,
// instead it will use whatever default layout that method is using.
//
// The layout should include nanosecond in order to avoid name conflict. Upon name conflict, the new log file will replace the old log.
//
// See more about Go time layout at [time package constants].
//
// [time package constants]: https://pkg.go.dev/time#pkg-constants
func WithTimeLayout(layout string) Opt {
	return func(k *Keeper) (*Keeper, error) {
		k.timeLayout = layout
		return k, nil
	}
}

// Maximum size in bytes per log file.
// Keeper will rotate the log file if its size exceeds this value.
// Set this value to zero or negative will disable this feature.
// The default value is 15[Mb].
func WithMaxsize(size int) Opt {
	return func(k *Keeper) (*Keeper, error) {
		k.maxsize = size
		return k, nil
	}
}

// Set the filename layout for the archived log file.
// If the layout is empty will use the default value.
// The default value is "{{ .time }}-{{ .name }}{{ .extension }}".
// The layout is parsed using the [text/template] package.
// The supported arguments are:
//   - {{ .time }} the time when the rotation happened.
//   - {{ .name }} the name of the Keeper.
//   - {{ .extension }} the extension of the file.
func WithArchiveNameLayout(layout string) Opt {
	return func(k *Keeper) (*Keeper, error) {
		if len(layout) == 0 {
			return k, nil
		}
		templ, err := template.New("lorekeeper-archive-template").Parse(layout)
		if err != nil {
			return nil, fmt.Errorf("failed to set archive name layout, caused by %w", err)
		}
		k.archiveNameLayout = templ
		return k, nil
	}
}
