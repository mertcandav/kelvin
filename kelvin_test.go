package kelvin

import "testing"

type test_structure struct {}

func TestNew(t *testing.T) {
	_ = New[test_structure]()

	defer func() {
		err := recover()
		if err == nil {
			t.Error("int is not a valid data type but kelvin is accept as valid")
		}
	}()

	_ = New[int]()
}
