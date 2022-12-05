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
	path   string
	buffer []T
	cipher Cipher
	locker *sync.Mutex
}

func (k *kelvin[T]) lock() { k.locker.Lock() }
func (k *kelvin[T]) unlock() { k.locker.Unlock() }

func (k *kelvin[T]) commit(buffer []T) error {
	f, err := os.Create(k.path)
	if err != nil {
		return err
	}

	var content []byte
	if buffer == nil {
		content = []byte(emptyContent)
	} else {
		content, err = json.Marshal(buffer)
		if err != nil {
			return err
		}
	}

	_, err = f.Write(readyToWrite(content, k.cipher))
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

// Commit writes content to disk.
// Only useable for in-memory mode.
func (k *kelvin[T]) Commit() error {
	if k.IsNoWrite() {
		return errors.New("no write mode enabled")
	}

	if k.mode != InMemory {
		return errors.New("mode is not setted as in-memory")
	}

	k.lock()
	defer k.unlock()
	return k.commit(k.buffer)
}

// IsNoWrite reports Kelvin instance is NoWrite mode.
func (k *kelvin[T]) IsNoWrite() bool { return k.path == NoWrite }

func (k *kelvin[T]) decode() ([]T, error) {
	info, err := os.Stat(k.path)
	if err != nil {
		panic("buffering failed: path is not exist: " + k.path)
	}
	if info.IsDir() {
		panic("buffering failed: path is direcotry: " + k.path)
	}
	bytes, err := os.ReadFile(k.path)
	if err != nil {
		panic("buffering failed: " + err.Error())
	}
	var buffer []T
	bytes = readyToProcess(bytes, k.cipher)
	err = json.Unmarshal(bytes, &buffer)
	if err != nil {
		panic("buffering failed: " + err.Error())
	}
	return buffer, err
}

// buff reads disk content of Kelvin database into buffer.
func (k *kelvin[T]) buff() {
	buffer, err := k.decode()
	if err == nil {
		k.buffer = buffer
	}
}

func (k *kelvin[T]) getBufferCopy() (_ []T, err error) {
	k.lock()
	defer k.unlock()
	var kbuffer []T
	if k.mode == InMemory {
		kbuffer = k.buffer
	} else {
		kbuffer, err = k.decode()
		if err != nil {
			return nil, err
		}
	}
	buffer := make([]T, len(kbuffer))
	_ = copy(buffer, kbuffer)
	return buffer, err
}

func (k *kelvin[T]) push(buffer []T) error {
	k.lock()
	defer k.unlock()
	if k.mode == Strict {
		if !k.IsNoWrite() {
			return k.commit(buffer)
		}
	} else {
		k.buffer = buffer
	}
	return nil
}

// Insert inserts items to database content.
func (k *kelvin[T]) Insert(items ...T) error {
	buffer, err := k.getBufferCopy()
	if err != nil {
		return err
	}
	buffer = append(buffer, items...)
	return k.push(buffer)
}
