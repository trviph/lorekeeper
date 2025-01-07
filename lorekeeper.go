package lorekeeper

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"text/template"
	"time"
)

// A [Keeper] is a log file manager that handles writing to log files and rotates them.
// Use [NewKeeper] to create a new Keeper.
type Keeper struct {
	// See [WithFolder] for documentation.
	folder string
	// See [WithName] for documentation.
	name string
	// See [WithExtension] for documentation.
	extension string
	// See [WithTimeLayout] for documentation
	timeLayout string
	// See [WithMaxsize] for documentation
	maxsize int
	// See [WithArchiveNameLayout] for documentation
	archiveNameLayout *template.Template

	fileMU      sync.Mutex
	currentFile io.WriteCloser
	currentSize int
}

// Make sure that keeper implements the [io.Writer] interface,
// so that it can be use with the [log] package.
var _ io.Writer = (*Keeper)(nil)

// Create a new [Keeper] with the provided options.
// This will create a [DefaultKeeper] if no option is provided.
// If at least one option is provided, this may also return an error if the option is invalid.
// See [Opt] for all available options.
//
// Example usage:
//
//		import "github.com/trviph/lorekeeper"
//
//		func main() {
//			keeper, err := lorekeeper.NewKeeper(
//				lorekeeper.WithName("Lorekeeper Example"),
//				lorekeeper.WithMaxByte(12 * lorekeeper.Kb),
//	 	)
//		}
func NewKeeper(opts ...Opt) (*Keeper, error) {
	defaultOpts := []Opt{
		WithFolder(os.TempDir()),
		WithName(defaultKeeperName()),
		WithExtension(".log"),
		WithTimeLayout("2006-01-02-15-04-05.000000000-0700"),
		WithMaxsize(15 * Mb),
		WithArchiveNameLayout("{{ .time }}-{{ .name }}{{ .extension }}"),
	}

	var err error
	keeper := new(Keeper)
	for _, opt := range append(defaultOpts, opts...) {
		keeper, err = opt(keeper)
		if err != nil {
			return nil, err
		}
	}

	file, err := keeper.getCurrentFile()
	if err != nil {
		return nil, err
	}
	keeper.currentFile = file

	return keeper, err
}

// Get the current log file descriptor.
func (k *Keeper) getCurrentFile() (*os.File, error) {
	return os.OpenFile(k.getCurrentFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// Get the path to the current log file.
func (k *Keeper) getCurrentFilePath() string {
	return path.Join(k.folder, fmt.Sprintf("%s%s", k.name, k.extension))
}

// Write the msg to the current log file.
func (k *Keeper) Write(msg []byte) (int, error) {
	k.fileMU.Lock()
	defer k.fileMU.Unlock()

	if k.shouldRotate(msg) {
		if err := k.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := k.currentFile.Write(msg)
	if err != nil {
		return 0, err
	}
	k.currentSize += n
	return n, nil
}

// Archive the current log file and create a new log file.
func (k *Keeper) rotate() error {
	// Close and rename the old file
	if err := k.currentFile.Close(); err != nil {
		return fmt.Errorf("failed to rotate log file, caused by %w", err)
	}
	archiveName, err := k.newArchiveName()
	if err != nil {
		return fmt.Errorf("failed to get new archive name, caused by %w", err)
	}
	if err := os.Rename(k.getCurrentFilePath(), archiveName); err != nil {
		return fmt.Errorf("failed to rotate log file, caused by %w", err)
	}

	// Create a new file
	file, err := k.getCurrentFile()
	if err != nil {
		return err
	}
	k.currentFile = file
	k.currentSize = 0
	return nil
}

func (k *Keeper) newArchiveName() (string, error) {
	var buff bytes.Buffer
	err := k.archiveNameLayout.Execute(
		&buff,
		map[string]any{
			"time":      time.Now().Format(k.timeLayout),
			"name":      k.name,
			"extension": k.extension,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to execute template, caused by %w", err)
	}
	return path.Join(k.folder, buff.String()), nil
}

func (k *Keeper) shouldRotate(nextMsg []byte) bool {
	return k.maxsize > 0 && k.currentSize+len(nextMsg) > k.maxsize
}
