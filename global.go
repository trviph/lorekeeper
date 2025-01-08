package lorekeeper

import "sync"

// Keeping track of all Keeper instances by their name.
var registry *sync.Map = new(sync.Map)

// Register the Keeper to the registry if it's not yet created,
// else return the registered one.
func register(name string, keeper *Keeper) (k *Keeper, new bool) {
	val, loaded := registry.LoadOrStore(name, keeper)
	return val.(*Keeper), !loaded
}

// Unregister the Keeper of a given name.
func unregister(name string) {
	registry.Delete(name)
}
