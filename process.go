package kelvin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

func readyToWrite(bytes []byte, cipher Cipher) []byte {
	if cipher != nil {
		bytes = cipher.Encrypt(bytes)
	}
	return bytes
}

func readyToProcess(bytes []byte, cipher Cipher) []byte {
	if cipher != nil {
		bytes = cipher.Decrpyt(bytes)
	}
	return bytes
}

func open[T any](path string, mode int, cipher Cipher) (k *kelvin[T], _ error) {
	var t T
	kind := reflect.TypeOf(t).Kind()
	if kind != reflect.Struct {
		panic("type error: kelvin supports only structures")
	}

	if mode != InMemory && mode != Strict {
		return nil, fmt.Errorf("'%d' is invalid kelvin mode", mode)
	}

	var buffer []T
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

			if mode == InMemory {
				defer func () { k.buff() }()
			}
		}
	}

	k = new(kelvin[T])
	k.path = path
	k.buffer = buffer
	k.mode = byte(mode)
	k.cipher = cipher
	k.locker = new(sync.Mutex)
	return k, nil
}

// OpenSafe returns new instance of Kelvin by data type.
// Returns error if any errors occur.
// Creates new Kelvin database if not exist in given path.
// Uses given cipher for encryption.
//
// Panics if;
//  - T is not structure
//  - decoding is failed
//  - buffering is failed
func OpenSafe[T any](path string, mode int, cipher Cipher) (*kelvin[T], error) {
	return open[T](path, mode, cipher)
}

// Open same as OpenSafe, but not uses cipher.
func Open[T any](path string, mode int) (*kelvin[T], error) {
	return open[T](path, mode, nil)
}
