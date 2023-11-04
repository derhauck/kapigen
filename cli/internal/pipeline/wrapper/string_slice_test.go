package wrapper

import (
	"fmt"
	"testing"
)

func TestNewStringSlice(t *testing.T) {
	t.Parallel()
	t.Run("create new string slice", func(t *testing.T) {
		slice := NewStringSlice()
		if slice == nil {
			t.Error("should be able to create new slice")
		}
		if len(slice.Get()) != 0 {
			t.Error("should be empty")
		}
	})
	t.Run("mutate new string slice", func(t *testing.T) {
		slice := NewStringSlice()

		if slice.Has("test") {
			t.Error("should not have test")
		}

		slice.Add("test")
		if !slice.Has("test") {
			t.Error("should have test")
		}
		slice.AddSeveral("test2", "test3")
		if !slice.Has("test2") {
			t.Error("should have test2")
		}

		if !slice.Has("test3") {
			t.Error("should have test3")
		}

		if length := len(slice.Get()); length != 3 {
			t.Error(fmt.Sprintf("should have exactly 3 elements, received: %d", length))
		}

		slice.AddSlice([]string{"test4", "test5"})

		if !slice.Has("test4") {
			t.Error("should have test4")
		}

		if !slice.Has("test5") {
			t.Error("should have test5")
		}
		if length := len(slice.Get()); length != 5 {
			t.Error(fmt.Sprintf("should have exactly 5 elements, received: %d", length))
		}
	})
}
