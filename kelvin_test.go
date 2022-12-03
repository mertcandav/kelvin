package kelvin

import "testing"

type test_structure struct {}

func TestOpen(t *testing.T) {
	_, _ = Open[test_structure](NoWrite, InMemory)

	defer func() {
		err := recover()
		if err == nil {
			t.Error("int is not a valid data type but kelvin is accept as valid")
		}
	}()

	_, _ = Open[int](NoWrite, InMemory)
}

func TestOpenSafe(t *testing.T) {
	_, _ = OpenSafe[test_structure](NoWrite, InMemory, nil)

	defer func() {
		err := recover()
		if err == nil {
			t.Error("int is not a valid data type but kelvin is accept as valid")
		}
	}()

	_, _ = OpenSafe[int](NoWrite, InMemory, nil)
}
