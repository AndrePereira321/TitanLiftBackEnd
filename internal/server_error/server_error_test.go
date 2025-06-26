package server_error

import (
	"errors"
	"testing"
)

func TestServerErrorComparing(t *testing.T) {
	t.Run("comparing server errors", func(t *testing.T) {
		error1 := New("TEST", "Test Error")
		error1Clone := New("TEST", "Test Error")
		error2 := New("TEST2", "Test Error 2")
		libError := errors.New("test Error")

		if !errors.Is(error1, error1Clone) {
			t.Error("Failing comparing clone errors")
		}

		if errors.Is(error1, error2) {
			t.Error("Failing comparing different errors")
		}

		if errors.Is(error1, libError) {
			t.Error("Failing comparing lib errors")
		}

	})

	t.Run("comparing server wrapped errors with cause", func(t *testing.T) {
		cause1 := errors.New("test Error")
		cause2 := errors.New("test Error 2")

		error1 := Wrap("TEST", "Test Error", cause1)
		error1Clone := Wrap("TEST", "Test Error", cause1)
		error2 := Wrap("TEST", "Test Error", cause2)
		error3Unwraped := New("TEST", "Test Error 3")

		if !errors.Is(error1, error1Clone) {
			t.Error("Failing comparing clone wrapped errors")
		}

		if errors.Is(error1, error2) {
			t.Error("Failing comparing different wrapped errors")
		}

		if errors.Is(error1, error3Unwraped) {
			t.Error("Failing comparing different unwrapped errors")
		}
	})

}
