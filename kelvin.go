package kelvin

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
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
const emptyContent = "[]"

// kelvin is kelvin database structure.
type kelvin[T any] struct {
	mode   byte
	stream *os.File
	buffer []T
	cipher Cipher
	locker *sync.Mutex
}

func (k *kelvin[T]) lock() { k.locker.Lock() }
func (k *kelvin[T]) unlock() { k.locker.Unlock() }

func (k *kelvin[T]) commit(buffer []T) {
	var content []byte
	var err error
	if buffer == nil {
		content = []byte(emptyContent)
	} else {
		content, err = json.Marshal(buffer)
		if err != nil {
			panic("comitting failed: " + err.Error())
		}
	}
	content = readyToWrite(content, k.cipher)
	const TRUNC = 0
	err = k.stream.Truncate(TRUNC)
	if err != nil {
		panic("comitting failed: " + err.Error())
	}
	_, err = k.stream.Seek(TRUNC, TRUNC)
	if err != nil {
		panic("comitting failed: " + err.Error())
	}
	_, err = k.stream.WriteAt(content, TRUNC)
	if err != nil {
		panic("comitting failed: " + err.Error())
	}
	err = k.stream.Sync()
	if err != nil {
		panic("comitting failed: " + err.Error())
	}
}

// Commit writes content to disk.
// Only useable for in-memory mode.
func (k *kelvin[T]) Commit() error {
	if k.IsNoWrite() {
		return errors.New("comitting failed: no-write mode enabled")
	}

	if k.mode != InMemory {
		return errors.New("comitting failed: mode is not setted as in-memory")
	}

	k.lock()
	k.commit(k.buffer)
	k.unlock()
	return nil
}

// IsNoWrite reports Kelvin instance is NoWrite mode.
func (k *kelvin[T]) IsNoWrite() bool { return k.stream == nil }

func (k *kelvin[T]) decode() []T {
	const TRUNC = 0
	const EOF   = 2
	n, err := k.stream.Seek(TRUNC, EOF)
	if err != nil {
		panic("buffering failed: " + err.Error())
	}
	bytes := make([]byte, n)
	_, err = k.stream.ReadAt(bytes, TRUNC)
	if err != nil {
		panic("buffering failed: " + err.Error())
	}
	var buffer []T
	bytes = readyToProcess(bytes, k.cipher)
	err = json.Unmarshal(bytes, &buffer)
	if err != nil {
		panic("buffering failed: " + err.Error())
	}
	return buffer
}

// buff reads disk content of Kelvin database into buffer.
func (k *kelvin[T]) buff() { k.buffer = k.decode() }

func (k *kelvin[T]) getCollection() []T {
	k.lock()
	defer k.unlock()
	if k.mode == Strict {
		return k.decode()
	}
	return k.buffer
}

func (k *kelvin[T]) getImmutableCollection() []T {
	buffer := k.getCollection()
	if k.mode == Strict {
		return buffer
	}
	cbuffer := make([]T, len(buffer))
	_ = copy(cbuffer, buffer)
	return cbuffer
}

func (k *kelvin[T]) push(buffer []T) {
	k.lock()
	defer k.unlock()
	if k.mode == Strict {
		if !k.IsNoWrite() {
			k.commit(buffer)
		}
	} else {
		k.buffer = buffer
	}
}

// Insert inserts items to database content.
func (k *kelvin[T]) Insert(items ...T) {
	buffer := k.getImmutableCollection()
	buffer = append(buffer, items...)
	k.push(buffer)
}

// GetCollection returns all collection.
// Not returns deep copy of collection.
func (k *kelvin[T]) GetCollection() []T { return k.getImmutableCollection() }
