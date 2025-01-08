package lorekeeper

import "sync"

// Keeping tracks if all Keeper instances by name
var registry *sync.Map = new(sync.Map)

// Register the Keeper to the registry if it not yet created,
// else return the registered one.
func register(name string, keeper *Keeper) (k *Keeper, new bool) {
	val, loaded := registry.LoadOrStore(name, keeper)
	return val.(*Keeper), !loaded
}
