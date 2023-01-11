package common

import (
	"sync"
	"sync/atomic"
)

var storeInitialized uint32

var dataStore DataStore

var dataStoreMu sync.Mutex

type DataStore struct {
	configFileName string
	baiduYunConfig BaiduYunConfig
	listenedEvents []ListenedEvent
}

func GetDataStore() DataStore {
	if atomic.LoadUint32(&storeInitialized) == 1 {
		return dataStore
	}
	dataStoreMu.Lock()
	defer dataStoreMu.Unlock()

	if storeInitialized == 0 {
		dataStore = DataStore{}
		atomic.StoreUint32(&storeInitialized, 1)
	}
	return dataStore
}
