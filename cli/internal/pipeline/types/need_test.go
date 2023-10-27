package types

import "testing"

func TestNeeds_RemoveNeed(t *testing.T) {
	t.Parallel()
	t.Run("Remove need", func(t *testing.T) {

		expectedNeed := &Need{
			Job:      nil,
			Optional: true,
		}
		anotherExpectedNeed := &Need{
			Job:      nil,
			Optional: false,
		}
		needs := Needs{
			&Need{
				Job:      nil,
				Optional: true,
			},
			expectedNeed,
			anotherExpectedNeed,
		}
		result := needs.RemoveNeed(expectedNeed)
		if len(needs.GetNeeds()) != 2 {
			t.Error("Should have 2 needs")
		}
		if result != true {
			t.Error("should have returned true for successful deletion")
		}

		result = needs.RemoveNeed(expectedNeed)

		if len(needs.GetNeeds()) != 2 {
			t.Error("Should have 2 needs")
		}
		if result != false {
			t.Error("should have returned false for failed deletion, need was already deleted")
		}
		result = needs.RemoveNeed(anotherExpectedNeed)
		if len(needs.GetNeeds()) != 1 {
			t.Error("Should have 1 needs")
		}
		if result != true {
			t.Error("should have returned true for successful deletion")
		}
	})
}
