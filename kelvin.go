package kelvin

import "reflect"

// kelvin is kelvin database structure.
type kelvin[T any] struct {
	stack []T
}

// New returns new instance of Kelvin by data type.
// Panics if T is not structure.
func New[T any]() *kelvin[T] {
	var t T
	kind := reflect.TypeOf(t).Kind()
	if kind != reflect.Struct {
		panic("type error: kelvin supports only structures")
	}

	k := new(kelvin[T])
	k.stack = nil
	return k
}
