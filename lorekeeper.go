package lorekeeper

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"text/template"

	"github.com/robfig/cron/v3"
	"github.com/trviph/collection"
)

// A [Keeper] is a log file manager that handles writing to log files and rotates them.
// Use [New] to create a new Keeper.
type Keeper struct {
	// See [WithFolder] for documentation.
	folder string
	// See [WithName] for documentation.
	name string
	// See [WithExtension] for documentation.
	extension string
	// See [WithTimeLayout] for documentation
	timeLayout string
	// See [WithMaxSize] for documentation
	maxSize int
	// See [WithArchiveNameLayout] for documentation
	archiveNameLayout *template.Template
	// See [WithMaxFiles] for documentation
	maxFiles int
	// See [WithCron] for documentation
	cronScheduler *cron.Cron
	cronEntryID   cron.EntryID
	// See [WithGzip], [WithGzipLevel] for documentation
	compressorContructor func(w io.Writer) (io.WriteCloser, error)
	compressionExt       string

	mu              sync.Mutex
	currentFile     io.WriteCloser
	currentFileSize int

	archives *collection.List[string]
}

// Make sure that keeper implements the [io.Writer] interface,
// so that it can be use with the [log] package.
var _ io.WriteCloser = (*Keeper)(nil)

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
//			keeper, err := lorekeeper.New(
//				lorekeeper.WithName("Lorekeeper Example"),
//				lorekeeper.WithMaxByte(12 * lorekeeper.Kb),
//	 	)
//		}
func New(opts ...Opt) (*Keeper, error) {
	defaultOpts := []Opt{
		WithFolder(os.TempDir()),
		WithName(defaultKeeperName()),
		WithExtension(".log"),
		WithTimeLayout("2006-01-02-15-04-05.000000000-0700"),
		WithMaxSize(15 * Mb),
		WithArchiveNameLayout("{{ .time }}-{{ .name }}{{ .extension }}"),
		WithMaxFiles(0),
		NoCron(),
		NoCompression(),
	}
	finalOpts := append(defaultOpts, opts...)

	keeper := new(Keeper)
	if err := keeper.applyOpts(finalOpts...); err != nil {
		return nil, fmt.Errorf("failed to create new keeper, caused by %w", err)
	}

	keeper, new := register(keeper.name, keeper)
	// If loaded old keeper from registry, update it configurations
	if !new {
		keeper.mu.Lock()
		defer keeper.mu.Unlock()
		if err := keeper.applyOpts(finalOpts...); err != nil {
			return nil, fmt.Errorf("failed to create new keeper, caused by %w", err)
		}
	}

	return keeper, nil
}

func (k *Keeper) applyOpts(opts ...Opt) error {
	var err error
	for _, opt := range opts {
		k, err = opt(k)
		if err != nil {
			return fmt.Errorf("failed to apply option, caused by %w", err)
		}
	}

	file, err := k.getCurrentFile()
	if err != nil {
		return fmt.Errorf("failed to apply option, caused by %w", err)
	}
	k.currentFile = file
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to apply option, caused by %w", err)
	}
	k.currentFileSize = int(stat.Size())

	if k.maxFiles > 0 {
		archived, err := k.getArchives()
		if err != nil {
			return fmt.Errorf("failed to apply option, caused by %w", err)
		}
		k.archives = archived
	}

	return nil
}

func (k *Keeper) getArchives() (*collection.List[string], error) {
	pattern, err := k.getArchiveGlobPattern()
	if err != nil {
		return nil, fmt.Errorf("failed to get archive pattern, caused by %w", err)
	}
	return getArchives(pattern)
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
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.shouldRotate(msg) {
		if err := k.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := k.currentFile.Write(msg)
	if err != nil {
		return 0, err
	}
	k.currentFileSize += n
	return n, nil
}

// Rotate the current log file and close the Keeper.
// Any subsequence writes after this may cause error.
func (k *Keeper) Close() error {
	k.mu.Lock()
	defer k.mu.Unlock()
	// Rotate the log
	if err := k.rotate(); err != nil {
		return fmt.Errorf("failed to rotate file, caused by %w", err)
	}
	// Remove this Keeper from the registry
	unregister(k.name)
	if k.cronScheduler != nil {
		// Stop the cron scheduler to prevent goroutine leak
		k.cronScheduler.Stop()
	}
	// Close the file newly created file from rotate
	return k.currentFile.Close()
}

// Rotate to a new file immediately without waiting for the rotation conditions to be met.
func (k *Keeper) Rotate() error {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.rotate()
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
	k.currentFileSize = 0

	// Compress if set
	if k.compressorContructor != nil {
		if err := k.compress(archiveName); err != nil {
			return fmt.Errorf("failed to compressed rotated log")
		}
		archiveName += k.compressionExt
	}

	// Remove oldest archive
	if k.maxFiles > 0 {
		k.archives.Append(archiveName)
		for k.archives.Length() > k.maxFiles {
			oldest, err := k.archives.Dequeue()
			if err != nil {
				return fmt.Errorf("failed to get oldest archive name, caused by %w", err)
			}
			if err := os.Remove(oldest); err != nil {
				return fmt.Errorf("failed to remove oldest archive name, caused by %w", err)
			}
		}
	}

	return nil
}

func (k *Keeper) compress(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open file, caused by %w", err)
	}
	defer f.Close()

	cf, err := os.OpenFile(name+k.compressionExt, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create compressed file, caused by %w", err)
	}
	defer cf.Close()

	compressor, err := k.compressorContructor(cf)
	if err != nil {
		return fmt.Errorf("failed to create compress algorithm, caused by %w", err)
	}
	defer compressor.Close()

	if _, err := f.WriteTo(compressor); err != nil {
		return fmt.Errorf("failed to write to compressed file, caused by %w", err)
	}

	if err := os.Remove(name); err != nil {
		return fmt.Errorf("failed to delete %s, caused by %w", name, err)
	}
	return nil
}

func (k *Keeper) newArchiveName() (string, error) {
	var buff bytes.Buffer
	err := k.archiveNameLayout.Execute(
		&buff,
		map[string]any{
			"time":      now().Format(k.timeLayout),
			"name":      k.name,
			"extension": k.extension,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to execute template, caused by %w", err)
	}
	return path.Join(k.folder, buff.String()), nil
}

func (k *Keeper) getArchiveGlobPattern() (string, error) {
	var buff bytes.Buffer
	err := k.archiveNameLayout.Execute(
		&buff,
		map[string]any{
			"time":      "*",
			"name":      k.name,
			"extension": k.extension,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to execute template, caused by %w", err)
	}

	pattern := buff.String()
	if k.compressorContructor != nil {
		return path.Join(k.folder, pattern+k.compressionExt), nil
	}
	return path.Join(k.folder, pattern), nil
}

func (k *Keeper) shouldRotate(nextMsg []byte) bool {
	return k.maxSize > 0 && k.currentFileSize+len(nextMsg) > k.maxSize
}
