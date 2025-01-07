package lorekeeper

import (
	"fmt"
	"os"
)

// Get default name for the [Keeper].
func defaultKeeperName() string {
	if len(os.Args) > 1 && len(os.Args[0]) > 2 {
		return fmt.Sprintf("lorekeeper-%s", os.Args[0][2:])
	}
	return "lorekeeper"
}
