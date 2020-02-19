package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Linter struct {
	Entrypoint string
	Errors     []*Error
	*sync.RWMutex
}

func (linter *Linter) getEntrypoint() string {
	linter.RLock()
	defer linter.RUnlock()

	return linter.Entrypoint
}

func (linter *Linter) getErrors() []*Error {
	linter.RLock()
	defer linter.RUnlock()

	return linter.Errors
}

func (linter *Linter) addError(error *Error) {
	linter.Lock()
	defer linter.Unlock()

	linter.Errors = append(linter.Errors, error)
}

func (linter *Linter) validateDir(config *Config, index index, path string) error {
	rules := config.getConfig(index, path)
	basename := filepath.Base(path)

	log.Printf("%s %s %+v", basename, path, rules[".dir"])
	return nil
}

func (linter *Linter) validateFile(config *Config, index index, entrypoint string, path string) error {
	ext := filepath.Ext(path)
	rules := config.getConfig(index, path)
	withoutExt := strings.TrimSuffix(filepath.Base(path), ext)

	log.Printf("%s %s %+v", ext, withoutExt, rules[ext])
	return nil
}

func (linter *Linter) Run(config *Config) error {
	var g = new(errgroup.Group)
	var ls = config.getLs()
	var index, err = config.getIndex(ls)

	if err != nil {
		return err
	}

	for entrypoint := range ls {
		g.Go(func() error {
			return filepath.Walk(entrypoint, func(path string, info os.FileInfo, err error) error {
				if info == nil {
					return fmt.Errorf("%s not found", entrypoint)
				}

				if info.IsDir() {
					return linter.validateDir(config, index, path)
				}

				return linter.validateFile(config, index, entrypoint, path)
			})
		})
	}

	return g.Wait()
}
