package driver

import (
	"fmt"
	"net/url"
	"sync"
)

type Driver interface {
	Open(dsn string) error
	Run(loader string) (<-chan map[string]interface{}, <-chan error, <-chan bool)
}

var (
	drivers     = map[string]Driver{}
	driversLock = &sync.RWMutex{}
)

func Register(name string, driver Driver) error {
	driversLock.Lock()
	defer driversLock.Unlock()

	if _, exists := drivers[name]; exists {
		return fmt.Errorf("duplicate driver name %s", name)
	}

	drivers[name] = driver

	return nil
}

func Open(dsn string) (Driver, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	driversLock.Lock()
	defer driversLock.Unlock()

	driver, exists := drivers[u.Scheme]

	if !exists {
		return nil, fmt.Errorf("driver %s not found", u.Scheme)
	}

	return driver, driver.Open(dsn)
}
