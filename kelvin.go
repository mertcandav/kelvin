package kelvin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// Kelvin mode.
const (
	InMemory = 1
	Strict   = 2
)

// kelvin no-write-to-disk mode.
const NoWrite = ""

// Kelvin database file extension.
const Ext = ".klvn"

// Empty Kelvin content.
const emptyContent = "{}"

// kelvin is kelvin database structure.
type kelvin[T any] struct {
	mode   byte
	path   string
	stack  []T
	cipher Cipher
}

func readyToWrite(content []byte, cipher Cipher) []byte {
	if cipher != nil {
		content = cipher.Encrypt(content)
	}
	return content
}

// OpenSafe returns new instance of Kelvin by data type.
// Panics if T is not structure.
// Returns error if any errors occur.
// Creates new Kelvin database if not exist in given path.
// Uses given cipher for encryption.
func OpenSafe[T any](path string, mode int, cipher Cipher) (*kelvin[T], error) {
	var t T
	kind := reflect.TypeOf(t).Kind()
	if kind != reflect.Struct {
		panic("type error: kelvin supports only structures")
	}

	if mode != InMemory && mode != Strict {
		return nil, fmt.Errorf("'%d' is invalid kelvin mode", mode)
	}

	if path != NoWrite {
		if filepath.Ext(path) != Ext {
			return nil, fmt.Errorf(`"%s" path is not a kelvin database file`, path)
		}

		info, err := os.Stat(path)
		if err != nil {
			// Create new Kelvin database because not exist.

			if !errors.Is(err, os.ErrNotExist) {
				return nil, err
			}
			f, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			content := []byte(emptyContent)
			_, err = f.Write(readyToWrite(content, cipher))
			if err != nil {
				return nil, err
			}
			err = f.Close()
			if err != nil {
				return nil, err
			}
		} else {
			if info.IsDir() {
				return nil, fmt.Errorf(`"%s" path is directory`, path)
			}
		}
	}

	k := new(kelvin[T])
	k.path = path
	k.stack = nil
	k.mode = byte(mode)
	k.cipher = cipher
	return k, nil
}

// Open same as OpenSafe, but not uses cipher.
func Open[T any](path string, mode int) (*kelvin[T], error) {
	return OpenSafe[T](path, mode, nil)
}
