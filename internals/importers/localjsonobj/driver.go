package localjsonobj

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Driver represents the main importer driver
type Driver struct {
	dir string
}

// Open implements Importer#open
func (d *Driver) Open(dsn string) error {
	parts := strings.Split(dsn, "://")
	if len(parts) < 1 {
		return fmt.Errorf("invalid dsn format for (%s)", dsn)
	}

	info, err := os.Stat(parts[1])
	if err != nil {
		return fmt.Errorf("unable to open (%s) due to: %s", parts[1], err)
	}

	if !info.IsDir() {
		return fmt.Errorf("the provided path (%s) isn't a directory", parts[1])
	}

	d.dir = parts[1]

	return nil
}

// Import implements Importer#import
func (d *Driver) Import(loaderRegexp string) (<-chan map[string]interface{}, <-chan error, <-chan bool) {
	resultChan := make(chan map[string]interface{})
	errChan := make(chan error)
	doneChan := make(chan bool)

	re, err := regexp.Compile(loaderRegexp)
	if err != nil {
		errChan <- err
		doneChan <- true
		goto eof
	}

	go (func() {
		errChan <- filepath.Walk(
			d.dir,
			func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}

				if !re.MatchString(path) {
					return nil
				}

				file, err := os.Open(path)
				if err != nil {
					errChan <- err
					return nil
				}

				decoder := json.NewDecoder(file)

				for {
					var m map[string]interface{}

					if decoder.Decode(&m) == io.EOF {
						break
					}

					resultChan <- m
				}

				doneChan <- true

				return nil
			},
		)
	})()

eof:
	return resultChan, errChan, doneChan
}
