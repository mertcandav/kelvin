package kelvin

import (
	"errors"
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

func open[T any](path string, mode Mode, cipher Cipher) (k *kelvin[T]) {
	var t T
	tt := reflect.TypeOf(t)
	if tt.Kind() != reflect.Struct {
		if tt.Kind() != reflect.Pointer || tt.Elem().Kind() != reflect.Struct {
			panic("type error: kelvin supports only structures")
		}
	}

	if mode != InMemory && mode != Strict {
		panic("database connection failed: invalid kelvin mode")
	}

	if path != NoWrite {
		if filepath.Ext(path) != Ext {
			panic("database connection failed: path is not a kelvin database file")
		}

		info, err := os.Stat(path)
		if err != nil {
			// Create new Kelvin database because not exist.
			if !errors.Is(err, os.ErrNotExist) {
				panic("database creation failed: " + err.Error())
			}

			f, err := os.Create(path)
			if err != nil {
				panic("database creation failed: " + err.Error())
			}
			content := []byte(emptyContent)
			_, err = f.Write(readyToWrite(content, cipher))
			if err != nil {
				panic("database creation failed: " + err.Error())
			}
			defer func() { k.stream = f }()
		} else {
			if info.IsDir() {
				panic("database connection failed: path is directory")
			}

			const PERM = 0666
			f, err := os.OpenFile(path, os.O_RDWR, PERM)
			if err != nil {
				panic("database connection failed: " + err.Error())
			}

			defer func() {
				k.stream = f
				if mode == InMemory {
					k.buff()
				}
			}()
		}
	} else if mode == Strict {
		panic("database connection failed: no-write mode is cannot combined with strict mode")
	}

	k = new(kelvin[T])
	k.mode = mode
	k.cipher = cipher
	k.locker = new(sync.Mutex)
	return k
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
func OpenSafe[T any](path string, mode Mode, cipher Cipher) Kelvin[T] {
	return open[T](path, mode, cipher)
}

// Open same as OpenSafe, but not uses cipher.
func Open[T any](path string, mode Mode) Kelvin[T] {
	return open[T](path, mode, nil)
}

// OpenNW same as Open, but opens with NoWrite and InMemory.
func OpenNW[T any]() Kelvin[T] {
	return open[T](NoWrite, InMemory, nil)
}
