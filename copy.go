package kelvin

import (
	"bytes"
	"encoding/gob"
)

func deepCopy[T any](t T) T {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(&t)
	if err != nil {
		panic("deep copy failed: " + err.Error())
	}
	dec := gob.NewDecoder(&buffer)
	var c T
	err = dec.Decode(&c)
	if err != nil {
		panic("deep copy failed: " + err.Error())
	}
	return c
}
