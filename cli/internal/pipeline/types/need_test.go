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

func TestNeeds_HasJob(t *testing.T) {
	t.Parallel()
	t.Run("Can find Need", func(t *testing.T) {
		job := &Job{}
		job2 := &Job{}
		expectedNeed := &Need{
			Job:      job,
			Optional: true,
		}
		anotherExpectedNeed := &Need{
			Job:      job2,
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
		if !needs.HasJob(job) {
			t.Error("should be able to find job")
		}

		if !needs.HasJob(job2) {
			t.Error("should be able to find job2")
		}
		if !needs.HasJob(nil) {
			t.Error("should be able to find nil job")
		}

	})
}
