package pipeline

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	gitlab2 "github.com/xanzy/go-gitlab"
	"gitlab.com/kateops/kapigen/cli/factory"
	"gitlab.com/kateops/kapigen/cli/internal/cli"
	"gitlab.com/kateops/kapigen/cli/internal/pipeline/wrapper"
	"gitlab.com/kateops/kapigen/cli/internal/version"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/environment"
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gopkg.in/yaml.v3"
)

var ReportsCmd = &cobra.Command{
	Use:              "reports",
	Short:            "Get pipeline",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.PreparePersistentFlags(cmd)
		logger.Debug("activated verbose mode")

		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		privateTokenName, err := cmd.Flags().GetString("private-token")
		if err != nil {
			return err
		}
		logger.Info("will create settings")

		logger.Info("will try to read pipeline config from: " + configPath)
		cmd.SilenceUsage = true
		if body, err := os.ReadFile(configPath); err != nil {
			logger.ErrorE(err)
		} else {
			var pipelineConfig types.PipelineConfig
			err = yaml.Unmarshal(body, &pipelineConfig)
			if err != nil {
				logger.ErrorE(err)
			}
			if privateTokenName == "" && err == nil {
				privateTokenName = pipelineConfig.PrivateTokenName
			}
		}

		logger.Debug("will use private token: " + privateTokenName)
		settings := cli.NewSettings(
			cli.SetMode(version.GetModeFromString(version.Gitlab.Name())),
			cli.SetPrivateToken(privateTokenName),
		)
		gitlab := factory.New(settings).GetGitlabClient()
		pipelineId, err := strconv.ParseInt(environment.CI_PIPELINE_ID.Get(), 10, 32)
		if err != nil {
			return err
		}
		bridges, res, err := gitlab.Jobs.ListPipelineBridges(environment.CI_PROJECT_ID.Get(), int(pipelineId), nil)
		if res.StatusCode != 200 {
			logger.Error(res.Status)
			return err
		}
		var downstreamPipelineIds []int
		var coverageValues []float64
		for _, bridge := range bridges {
			if bridge.DownstreamPipeline != nil {
				downstreamPipelineIds = append(downstreamPipelineIds, bridge.DownstreamPipeline.ID)
				logger.Debug(fmt.Sprintf("found trigger job: %s", bridge.Name))
			}
		}
		var reportJobs wrapper.Array[gitlab2.Job]
		for _, downstreamPipelineId := range downstreamPipelineIds {
			jobs, res, err := gitlab.Jobs.ListPipelineJobs(environment.CI_PROJECT_ID.Get(), downstreamPipelineId, nil)
			if res.StatusCode != 200 {
				logger.Error(res.Status)
				return err
			}
			for _, job := range jobs {
				if job.Artifacts != nil {
					for _, artifact := range job.Artifacts {
						if artifact.FileType == "junit" {
							reportJobs.Push(*job)
							logger.Info(fmt.Sprintf("found reports in job: %s", job.Name))
						}
					}
				}
				if job.Coverage > 0 {
					coverageValues = append(coverageValues, job.Coverage)
				}
			}
		}
		if reportJobs.Length() == 0 {
			logger.Error("no reports found")
			return nil
		}
		var totalCoverage float64 = 0
		for _, coverageValue := range coverageValues {
			totalCoverage += coverageValue
		}
		if totalCoverage > 0 {
			logger.Info(fmt.Sprintf("coverage: %f", totalCoverage/float64(len(coverageValues))))
		}
		_ = os.Mkdir(".reports", 0750)
		logger.Info("unzipping archives")
		reportJobs.ForEach(func(e *gitlab2.Job) {
			jobArtifact, res, err := gitlab.Jobs.GetJobArtifacts(environment.CI_PROJECT_ID.Get(), e.ID, nil)
			if err != nil {
				logger.Error(e.Name)
				logger.ErrorE(err)
				return
			}
			artifactsDir := fmt.Sprintf(".reports/%d", e.ID)
			_ = os.Mkdir(artifactsDir, 0750)
			artifactFile := fmt.Sprintf("%s/reports.zip", artifactsDir)
			if res.StatusCode != 200 {
				logger.Error(res.Status)
				logger.Error(e.Name)
				return
			}
			fi, err := os.Create(artifactFile)
			if err != nil {
				logger.Error(e.Name)
				logger.ErrorE(err)
				return
			}
			defer func() {
				_ = fi.Close()
			}()
			buffer := make([]byte, 1024)
			for {
				n, err := jobArtifact.Read(buffer)
				if err != nil && err != io.EOF {
					logger.Error(e.Name)
					logger.Error(res.Status)
					logger.ErrorE(err)
					return
				}
				if _, err = fi.Write(buffer[:n]); err != nil {
					logger.Error(e.Name)
					logger.ErrorE(err)
					return
				}

				if n == 0 {
					break
				}
			}

			archive, err := zip.OpenReader(artifactFile)
			if err != nil {
				panic(err)
			}
			defer func(archive *zip.ReadCloser) {
				_ = archive.Close()
			}(archive)

			for _, f := range archive.File {
				filePath := filepath.Join(artifactsDir, f.Name)

				if !strings.HasPrefix(filePath, filepath.Clean(artifactsDir)+string(os.PathSeparator)) {
					logger.Error(fmt.Sprintf("invalid file path: %s", filePath))
					return
				}
				if f.FileInfo().IsDir() {
					logger.Debug("creating directory: %s", filePath)
					_ = os.MkdirAll(filePath, os.ModePerm)
					continue
				}

				if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
					logger.Error(e.Name)
					logger.ErrorE(err)
				}

				dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				logger.Debug(fmt.Sprintf("writing file: %s", filePath))
				if err != nil {
					logger.Error(e.Name)
					logger.ErrorE(err)
				}

				fileInArchive, err := f.Open()
				if err != nil {
					logger.Error(e.Name)
					logger.ErrorE(err)
				}

				if _, err := io.Copy(dstFile, fileInArchive); err != nil {
					logger.Error(e.Name)
					logger.ErrorE(err)
				}

				_ = dstFile.Close()
				_ = fileInArchive.Close()
			}

		})
		return err
	},
}

func init() {
	ReportsCmd.Flags().String("config", "config.kapigen.yaml", "config to use")
}
