package types

import (
	"testing"
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
