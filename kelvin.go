package kelvin

import (
	"encoding/json"
	"errors"
	"os"
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
	buffer []T
	cipher Cipher
}

func (k *kelvin[T]) commit() error {
	f, err := os.Create(k.path)
	if err != nil {
		return err
	}
	content := []byte(emptyContent)
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
	if k.path == NoWrite {
		return errors.New("no write mode enabled")
	}

	if k.mode != InMemory {
		return errors.New("mode is not setted as in-memory")
	}

	return k.commit()
}

// buff reads disk content of Kelvin database into buffer.
func (k *kelvin[T]) buff() {
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
	bytes = readyToProcess(bytes, k.cipher)
	err = json.Unmarshal(bytes, &k.buffer)
	if err != nil {
		panic("buffering failed: " + err.Error())
	}
}

// Insert inserts items to database content.
func (k *kelvin[T]) Insert(items ...T) {

}
