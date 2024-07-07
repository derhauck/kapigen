package wrapper

import (
	"reflect"
	"testing"
)

func TestNewArray(t *testing.T) {
	t.Run("create new array", func(t *testing.T) {
		array := NewArray[string]()
		if array == nil {
			t.Error("should be able to create new array")
		}
		if array == nil && array.slice != nil && len(array.slice) != 0 {
			t.Error("should be empty")
		}
	})
}
func TestArray_Push(t *testing.T) {
	t.Run("use array push", func(t *testing.T) {
		array := NewArray[string]()
		if len(array.slice) != 0 {
			t.Error("should be empty")
		}

		array.Push("test")
		array.Push("test2")
		array.Push("test3")
		array.Push("test4")
		array.Push("test5")

		if len(array.slice) != 5 {
			t.Errorf("should have exactly 5 elements, received: %d", len(array.slice))
		}
	})
}

func TestArray_Find(t *testing.T) {
	t.Parallel()
	t.Run("can not find element", func(t *testing.T) {
		array := NewArray[string]()
		array.slice = []string{"test", "test2", "test3", "test4", "test5"}

		element, index := array.Find(func(element *string) bool {
			return *element == "test2222"
		})

		if element != nil {
			t.Errorf("should not have found element, received: %s", *element)
		}

		if index != -1 {
			t.Errorf("should not have found element and returned index -1, received: %d", index)
		}

	})
	t.Run("can find element", func(t *testing.T) {
		array := NewArray[string]()
		array.slice = []string{"test", "test2", "test3", "test4", "test5"}

		element, index := array.Find(func(element *string) bool {
			return *element == "test2"
		})

		if element == nil {
			t.Error("should have found element")
		}

		if index != 1 {
			t.Errorf("should have found element at index 1, received: %d", index)
		}

		if element != nil && *element != "test2" {
			t.Errorf("should have found element test2, received: %s", *element)
		}

	})
}

func TestArray_Length(t *testing.T) {
	t.Run("get length", func(t *testing.T) {
		array := NewArray[string]()
		array.slice = []string{"test", "test2", "test3", "test4", "test5"}

		if array.Length() != len(array.slice) {
			t.Errorf("should have same length, received: %d", array.Length())
		}
	})
}

func TestArray_ForEach(t *testing.T) {
	array := NewArray[string]()
	array.slice = []string{"test", "test2", "test3", "test4", "test5"}

	result := []string{}
	array.ForEach(func(element *string) {
		result = append(result, *element)
	})

	if len(result) != len(array.slice) {
		t.Errorf("should have same length, received: %d", len(result))
	}
	if reflect.DeepEqual(result, array.slice) == false {
		t.Error("should have same elements")
	}

}
func TestArray_Map(t *testing.T) {
	array := NewArray[string]()
	array.slice = []string{"test", "test2", "test3", "test4", "test5"}

	expectation := []string{"test6", "test7", "test8", "test9", "test10"}
	counter := -1
	result := array.Map(func(element *string) string {
		counter++
		return expectation[counter]
	})

	if len(expectation) != len(array.slice) {
		t.Errorf("should have same length, received: %d", len(expectation))
	}
	if reflect.DeepEqual(expectation, array.slice) == true {
		t.Error("should not have same elements as initial array")
	}

	if reflect.DeepEqual(result.slice, expectation) == false {
		t.Errorf("should have same elements %v, received: %v", expectation, result.slice)
	}

}

func TestArray_Has(t *testing.T) {
	t.Run("has element", func(t *testing.T) {
		array := NewArray[string]()
		array.slice = []string{"test", "test2", "test3", "test4", "test5"}

		if array.Has("test2222") {
			t.Error("should not have found element")
		}
		if !array.Has("test2") {
			t.Error("should have found element")
		}
	})

}
