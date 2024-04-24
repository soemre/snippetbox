package main

import (
	"log"
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	http.FileSystem
}

func (nfs neuteredFileSystem) Open(name string) (http.File, error) {
	f, err := nfs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		log.Printf("\"%s\" stat operation failed: %v", name, err)
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(name, "index.html")
		if _, err := nfs.FileSystem.Open(index); err != nil {
			if err := f.Close(); err != nil {
				log.Printf("\"%s\" couldn't be closed: %v", name, err)
				return nil, err
			}
			return nil, err
		}
	}

	return f, nil
}
