package fileutil

import (
	"errors"
	"log"
	"os"
)

func PathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		log.Fatal(err)
		return false
	}
}
