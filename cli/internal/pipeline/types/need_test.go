package types

import (
	"testing"
)

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

func TestNeeds_HasNeed(t *testing.T) {
	t.Run("has need", func(t *testing.T) {
		jobOne := Job{}
		jobTwo := Job{}
		jobThree := Job{}

		needOne := Need{Job: &jobOne}
		needTwo := Need{Job: &jobTwo}
		needThree := Need{Job: &jobThree}
		needs := Needs{
			&needOne,
			&needTwo,
		}

		if !needs.HasNeed(&needOne) {
			t.Error("list 'needs' should have 'needOne'")
		}
		if !needs.HasNeed(&needTwo) {
			t.Error("list 'needs' should have 'needTwo'")
		}
		if needs.HasNeed(&needThree) {
			t.Error("list 'needs' should not have 'needThree'")
		}

	})
}

func TestNeed_NotOptional(t *testing.T) {
	t.Run("not optional", func(t *testing.T) {
		need := Need{
			Job:      &Job{},
			Optional: true,
		}
		need.NotOptional()
		if need.Optional {
			t.Error("should not be optional")
		}
	})
	t.Run("stays not optional", func(t *testing.T) {
		need := Need{
			Job:      &Job{},
			Optional: false,
		}
		if need.Optional {
			t.Error("should not be optional")
		}
	})
}
