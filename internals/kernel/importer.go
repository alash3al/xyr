package kernel

import (
	"fmt"
	"net/url"
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
func RegisterImporter(name string, driver Importer) error {
	importersLock.Lock()
	defer importersLock.Unlock()

	if _, exists := importers[name]; exists {
		return fmt.Errorf("duplicate driver name %s", name)
	}

	importers[name] = driver

	return nil
}

// OpenImporter opens the specified importer via its dsn
func OpenImporter(dsn string) (Importer, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	importersLock.Lock()
	defer importersLock.Unlock()

	driver, exists := importers[u.Scheme]

	if !exists {
		return nil, fmt.Errorf("driver %s not found", u.Scheme)
	}

	return driver, driver.Open(dsn)
}
