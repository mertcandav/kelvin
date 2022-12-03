package kelvin

import (
	"encoding/json"
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
