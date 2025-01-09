package lorekeeper

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// This test demonstrates on how to use [lorekeeper.Keeper] with the std [log].
func Test(t *testing.T) {
	// Create a Keeper
	keeper, err := New(
		// Set the Keeper name, this will be used when generate log files.
		WithName("Test"),
		// Set the extension of archived logs.
		WithExtension(".log"),
		// Set the time layout of archived logs.
		WithTimeLayout("20060102150405.000"),
		// Specify the folder where the log files will be stored.
		WithFolder("."),
		// Each log file hold a maximum of 50 Kibibyte before being rotated.
		WithMaxSize(50*Kb),
		// Set the name layout of archived logs.
		WithArchiveNameLayout("test-output-{{ .name }}{{.extension}}{{.time}}"),
		// Set the maximum number of archives to keep.
		WithMaxFiles(2),
		// Set archives to be compressed with gzip
		WithGzip(),
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

func BenchmarkKeeperWrite(b *testing.B) {
	const lorem = "Culpa sequi esse et et expedita aut qui quia. Error minus modi sunt beatae asperiores qui rem. Quia minima cumque laudantium sed rerum. Sunt delectus nesciunt dolor veniam soluta provident porro deserunt. Ullam illo beatae et quos unde maxime repellendus. Beatae itaque totam eum itaque velit et. Sit molestias dolore deserunt rerum amet. Molestiae rem provident minima autem nulla numquam. Illum voluptas ea nam suscipit. Corporis molestias necessitatibus dolore facilis. Nostrum cum nemo vero. Enim dolorem esse ad. Sed numquam odio eum ex. Praesentium incidunt quod perferendis sit est omnis sapiente. Sed rem itaque laboriosam minus eos. Sed fugiat dolores ut. Nam veniam nihil voluptatem accusamus molestias ducimus. Minima aut consequuntur dolores facere inventore libero tempore omnis. Suscipit et aut nostrum. Porro sapiente dignissimos nisi error. Et nulla vel molestiae veniam molestiae eum. Est similique sapiente aperiam voluptate cum occaecati et laboriosam. Praesentium cupiditate et laboriosam aperiam neque ut ut. Provident blanditiis autem pariatur autem animi et sint dicta."
	k, _ := New(
		WithFolder("."),
		WithMaxSize(100*KB),
		WithMaxFiles(1),
		WithGzip(),
		WithName("BenchmarkKeeperWrite"),
		WithArchiveNameLayout("test-output-{{ .name }}{{.extension}}{{.time}}"),
	)
	logger := log.New(k, "[Benchmark] ", log.LstdFlags|log.Lmsgprefix)
	b.Run(
		"Write to Keeper",
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				logger.Println(lorem)
			}
		},
	)
}

func TestKeeperNewArchiveName(t *testing.T) {
	now = func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		return t
	}

	defaultOpts := []Opt{
		WithTimeLayout("20060102"),
		WithName("testcase 1"),
		WithExtension(".log"),
		WithArchiveNameLayout("log.old"),
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		opts    []Opt
		want    string
		wantErr bool
	}{
		{
			name: "with fixed name",
			opts: defaultOpts,
			want: filepath.Join(os.TempDir(), "log.old"),
		},
		{
			name: "with configured name",
			opts: append(
				defaultOpts,
				WithName("testcase 2"),
				WithArchiveNameLayout("{{ .name }}.old"),
			),
			want: filepath.Join(os.TempDir(), "testcase-2.old"),
		},
		{
			name: "with configured name, extension",
			opts: append(
				defaultOpts,
				WithName("testcase 3"),
				WithArchiveNameLayout("{{ .name }}{{ .extension }}"),
			),
			want: filepath.Join(os.TempDir(), "testcase-3.log"),
		},
		{
			name: "with configured name, extension, and time",
			opts: append(
				defaultOpts,
				WithName("testcase 4"),
				WithArchiveNameLayout("{{ .name }}{{ .extension }}{{ .time }}"),
			),
			want: filepath.Join(os.TempDir(), "testcase-4.log20060102"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := New(tt.opts...)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := k.newArchiveName()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("newArchiveName() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("newArchiveName() succeeded unexpectedly")
			}
			if tt.want != got {
				t.Errorf("newArchiveName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeperGetArchiveBlobPattern(t *testing.T) {
	defaultOpts := []Opt{
		WithTimeLayout("20060102"),
		WithName("testcase 1"),
		WithExtension(".log"),
		WithArchiveNameLayout("log.old"),
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		opts    []Opt
		want    string
		wantErr bool
	}{
		{
			name: "with fixed name",
			opts: defaultOpts,
			want: filepath.Join(os.TempDir(), "log.old"),
		},
		{
			name: "with configured name",
			opts: append(
				defaultOpts,
				WithName("testcase 2"),
				WithArchiveNameLayout("{{ .name }}.old"),
			),
			want: filepath.Join(os.TempDir(), "testcase-2.old"),
		},
		{
			name: "with configured name, extension",
			opts: append(
				defaultOpts,
				WithName("testcase 3"),
				WithArchiveNameLayout("{{ .name }}{{ .extension }}"),
			),
			want: filepath.Join(os.TempDir(), "testcase-3.log"),
		},
		{
			name: "with configured name, extension, and time",
			opts: append(
				defaultOpts,
				WithName("testcase 4"),
				WithArchiveNameLayout("{{ .name }}{{ .extension }}{{ .time }}"),
			),
			want: filepath.Join(os.TempDir(), "testcase-4.log*"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := New(tt.opts...)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := k.getArchiveGlobPattern()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getArchiveBlobPattern() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getArchiveBlobPattern() succeeded unexpectedly")
			}
			if tt.want != got {
				t.Errorf("getArchiveBlobPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeperGetCurrentFilePath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		opts []Opt
		want string
	}{
		{
			name: "default configuration",
			want: filepath.Join(
				os.TempDir(), "lorekeeper-lorekeeper.test.log",
			),
		},
		{
			name: "configured name",
			opts: []Opt{
				WithName("test"),
			},
			want: filepath.Join(
				os.TempDir(), "test.log",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := New(tt.opts...)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := k.getCurrentFilePath()
			if tt.want != got {
				t.Errorf("getCurrentFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeperClose(t *testing.T) {
	k, err := New(
		WithName("Test-Close"),
	)
	if _, ok := registry.Load(k.name); !ok {
		t.Error("expected the keeper to be registered")
	}

	if err != nil {
		t.Errorf("expected no error got %v", err)
	}
	if err := k.Close(); err != nil {
		t.Errorf("expected no error got %v", err)
	}
	if _, err := k.Write([]byte{}); err == nil {
		t.Errorf("expected error since Keeper is close got %v", err)
	}
	if val, ok := registry.Load(k.name); ok {
		t.Errorf("expected the keeper to be gone from the registry but got %v", val)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		opts    []Opt
		wantErr bool
	}{
		{
			name: "default",
		},
		{
			name: "fully configured",
			opts: []Opt{
				WithFolder("."),
				WithName("fully configured"),
				WithExtension(".log"),
				WithMaxSize(10),
				WithMaxFiles(5),
				WithTimeLayout("20060102"),
				WithArchiveNameLayout("{{ .time }}{{ .extension }}"),
				WithCron("* * * * *"),
				WithGzip(),
			},
		},
		{
			name: "empty extension",
			opts: []Opt{
				WithExtension(""),
			},
		},
		{
			name: "folder not existed",
			opts: []Opt{
				WithFolder("/lorem-ipsum-jada-jada"),
			},
			wantErr: true,
		},
		{
			name: "invalid archive name template",
			opts: []Opt{
				WithArchiveNameLayout("{{ time }}{{ name }}{{ .extension }}"),
			},
			wantErr: true,
		},
		{
			name: "invalid cron spec",
			opts: []Opt{
				WithCron("999 * * * *"),
			},
			wantErr: true,
		},
		{
			name: "invalid gzip level",
			opts: []Opt{
				WithGzipLevel(1000),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := New(tt.opts...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("New() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("New() succeeded unexpectedly")
			}
		})
	}
}
