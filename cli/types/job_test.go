package types

import (
	"reflect"
	"testing"

	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
)

func TestJobs_FindJobsByPipelineId(t *testing.T) {
	pipelineIdIFirst := "pipelineIdFirst"
	pipelineIdSecond := "pipelineIdSecond"
	pipelineJobFirstOne := &Job{PipelineId: pipelineIdIFirst, Names: []string{"pipelineJobFirstOne"}}
	pipelineJobFirstTwo := &Job{PipelineId: pipelineIdIFirst, Names: []string{"pipelineJobFirstTwo"}}
	pipelineJobFirstThree := &Job{PipelineId: pipelineIdIFirst, Names: []string{"pipelineJobFirstThree"}}
	pipelineJobSecondOne := &Job{PipelineId: pipelineIdSecond, Names: []string{"pipelineJobSecondOne"}}
	pipelineJobSecondTwo := &Job{PipelineId: pipelineIdSecond, Names: []string{"pipelineJobSecondTwo"}}
	pipelineJobs := Jobs{
		pipelineJobFirstOne,
		pipelineJobFirstTwo,
		pipelineJobFirstThree,
		pipelineJobSecondOne,
		pipelineJobSecondTwo,
	}
	conditions := []struct {
		name         string
		pipelineId   string
		expectedJobs Jobs
		result       bool
	}{
		{
			name:       "find all jobs in pipelineIdFirst",
			pipelineId: pipelineIdIFirst,
			expectedJobs: Jobs{
				pipelineJobFirstOne,
				pipelineJobFirstTwo,
				pipelineJobFirstThree,
			},
			result: true,
		},
		{
			name:       "find all jobs in pipelineIdFirst",
			pipelineId: pipelineIdSecond,
			expectedJobs: Jobs{
				pipelineJobSecondOne,
				pipelineJobSecondTwo,
			},
			result: true,
		},
		{
			name:       "missing all jobs in pipelineIdSecond",
			pipelineId: pipelineIdSecond,
			expectedJobs: Jobs{
				pipelineJobFirstOne,
				pipelineJobFirstTwo,
				pipelineJobFirstThree,
			},
			result: false,
		},
		{
			name:       "find no jobs in pipelineIdFirst with typo",
			pipelineId: "pipelineIdFirs",
			expectedJobs: Jobs{
				pipelineJobFirstOne,
				pipelineJobFirstTwo,
				pipelineJobFirstThree,
			},
			result: false,
		},
	}

	t.Run("Can find Jobs by id", func(t *testing.T) {
		for _, condition := range conditions {
			t.Run(condition.name, func(t *testing.T) {
				foundJobs, err := pipelineJobs.FindJobsByPipelineId(condition.pipelineId)
				if err != nil && condition.result {
					t.Error(err.Error())
				}
				if found := len(condition.expectedJobs.GetJobs()); found != len(foundJobs.GetJobs()) && condition.result {
					t.Errorf("Found only %d job(s), should be %d, id: %s", found, len(foundJobs.GetJobs()), condition.pipelineId)
				}
				for _, currentJob := range foundJobs.GetJobs() {
					hasFoundJob := false
					for _, expectedJob := range condition.expectedJobs.GetJobs() {
						if currentJob == expectedJob {
							hasFoundJob = true
						}
					}
					if condition.result == false && hasFoundJob {
						t.Errorf("Did find job: %v, but should not id: %s", currentJob.Names, condition.pipelineId)
					} else if condition.result && hasFoundJob == false {
						t.Errorf("Did not find job: %v, id: %s", currentJob.Names, condition.pipelineId)
					}
				}

			})

		}
	})
}

func TestNewJob(t *testing.T) {
	t.Parallel()

	t.Run("creates job with correct parameters", func(t *testing.T) {
		name := "Test Job"
		imageName := "golang:1.16"
		fn := func(ciJob *job.CiJob) {}

		newJob := NewJob(name, imageName, fn)

		if newJob.Names[0] != name {
			t.Errorf("unexpected job name: %s", newJob.Names[0])
		}

		if newJob.CiJob.Image.Name != imageName {
			t.Errorf("unexpected image name: %s", newJob.CiJob.Image.Name)
		}
	})
}

func TestJob_AddName(t *testing.T) {
	t.Parallel()

	t.Run("adds name to job", func(t *testing.T) {
		newJob := NewJob("Test Job", "golang:1.16", nil)
		newJob.AddName("Another Name")

		if len(newJob.Names) != 2 {
			t.Errorf("unexpected number of names: %d", len(newJob.Names))
		}

		if newJob.Names[1] != "Another Name" {
			t.Errorf("unexpected second name: %s", newJob.Names[1])
		}
	})
}

func TestJob_Render(t *testing.T) {
	t.Parallel()

	t.Run("renders job correctly", func(t *testing.T) {
		newJob := NewJob("Test Job", "golang:1.16", func(ciJob *job.CiJob) {
			ciJob.TagMediumPressure().
				Rules.AddRules(*job.DefaultMainBranchRules())
		})
		err := newJob.Render()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			t.FailNow()
		}

		if newJob.CiJobYaml == nil {
			t.Error("expected CiJobYaml to be set after rendering")
			t.FailNow()
		}
		if newJob.CiJobYaml.Image.Name != "golang:1.16" {
			t.Errorf("unexpected image name: %s", newJob.CiJobYaml.Image.Name)
		}

		if !contains(newJob.CiJobYaml.Tags, enum.TagPressureMedium.String()) {
			t.Errorf("expected to contain tag: %s", enum.TagPressureMedium.String())
		}

		if len(newJob.CiJobYaml.Tags) != 1 {
			t.Errorf("unexpected number of tags: %d", len(newJob.CiJobYaml.Tags))
		}
		if newJob.CiJobYaml.Stage != stages.Enum().ValueSafe(stages.DYNAMIC) {
			t.Errorf("unexpected stage: %s", newJob.CiJobYaml.Stage)
		}
	})
}
func TestJob_DynamicMerge(t *testing.T) {
	t.Parallel()

	t.Run("merges job with compatible job", func(t *testing.T) {
		job1 := NewJob("Job 1", "golang:1.16", func(ciJob *job.CiJob) {
			ciJob.Tags.Add(enum.TagPressureMedium)
			ciJob.TagMediumPressure().
				Rules.AddRules(*job.DefaultMainBranchRules())
		})
		job2 := NewJob("Job 2", "golang:1.16", func(ciJob *job.CiJob) {
			ciJob.Tags.Add(enum.TagPressureMedium)
			ciJob.TagMediumPressure().
				Rules.AddRules(*job.DefaultMainBranchRules())
		})
		jobs := Jobs{job1, job2}
		job2.AddJobAsNeed(job1)
		merged, err := job1.DynamicMerge(&jobs)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if merged != nil {
			t.Error("expected merged job to be nil")
		}

		if len(job1.Needs.GetNeeds()) != 0 {
			t.Errorf("unexpected number of needs for job1: %d", len(job1.Needs.GetNeeds()))
		}

		if len(job2.Needs.GetNeeds()) != 1 {
			t.Errorf("unexpected number of needs for job2: %d", len(job2.Needs.GetNeeds()))
		}

		if job2.Needs.HasJob(job1) {
			t.Error("job2 needs should include job1")
		}
	})

	t.Run("does not merge job with incompatible job", func(t *testing.T) {
		job1 := NewJob("Job 1", "golang:1.16", func(ciJob *job.CiJob) {
			ciJob.Tags.Add(enum.TagPressureMedium)
			ciJob.TagMediumPressure().
				Rules.AddRules(*job.DefaultMainBranchRules())
		})
		job2 := NewJob("Job 2", "node:14", func(ciJob *job.CiJob) {
			ciJob.Tags.Add(enum.TagPressureMedium)
			ciJob.TagMediumPressure().
				Rules.AddRules(*job.DefaultMainBranchRules())
		})
		jobs := Jobs{job1, job2}

		merged, err := job1.DynamicMerge(&jobs)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if merged != job1 {
			t.Error("expected merged job to be job1")
		}

		if len(job1.Needs.GetNeeds()) != 0 {
			t.Errorf("unexpected number of needs for job1: %d", len(job1.Needs.GetNeeds()))
		}

		if len(job2.Needs.GetNeeds()) != 0 {
			t.Errorf("unexpected number of needs for job2: %d", len(job2.Needs.GetNeeds()))
		}
	})
}

func TestJobs_DynamicMerge(t *testing.T) {
	t.Parallel()

	t.Run("merges compatible jobs", func(t *testing.T) {
		job1 := NewJob("Job 1", "golang:1.16", func(ciJob *job.CiJob) {
			ciJob.Tags.Add(enum.TagPressureMedium)
			ciJob.TagMediumPressure().
				Rules.AddRules(*job.DefaultMainBranchRules())
		})
		job2 := NewJob("Job 2", "golang:1.16", func(ciJob *job.CiJob) {
			ciJob.Tags.Add(enum.TagPressureMedium)
			ciJob.TagMediumPressure().
				Rules.AddRules(*job.DefaultMainBranchRules())
		})
		job3 := NewJob("Job 3", "node:14", func(ciJob *job.CiJob) {
			ciJob.Tags.Add(enum.TagPressureMedium)
			ciJob.TagMediumPressure().
				Rules.AddRules(*job.DefaultMainBranchRules())
		})
		jobs := Jobs{job1, job2, job3}
		job3.AddJobAsNeed(job1)
		merged, err := jobs.DynamicMerge()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			t.FailNow()
		}

		if len(merged.GetJobs()) != 2 {
			t.Errorf("unexpected number of merged jobs: %d", len(merged.GetJobs()))
		}

		mergedJob1 := merged.GetJobs()[0]
		mergedJob2 := merged.GetJobs()[1]

		if len(mergedJob1.Needs.GetNeeds()) != 0 {
			t.Errorf("unexpected number of needs for mergedJob1: %d", len(mergedJob1.Needs.GetNeeds()))
		}

		if len(mergedJob2.Needs.GetNeeds()) != 1 {
			t.Errorf("unexpected number of needs for mergedJob2: %d", len(mergedJob2.Needs.GetNeeds()))
		}

		if mergedJob1.Needs.HasJob(mergedJob2) {
			t.Error("mergedJob1 needs should include mergedJob2")
		}

		if merged.GetJobs()[1] != job3 {
			t.Error("job3 should not be merged")
		}
	})
}
func contains(slice []string, element string) bool {
	for _, s := range slice {
		if s == element {
			return true
		}
	}
	return false
}

func TestJob_EvaluateName(t *testing.T) {
	t.Run("do not create different names if they are already unique", func(t *testing.T) {
		job1 := NewJob("Job 1", "golang:1.16", nil)
		job2 := NewJob("Job 2", "golang:1.16", nil)
		job3 := NewJob("Job 3", "golang:1.16", nil)
		expected := &Jobs{job1, job2, job3}
		jobs := *expected

		result, err := jobs.EvaluateNames()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected result to be: %v, got: %v", expected, result)
		}
	})
	t.Run("create different names if they are not unique", func(t *testing.T) {
		expected := &Jobs{
			NewJob("Job 1", "golang:1.16", nil),
			NewJob("Job 1", "golang:1.16", nil),
			NewJob("Job 2", "golang:1.16", nil),
		}
		for _, currentJob := range expected.GetJobs() {
			currentJob.AddName("test")
			currentJob.AddName("test")
		}
		jobs := &Jobs{
			NewJob("Job 1", "golang:1.16", nil),
			NewJob("Job 1", "golang:1.16", nil),
			NewJob("Job 2", "golang:1.16", nil),
		}
		for _, currentJob := range jobs.GetJobs() {
			currentJob.AddName("test")
			currentJob.AddName("test")
		}
		result, err := jobs.EvaluateNames()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result.GetJobs()[0].GetName() == expected.GetJobs()[0].GetName() {
			t.Errorf("expected job name: '%s', received: '%s'", expected.GetJobs()[0].GetName(), result.GetJobs()[0].GetName())
		}

		if result.GetJobs()[1].GetName() != expected.GetJobs()[1].GetName() {
			t.Errorf("expected job name not to be identical received: '%s'", result.GetJobs()[1].GetName())
		}

		if result.GetJobs()[2].GetName() != expected.GetJobs()[2].GetName() {
			t.Errorf("expected job name not to be identical received: '%s'", result.GetJobs()[2].GetName())
		}
	})
}

func TestJobs_OverwriteTags(t *testing.T) {
	t.Run("overwrite tags", func(t *testing.T) {
		expectation := []string{"overwritten"}
		oldTags := []string{"tag1", "tag2"}
		jobs := Jobs{
			&Job{CiJobYaml: &job.CiJobYaml{Tags: oldTags}},
			&Job{CiJobYaml: &job.CiJobYaml{Tags: oldTags}},
			&Job{CiJobYaml: &job.CiJobYaml{Tags: oldTags}},
		}

		jobs.OverwriteTags(expectation)
		for _, finalJob := range jobs.GetJobs() {
			if !reflect.DeepEqual(finalJob.CiJobYaml.Tags, expectation) {
				t.Errorf("expected job to have tags: ['overwritten'], but it does not have %v", expectation)
			}
			if reflect.DeepEqual(finalJob.CiJobYaml.Tags, oldTags) {
				t.Error("expected job to not have old tags")
			}
		}
	})
}

func TestJob_RenderNeeds(t *testing.T) {
	t.Run("render needs", func(t *testing.T) {
		job1 := NewJob("Job 1", "golang:1.16", func(ciJob *job.CiJob) {
			ciJob.AddScript("hello world").
				SetStage(stages.TEST).
				TagMediumPressure().
				Rules.AddRules(*job.DefaultMainBranchRules())
		})

		result := job1.RenderNeeds()
		if result == nil {
			t.Error("expected rendered job, received nil")
			t.FailNow()
		}

		if result.CiJobYaml == nil {
			t.Error("expected CiJobYaml, received nil")
			t.FailNow()
		}
		if contains(result.CiJobYaml.Tags, stages.Enum().ValueSafe(stages.DYNAMIC)) {
			t.Error("expected stage to be set to dynamic")
		}
		if result.CiJobYaml.Needs.GetNeeds() != nil {
			t.Error("expected needs to be empty, received nil")
		}
		if len(result.CiJobYaml.Needs.GetNeeds()) != 0 {
			t.Errorf("expected needs to be empty, received %v", result.CiJobYaml.Needs.GetNeeds())
		}

	})
}
