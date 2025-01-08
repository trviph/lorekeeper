package lorekeeper

import (
	"fmt"
	"os"
	"path"
	"time"
)

var now func() time.Time = time.Now

// Get default name for the [Keeper].
func defaultKeeperName() string {
	if len(os.Args) > 1 && len(os.Args[0]) > 1 {
		_, execName := path.Split(os.Args[0])
		return fmt.Sprintf("lorekeeper-%s", execName)
	}
	return "lorekeeper"
}
