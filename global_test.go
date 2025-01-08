package lorekeeper

import (
	"log"
	"sync"
	"testing"
)

func TestRegistry(t *testing.T) {
	keeper1, err := NewKeeper(
		WithName("unique-name"),
		WithFolder("."),
		WithMaxSize(10*Mb),
	)
	if err != nil {
		t.Errorf("expect no error but got %v", err)
	}

	keeper2, err := NewKeeper(
		WithName("unique-name"),
		WithMaxSize(20*Mb),
	)
	if err != nil {
		t.Errorf("expect no error but got %v", err)
	}

	// Now both keeper1 and keeper2 will have maxSize of 20*MB.
	if keeper1 != keeper2 {
		t.Errorf("expect to be the same instance")
	}
	if keeper1.maxSize != keeper2.maxSize {
		t.Errorf("expect maxSize to be the same")
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
	keeper3, err := NewKeeper(
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
