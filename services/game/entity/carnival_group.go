package entity

import (
	"sync"
)

var carnivalRegistryMu sync.RWMutex
var carnivalRegistryInstance *CarnivalRegistry

func GlobalCarnivalRegistry() *CarnivalRegistry {
	carnivalRegistryMu.RLock()
	reg := carnivalRegistryInstance
	carnivalRegistryMu.RUnlock()
	if reg != nil {
		return reg
	}
	carnivalRegistryMu.Lock()
	defer carnivalRegistryMu.Unlock()
	if carnivalRegistryInstance == nil {
		carnivalRegistryInstance = NewCarnivalRegistry()
	}
	return carnivalRegistryInstance
}
