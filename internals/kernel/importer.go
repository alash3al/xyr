package kernel

import (
	"fmt"
	"sync"
)

// Importer is a contract for all drivers that imports data into xyr internal storage
type Importer interface {
	Open(dsn string) error
	Import(loader string) (<-chan map[string]interface{}, <-chan error, <-chan bool)
}

var (
	importers     = map[string]Importer{}
	importersLock = &sync.RWMutex{}
)

// RegisterImporter registers the given importer into the global imports registry
func RegisterImporter(driverName string, driver Importer) error {
	importersLock.Lock()
	defer importersLock.Unlock()

	if _, exists := importers[driverName]; exists {
		return fmt.Errorf("duplicate driver name %s", driverName)
	}

	importers[driverName] = driver

	return nil
}

// OpenImporter opens the specified importer via its dsn
func OpenImporter(driverName, dsn string) (Importer, error) {
	importersLock.Lock()
	defer importersLock.Unlock()

	driver, exists := importers[driverName]

	if !exists {
		return nil, fmt.Errorf("driver %s not found", driverName)
	}

	return driver, driver.Open(dsn)
}
