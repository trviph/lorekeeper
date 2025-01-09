package lorekeeper

import (
	"log"
	"sync"
	"testing"
)

func TestRegistry(t *testing.T) {
	keeper1, err := New(
		WithName("unique-name"),
		WithFolder("."),
		WithMaxSize(10*Mb),
		WithCron("* * * * *"),
		WithArchiveNameLayout("test-output-{{ .name }}{{.extension}}{{ .time }}"),
	)
	if err != nil {
		t.Errorf("expect no error but got %v", err)
	}

	keeper2, err := New(
		WithName("unique-name"),
		WithMaxSize(20*Mb),
		WithArchiveNameLayout("test-output-{{ .name }}-{{.extension}}{{.time}}"),
	)
	if err != nil {
		t.Errorf("expect no error but got %v", err)
	}

	if keeper1 != keeper2 {
		t.Errorf("expect to be the same instance")
	}
	if keeper1.maxSize != keeper2.maxSize && keeper1.maxSize == 20*Mb {
		t.Errorf("expect maxSize to be the same")
	}
	if keeper1.cronScheduler != nil {
		t.Errorf("expect cron to be stop")
	}

	// Expect no race
	var wg sync.WaitGroup
	wg.Add(2)
	run := func(k *Keeper, id int) {
		defer wg.Done()
		debugLogger := log.New(k, "[DEBUG] ", log.Lmsgprefix|log.LstdFlags|log.Llongfile)
		for i := 0; i < 1000; i++ {
			go debugLogger.Printf("[%d] flooding the log with debug information...", id)
		}
	}
	run(keeper1, 1)
	run(keeper2, 2)

	// Create new keeper
	keeper3, err := New(
		WithName("another-unique-name"),
		WithFolder("."),
		WithMaxSize(10*Mb),
	)
	if err != nil {
		t.Errorf("expect no error but got %v", err)
	}
	if (keeper1 == keeper3) || (keeper2 == keeper3) {
		t.Errorf("expect to be not the same instance")
	}
	wg.Add(1)
	run(keeper3, 3)

	wg.Wait()
}
