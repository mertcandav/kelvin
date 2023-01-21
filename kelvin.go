package kelvin

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"sync"
)

// Mode is types of the kelvin database modes.
type Mode byte

// Kelvin mode.
const (
	InMemory = Mode(1)
	Strict   = Mode(2)
)

// kelvin no-write-to-disk mode.
const NoWrite = ""

// Kelvin database file extension.
const Ext = ".klvn"

// Empty Kelvin content.
const emptyContent = "[]"

// Kelvin is an interface for static typing.
type Kelvin[T any] interface {
	Commit() error
	IsNoWrite() bool
	Fill(...T)
	Insert(...T)
	Drop(...T)
	DropWhere(func(T) bool)
	GetCollection() []T
	Map(func(*T))
	Where(func(*T) bool) []T
	UWhere(func(*T) bool) []T
}

// kelvin is kelvin database structure.
type kelvin[T any] struct {
	mode   Mode
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

// Lock needed.
func (k *kelvin[T]) getCollection() []T {
	if k.mode == Strict {
		return k.decode()
	}
	return k.buffer
}

// Lock needed.
func (k *kelvin[T]) getImmutableCollection() []T {
	buffer := k.getCollection()
	if k.mode == Strict {
		return buffer
	}
	n := len(buffer)
	cbuffer := make([]T, n)
	for i := 0; i < n; i++ {
		cbuffer[i] = deepCopy(buffer[i])
	}
	return cbuffer
}

// Lock needed.
func (k *kelvin[T]) push(buffer []T) {
	if k.mode == Strict {
		if !k.IsNoWrite() {
			k.commit(buffer)
		}
	} else {
		k.buffer = buffer
	}
}

// Fill removes all datas and inserts given datas.
// Fills each item by deep immutable copy.
func (k *kelvin[T]) Fill(items ...T) {
	k.lock()
	defer k.unlock()

	if len(items) == 0 {
		k.push(nil)
		return
	}

	if k.mode == Strict {
		k.push(items)
		return
	}

	buffer := make([]T, len(items))
	for i, item := range items {
		buffer[i] = deepCopy(item)
	}

	k.push(buffer)
}

// Insert inserts items to database content.
func (k *kelvin[T]) Insert(items ...T) {
	k.lock()
	defer k.unlock()

	buffer := k.getCollection()
	buffer = append(buffer, items...)
	k.push(buffer)
}

// Drop removes items from database content.
// Drops data if all fields (both exported and unexported) deeply equals.
func (k *kelvin[T]) Drop(items ...T) {
	k.lock()
	defer k.unlock()

	buffer := k.getCollection()
	changed := false
	for _, item := range items {
		for i, data := range buffer {
			if reflect.DeepEqual(data, item) {
				buffer[i] = buffer[len(buffer)-1]
				buffer = buffer[:len(buffer)-1]
				changed = true
			}
		}
	}
	if changed {
		k.push(buffer)
	}
}

// Drop removes items from database content.
// Drops data if handler function returns true.
// Does not nothing if handler is nil.
func (k *kelvin[T]) DropWhere(handler func(T) bool) {
	if handler == nil {
		return
	}

	k.lock()
	defer k.unlock()

	buffer := k.getCollection()
	changed := false
	for i := 0; i < len(buffer); {
		item := buffer[i]
		if handler(item) {
			buffer[i] = buffer[len(buffer)-1]
			buffer = buffer[:len(buffer)-1]
			changed = true
			continue
		}
		i++
	}
	if changed {
		k.push(buffer)
	}
}

// GetCollection returns all collection.
func (k *kelvin[T]) GetCollection() []T {
	k.lock()
	defer k.unlock()
	return k.getImmutableCollection()
}

// Map iterates into all collection and commits changes.
// Does not nothing if handler is nil.
func (k *kelvin[T]) Map(handler func(*T)) {
	if handler == nil {
		return
	}

	k.lock()
	defer k.unlock()

	buffer := k.getCollection()
	if len(buffer) == 0 {
		return
	}

	for i := 0; i < len(buffer); i++ {
		element := &buffer[i]
		handler(element)
	}

	k.push(buffer)
}

func (k *kelvin[T]) where(buffer []T, handler func(*T) bool) []T {
	result := make([]T, 0, len(buffer)/2)
	for i := 0; i < len(buffer); i++ {
		element := &buffer[i]
		if handler(element) {
			result = append(result, *element)
		}
	}

	return result
}

// Where returns a collection containing only data for which the handler returns true.
// Returns nil if handler is nil.
func (k *kelvin[T]) Where(handler func(*T) bool) []T {
	if handler == nil {
		return nil
	}

	k.lock()
	buffer := k.getImmutableCollection()
	k.unlock()

	return k.where(buffer, handler)
}

// UWhere returns a collection containing only data for which the handler returns true.
// Returns nil if handler is nil.
//
// This function is unsafe.
// You can change original buffer of Kelvin instance becuase
// this function doesn't use immutable copy of buffer.
// Therefore this function theorically fast than Where function
// but not safe as Where.
func (k *kelvin[T]) UWhere(handler func(*T) bool) []T {
	if handler == nil {
		return nil
	}

	k.lock()
	buffer := k.getCollection()
	k.unlock()

	return k.where(buffer, handler)
}
