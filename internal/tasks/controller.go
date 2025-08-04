package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/internal/exporter"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/rest"
	"go.cluttr.dev/gitlab-exporter/internal/metaerr"
	"go.cluttr.dev/gitlab-exporter/internal/types"
)

const (
	maxRecordsPerPage int = 100
)

type ProjectsSettings struct {
	settings map[int64]ProjectSettings
	mu       sync.RWMutex
}

func (ps *ProjectsSettings) Len() int {
	return len(ps.settings)
}

func (ps *ProjectsSettings) Set(m map[int64]ProjectSettings) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.settings = m
}

func (ps *ProjectsSettings) Get(id int64) (ProjectSettings, bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	s, ok := ps.settings[id]
	return s, ok
}

func (ps *ProjectsSettings) Add(id int64, settings ProjectSettings) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.settings == nil {
		ps.settings = make(map[int64]ProjectSettings)
	}

	if _, ok := ps.settings[id]; ok { // already exists
		return false
	}
	ps.settings[id] = settings
	return true
}

func (ps *ProjectsSettings) Delete(id int64) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.settings == nil {
		return false
	}

	if _, ok := ps.settings[id]; !ok {
		return false
	}

	delete(ps.settings, id)
	return true
}

func (ps *ProjectsSettings) GetBatches(size int, filter func(ProjectSettings) bool) [][]ProjectSettings {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	settings := make([]ProjectSettings, 0, len(ps.settings))
	for _, v := range ps.settings {
		if filter != nil && !filter(v) {
			continue
		}
		settings = append(settings, v)
	}

	var batches [][]ProjectSettings
	for i := 0; i < len(settings); i += size {
		j := i + size
		if j > len(settings) {
			j = len(settings)
		}
		batches = append(batches, settings[i:j])
	}

	return batches
}

func (ps *ProjectsSettings) ExportDeployments(id int64) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	cfg, ok := ps.settings[id]
	if !ok {
		return false
	}

	if cfg.AccessLevels.Environment == ProjectAccessLevelDisabled {
		return false
	}

	return cfg.Export.Deployments.Enabled
}

func (ps *ProjectsSettings) ExportIssues(id int64) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	cfg, ok := ps.settings[id]
	if !ok {
		return false
	}

	if cfg.AccessLevels.Issues == ProjectAccessLevelDisabled {
		return false
	}

	return cfg.Export.Issues.Enabled
}

func (ps *ProjectsSettings) ExportTestReports(id int64) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	cfg, ok := ps.settings[id]
	if !ok {
		return false
	}
	return cfg.Export.TestReports.Enabled
}

func (ps *ProjectsSettings) ExportLogData(id int64) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	cfg, ok := ps.settings[id]
	if !ok {
		return false
	}
	return cfg.Export.Sections.Enabled || cfg.Export.Metrics.Enabled || cfg.Export.Jobs.Properties.Enabled
}

func (ps *ProjectsSettings) ExportTraces(id int64) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	cfg, ok := ps.settings[id]
	if !ok {
		return false
	}
	return cfg.Export.Traces.Enabled
}

func (ps *ProjectsSettings) ExportMergeRequests(id int64) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	cfg, ok := ps.settings[id]
	if !ok {
		return false
	}
	return cfg.Export.MergeRequests.Enabled
}

func (ps *ProjectsSettings) ExportJunitReports(id int64) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	cfg, ok := ps.settings[id]
	if !ok {
		return false
	}

	return cfg.Export.Reports.Junit.Enabled
}

type ProjectSettings struct {
	Id       int64
	FullPath string

	Export config.ProjectExport

	CatchUp struct {
		Enabled       bool
		UpdatedAfter  *time.Time
		UpdatedBefore *time.Time
	}

	AccessLevels ProjectAccessLevels
}

type ProjectAccessLevels struct {
	Builds      ProjectAccessLevel
	Environment ProjectAccessLevel
	Issues      ProjectAccessLevel
}

type ProjectAccessLevel string

const ( // https://gitlab.com/gitlab-org/gitlab/-/blob/master/app/models/concerns/featurable.rb
	ProjectAccessLevelUnspecified ProjectAccessLevel = ""
	ProjectAccessLevelDisabled    ProjectAccessLevel = "disabled"
	ProjectAccessLevelPrivate     ProjectAccessLevel = "private"
	ProjectAccessLevelEnabled     ProjectAccessLevel = "enabled"
	ProjectAccessLevelPublic      ProjectAccessLevel = "public"
)

type ControllerConfig struct {
	GitLab     config.GitLab
	Projects   []config.Project
	Namespaces []config.Namespace

	ResolveInterval time.Duration
	ExportInterval  time.Duration
	CatchUpInterval time.Duration
}

type Controller struct {
	GitLab   *gitlab.Client
	Exporter *exporter.Exporter

	config ControllerConfig

	projectsSettings      ProjectsSettings
	projectsSettingsMutex sync.RWMutex
}

func NewController(glab *gitlab.Client, exp *exporter.Exporter, cfg ControllerConfig) *Controller {
	if cfg.ResolveInterval == 0 {
		cfg.ResolveInterval = 30 * time.Minute
	}
	if cfg.ExportInterval == 0 {
		cfg.ExportInterval = 5 * time.Minute
	}
	if cfg.CatchUpInterval == 0 {
		cfg.CatchUpInterval = 24 * time.Hour
	}

	return &Controller{
		GitLab:   glab,
		Exporter: exp,

		config: cfg,
		projectsSettings: ProjectsSettings{
			settings: make(map[int64]ProjectSettings),
		},
	}
}

func (c *Controller) Run(ctx context.Context) error {
	var (
		resolveTicker *time.Ticker = time.NewTicker(c.config.ResolveInterval)
		exportTicker  *time.Ticker = time.NewTicker(c.config.ExportInterval)

		after  time.Time
		before time.Time = time.Now().UTC()
	)

	var iteration int
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-resolveTicker.C:
			slog.Debug("[RUN] Resolving projects...")
			n, err := c.ResolveProjects(ctx)
			if err != nil {
				slog.Error("[RUN] error resolving projects", "error", err)
				continue
			}
			slog.Debug("[RUN] Resolving projects... done", "found", n)
		case now := <-exportTicker.C:
			iteration++
			firstIteration := (iteration == 1)

			after = before
			before = now

			slog.Debug("[RUN] Processing projects...", "iteration", iteration, "after", after.Format(time.RFC3339), "before", before.Format(time.RFC3339))

			var wg sync.WaitGroup
			projectBatches := c.projectsSettings.GetBatches(maxRecordsPerPage, nil)
			for i, batch := range projectBatches {
				wg.Add(1)
				go func() {
					defer wg.Done()

					if err := c.process(ctx, batch, &after, &before, firstIteration); err != nil {
						slog.
							With(
								slog.String("error", err.Error()),
								slog.Int("iteration", iteration),
								slog.String("after", after.Format(time.RFC3339)), slog.String("before", before.Format(time.RFC3339)),
								slog.Int("batchIndex", i), slog.Int("batchCount", len(projectBatches)),
							).
							With(metaerr.GetMetadata(err)...).
							Error("[RUN] error processing projects")
					}
				}()
			}
			wg.Wait()

			slog.Debug("[RUN] Processing projects... done", "iteration", iteration)
		}
	}
}

func (c *Controller) CatchUp(ctx context.Context) error {
	var (
		interval time.Duration = 24 * time.Hour
		after    time.Time     = time.Now().UTC()
		before   time.Time
	)

	var iteration int
	for {
		iteration++
		firstIteration := (iteration == 1)

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			before = after
			after = before.Add(-interval)
		}

		projectBatches := c.projectsSettings.GetBatches(maxRecordsPerPage, func(ps ProjectSettings) bool {
			if !ps.CatchUp.Enabled {
				return false
			}
			if ps.CatchUp.UpdatedAfter != nil && ps.CatchUp.UpdatedAfter.After(before) {
				return false
			}
			if ps.CatchUp.UpdatedBefore != nil && ps.CatchUp.UpdatedBefore.Before(after) {
				return false
			}
			return true
		})
		if len(projectBatches) == 0 { // nothing to catch up on
			break
		}

		slog.Debug("[CATCHUP] Processing projects...", "iteration", iteration, "after", after.Format(time.RFC3339), "before", before.Format(time.RFC3339))

		var wg sync.WaitGroup
		for i, batch := range projectBatches {
			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := c.process(ctx, batch, &after, &before, firstIteration); err != nil {
					slog.
						With(
							slog.String("error", err.Error()),
							slog.Int("iteration", iteration),
							slog.String("after", after.Format(time.RFC3339)), slog.String("before", before.Format(time.RFC3339)),
							slog.Int("batchIndex", i), slog.Int("batchCount", len(projectBatches)),
						).
						With(metaerr.GetMetadata(err)...).
						Error("[CATCHUP] error processing projects")
				}
			}()
		}
		wg.Wait()

		slog.Debug("[CATCHUP] Processing projects... done", "iteration", iteration)
	}

	return nil
}

func (c *Controller) process(ctx context.Context, projectSettings []ProjectSettings, updatedAfter *time.Time, updatedBefore *time.Time, firstIteration bool) error {
	result, err := c.getUpdatedProjects(ctx, projectSettings, updatedAfter, updatedBefore, firstIteration)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := c.processProjects(ctx, result.UpdatedProjects); err != nil {
			errChan <- fmt.Errorf("process projects: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := c.processPipelines(ctx, result.ProjectsWithUpdatedPipelines, updatedAfter, updatedBefore); err != nil {
			errChan <- fmt.Errorf("process pipelines: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := c.processProjectMergeRequests(ctx, result.ProjectsWithUpdatedMergeRequests, updatedAfter, updatedBefore); err != nil {
			errChan <- fmt.Errorf("process merge requests: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		projectIds := make([]int64, 0, len(result.UpdatedProjects))
		for _, p := range result.UpdatedProjects {
			projectIds = append(projectIds, p.Id)
		}

		if err := c.processProjectsIssues(ctx, projectIds, updatedAfter, updatedBefore); err != nil {
			errChan <- fmt.Errorf("process issues: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := c.processProjectsDeployments(ctx, result.ProjectsWithUpdatedPipelines, updatedAfter, updatedBefore); err != nil {
			errChan <- fmt.Errorf("process deployments: %w", err)
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	var errs error
loop:
	for {
		select {
		case <-done:
			break loop
		case err := <-errChan:
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

type getUpdatedProjectsResult struct {
	UpdatedProjects                  []types.Project
	ProjectsWithUpdatedPipelines     []int64
	ProjectsWithUpdatedMergeRequests []int64
}

func (c *Controller) getUpdatedProjects(ctx context.Context, projects []ProjectSettings, updatedAfter *time.Time, updatedBefore *time.Time, firstIteration bool) (getUpdatedProjectsResult, error) {
	// Turn internal project ids into global ones.
	gids := make([]string, 0, len(projects))
	for _, project := range projects {
		gids = append(gids, graphql.FormatId(project.Id, graphql.GlobalIdProjectPrefix))
	}

	// Get project fields and number of updated pipelines and merge requests
	results, err := c.GitLab.GraphQL.GetProjects(ctx, gids, updatedAfter, updatedBefore)
	if err != nil {
		return getUpdatedProjectsResult{}, fmt.Errorf("get projects: %w", err)
	}

	result := getUpdatedProjectsResult{}
	for _, r := range results {
		project, err := graphql.ConvertProject(r.ProjectFields)
		if err != nil {
			return getUpdatedProjectsResult{}, fmt.Errorf("convert project fields: %w", err)
		}

		// Only export projects that were updated
		if isUpdated(project, updatedAfter, updatedBefore) || firstIteration {
			result.UpdatedProjects = append(result.UpdatedProjects, project)
		}

		if r.PipelinesCount > 0 {
			result.ProjectsWithUpdatedPipelines = append(result.ProjectsWithUpdatedPipelines, project.Id)
		}
		if r.MergeRequestsCount > 0 {
			result.ProjectsWithUpdatedMergeRequests = append(result.ProjectsWithUpdatedMergeRequests, project.Id)
		}
	}

	return result, nil
}

func (s *Controller) processProjects(ctx context.Context, projects []types.Project) error {
	return s.Exporter.ExportProjects(ctx, projects)
}

func (c *Controller) processProjectsIssues(ctx context.Context, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) error {
	pids := make([]int64, 0, len(projectIds))
	for _, pid := range projectIds {
		if c.projectsSettings.ExportIssues(pid) {
			pids = append(pids, pid)
		}
	}

	issues, err := FetchProjectsIssues(ctx, c.GitLab, pids, updatedAfter, updatedBefore)
	if errors.Is(err, context.Canceled) {
		return err
	} else if err != nil {
		err = fmt.Errorf("fetch issues: %w", err)
	}

	if err := c.Exporter.ExportIssues(ctx, issues); err != nil {
		return fmt.Errorf("export issues: %w", err)
	}

	return err
}

func (c *Controller) processProjectsDeployments(ctx context.Context, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) error {
	pids := make([]int64, 0, len(projectIds))
	for _, pid := range projectIds {
		if c.projectsSettings.ExportDeployments(pid) {
			pids = append(pids, pid)
		}
	}

	deployments, err := FetchProjectsDeployments(ctx, c.GitLab, pids, updatedAfter, updatedBefore)
	if err != nil {
		return fmt.Errorf("fetch deployments: %w", err)
	}

	if err := c.Exporter.ExportDeployments(ctx, deployments); err != nil {
		return fmt.Errorf("export deployments: %w", err)
	}

	return nil
}

func (c *Controller) handleError(errs *error, err error, msg string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		return err
	}

	*errs = errors.Join(*errs, fmt.Errorf("%s: %w", msg, err))
	return nil
}

func (c *Controller) processPipelines(ctx context.Context, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) error {
	var errs error

	pipelines, err := FetchProjectsPipelines(ctx, c.GitLab, projectIds, updatedAfter, updatedBefore)
	if err := c.handleError(&errs, err, "fetch projects pipelines"); err != nil {
		return err
	}

	err = c.Exporter.ExportPipelines(ctx, pipelines)
	if err := c.handleError(&errs, err, "export pipelines"); err != nil {
		return err
	}

	var tracePipelines []types.Pipeline
	for _, p := range pipelines {
		if c.projectsSettings.ExportTraces(p.Project.Id) {
			tracePipelines = append(tracePipelines, p)
		}
	}
	err = c.Exporter.ExportPipelineSpans(ctx, tracePipelines)
	if err := c.handleError(&errs, err, "export pipeline spans"); err != nil {
		return err
	}

	err = c.exportJobs(ctx, projectIds, updatedAfter, updatedBefore)
	if err := c.handleError(&errs, err, "export jobs"); err != nil {
		return err
	}

	err = c.exportReports(ctx, pipelines)
	if err := c.handleError(&errs, err, "export reports"); err != nil {
		return err
	}

	return errs
}

func (c *Controller) exportJobs(ctx context.Context, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) error {
	var errs error

	jobs, err := FetchProjectsPipelinesJobs(ctx, c.GitLab, projectIds, updatedAfter, updatedBefore)
	if err := c.handleError(&errs, err, "fetch projects pipelines jobs"); err != nil {
		return err
	}

	var logDataProjectJobs []types.Job
	for _, job := range jobs {
		if job.Kind == types.JobKindBridge { // bridges don't have logs
			continue
		}
		if !c.projectsSettings.ExportLogData(job.Pipeline.Project.Id) {
			continue
		}

		logDataProjectJobs = append(logDataProjectJobs, job)
	}

	sections, metrics, properties, err := FetchProjectsJobsLogData(ctx, c.GitLab, logDataProjectJobs)
	if err := c.handleError(&errs, err, "fetch projects job log data"); err != nil {
		return err
	}
	for jobId, props := range properties {
		for i := 0; i < len(jobs); i++ {
			if jobs[i].Id == jobId {
				jobs[i].Properties = props
				break // inner loop
			}
		}
	}

	err = c.Exporter.ExportJobs(ctx, jobs)
	if err := c.handleError(&errs, err, "export jobs"); err != nil {
		return err
	}
	err = c.Exporter.ExportSections(ctx, sections)
	if err := c.handleError(&errs, err, "export sections"); err != nil {
		return err
	}

	var traceJobs []types.Job
	var traceSections []types.Section
	for _, j := range jobs {
		if c.projectsSettings.ExportTraces(j.Pipeline.Project.Id) {
			traceJobs = append(traceJobs, j)

			for _, s := range sections {
				if s.Job.Id == j.Id {
					traceSections = append(traceSections, s)
				}
			}
		}
	}
	err = c.Exporter.ExportJobSpans(ctx, traceJobs)
	if err := c.handleError(&errs, err, "export job spans"); err != nil {
		return err
	}
	err = c.Exporter.ExportSectionSpans(ctx, traceSections)
	if err := c.handleError(&errs, err, "export section spans"); err != nil {
		return err
	}

	err = c.Exporter.ExportMetrics(ctx, metrics)
	if err := c.handleError(&errs, err, "export metrics"); err != nil {
		return err
	}

	return errs
}

func (c *Controller) exportReports(ctx context.Context, pipelines []types.Pipeline) error {
	testReportProjectPipelines := []types.Pipeline{}
	junitReportProjectPipelines := make(map[string][]string)
	junitReportProjectArtifactPaths := make(map[string][]string)
	coberturaReportProjectPipelines := make(map[string][]string)
	coberturaReportProjectArtifactPaths := make(map[string][]string)
	for _, p := range pipelines {
		settings, ok := c.projectsSettings.Get(p.Project.Id)
		if !ok {
			continue
		}

		if settings.Export.Reports.Enabled {
			projectPath := p.Project.FullPath
			pipelineIid := strconv.FormatInt(p.Iid, 10)
			if settings.Export.Reports.Junit.Enabled {
				junitReportProjectPipelines[projectPath] = append(junitReportProjectPipelines[projectPath], pipelineIid)
				junitReportProjectArtifactPaths[projectPath] = settings.Export.Reports.Junit.Paths
			}
			if settings.Export.Reports.Coverage.Enabled {
				coberturaReportProjectPipelines[projectPath] = append(coberturaReportProjectPipelines[projectPath], pipelineIid)
				coberturaReportProjectArtifactPaths[projectPath] = settings.Export.Reports.Coverage.Paths
			}
		}

		if (!settings.Export.Reports.Enabled || !settings.Export.Reports.Junit.Enabled) && c.projectsSettings.ExportTestReports(p.Project.Id) {
			testReportProjectPipelines = append(testReportProjectPipelines, p)
		}
	}

	var joinedErr error

	// fetch test reports
	testReports, testSuites, testCases, err := FetchProjectsPipelinesTestReports(ctx, c.GitLab, testReportProjectPipelines)
	if err := c.handleError(&joinedErr, err, "fetch test reports"); err != nil {
		return err
	}
	junitReports, junitSuites, junitCases, err := FetchProjectsPipelinesJunitReports(ctx, c.GitLab, junitReportProjectPipelines, junitReportProjectArtifactPaths)
	if err := c.handleError(&joinedErr, err, "fetch junit reports"); err != nil {
		return err
	}
	// export test reports
	err = c.Exporter.ExportTestReports(ctx, append(junitReports, testReports...))
	if herr := c.handleError(&joinedErr, err, "test reports"); herr != nil {
		return herr
	}
	err = c.Exporter.ExportTestSuites(ctx, append(junitSuites, testSuites...))
	if herr := c.handleError(&joinedErr, err, "test suites"); herr != nil {
		return herr
	}
	err = c.Exporter.ExportTestCases(ctx, append(junitCases, testCases...))
	if herr := c.handleError(&joinedErr, err, "test cases"); herr != nil {
		return herr
	}

	// fetch coverage reports
	covReports, covPackages, covClasses, covMethods, err := FetchProjectsPipelinesCoberturaReports(ctx, c.GitLab, coberturaReportProjectPipelines, coberturaReportProjectArtifactPaths)
	if err := c.handleError(&joinedErr, err, "fetch coverage reports"); err != nil {
		return err
	}
	// export coverage reports
	err = c.Exporter.ExportCoverageReports(ctx, covReports)
	if herr := c.handleError(&joinedErr, err, "coverage reports"); herr != nil {
		return herr
	}
	err = c.Exporter.ExportCoveragePackages(ctx, covPackages)
	if herr := c.handleError(&joinedErr, err, "coverage packages"); herr != nil {
		return herr
	}
	err = c.Exporter.ExportCoverageClasses(ctx, covClasses)
	if herr := c.handleError(&joinedErr, err, "coverage packages"); herr != nil {
		return herr
	}
	err = c.Exporter.ExportCoverageMethods(ctx, covMethods)
	if herr := c.handleError(&joinedErr, err, "coverage reports"); herr != nil {
		return herr
	}

	return joinedErr
}

func (c *Controller) processProjectMergeRequests(ctx context.Context, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) error {
	var errs []error

	mergeRequests, err := FetchProjectsMergeRequests(ctx, c.GitLab, projectIds, updatedAfter, updatedBefore)
	if err != nil {
		errs = append(errs, fmt.Errorf("fetch merge requests: %w", err))
	}

	mergeRequestNoteEvents, err := FetchProjectsMergeRequestsNotes(ctx, c.GitLab, projectIds, updatedAfter, updatedBefore)
	if err != nil {
		errs = append(errs, fmt.Errorf("fetch merge request note events: %w", err))
	}

	if err := c.Exporter.ExportMergeRequests(ctx, mergeRequests); err != nil {
		errs = append(errs, fmt.Errorf("export merge requests: %w", err))
	}

	if err := c.Exporter.ExportMergeRequestNoteEvents(ctx, mergeRequestNoteEvents); err != nil {
		errs = append(errs, fmt.Errorf("export merge request note events: %w", err))
	}

	return errors.Join(errs...)
}

func (c *Controller) ResolveProjects(ctx context.Context) (int, error) {
	c.projectsSettingsMutex.Lock()
	defer c.projectsSettingsMutex.Unlock()

	projectsSettings := make(map[int64]ProjectSettings)

	opt := rest.ListNamespaceProjectsOptions{}
	for _, namespace := range c.config.Namespaces {
		opt.Kind = namespace.Kind
		if namespace.Visibility != "" {
			opt.Visibility = &namespace.Visibility
		}

		opt.WithShared = gitlab.Ptr(namespace.WithShared)
		opt.IncludeSubgroups = gitlab.Ptr(namespace.IncludeSubgroups)

		var updatedAfter *time.Time
		if namespace.ProjectSettings.CatchUp.UpdatedAfter != "" {
			after, err := parseTimeISO8601(namespace.ProjectSettings.CatchUp.UpdatedAfter)
			if err != nil {
				return 0, fmt.Errorf("error parsing catchup update_after: %w", err)
			}
			updatedAfter = &after
		}
		var updatedBefore *time.Time
		if namespace.ProjectSettings.CatchUp.UpdatedBefore != "" {
			before, err := parseTimeISO8601(namespace.ProjectSettings.CatchUp.UpdatedBefore)
			if err != nil {
				return 0, fmt.Errorf("error parsing catchup update_before: %w", err)
			}
			updatedBefore = &before
		}

		err := c.GitLab.Rest.ListNamespaceProjects(ctx, namespace.Id, opt, func(projects []*rest.Project) bool {
			for _, project := range projects {
				ps := ProjectSettings{
					Id:       int64(project.ID),
					FullPath: project.PathWithNamespace,

					Export: namespace.ProjectSettings.Export,

					AccessLevels: ProjectAccessLevels{
						Builds:      ProjectAccessLevel(project.BuildsAccessLevel),
						Environment: ProjectAccessLevel(project.EnvironmentsAccessLevel),
						Issues:      ProjectAccessLevel(project.IssuesAccessLevel),
					},
				}
				ps.CatchUp.Enabled = namespace.ProjectSettings.CatchUp.Enabled
				if updatedAfter != nil {
					ps.CatchUp.UpdatedAfter = updatedAfter
				} else {
					ps.CatchUp.UpdatedAfter = project.CreatedAt
				}
				if updatedBefore != nil {
					ps.CatchUp.UpdatedBefore = updatedBefore
				}

				projectsSettings[int64(project.ID)] = ps
			}

			for _, pid := range namespace.ExcludeProjects {
				p, err := c.GitLab.Rest.GetProject(ctx, pid)
				if err != nil {
					slog.Error("error getting namespace project", "namespace_id", namespace.Id, "project", pid, "error", err)
					return false
				}
				delete(projectsSettings, int64(p.ID))
			}

			return true
		})
		if err != nil {
			return 0, err
		}
	}

	// overwrite with explicitly configured projects
	for _, p := range c.config.Projects {
		project, err := c.GitLab.Rest.GetProject(ctx, int(p.Id))
		if err != nil {
			return 0, err
		}

		var updatedAfter *time.Time = project.CreatedAt
		if p.ProjectSettings.CatchUp.UpdatedAfter != "" {
			after, err := parseTimeISO8601(p.ProjectSettings.CatchUp.UpdatedAfter)
			if err != nil {
				return 0, fmt.Errorf("error parsing catchup update_after: %w", err)
			}
			updatedAfter = &after
		}
		var updatedBefore *time.Time
		if p.ProjectSettings.CatchUp.UpdatedBefore != "" {
			before, err := parseTimeISO8601(p.ProjectSettings.CatchUp.UpdatedBefore)
			if err != nil {
				return 0, fmt.Errorf("error parsing catchup update_before: %w", err)
			}
			updatedBefore = &before
		}

		ps := ProjectSettings{
			Id:       int64(project.ID),
			FullPath: project.PathWithNamespace,

			Export: p.ProjectSettings.Export,

			AccessLevels: ProjectAccessLevels{
				Builds:      ProjectAccessLevel(project.BuildsAccessLevel),
				Environment: ProjectAccessLevel(project.EnvironmentsAccessLevel),
			},
		}

		ps.CatchUp.Enabled = p.ProjectSettings.CatchUp.Enabled
		ps.CatchUp.UpdatedAfter = updatedAfter
		ps.CatchUp.UpdatedBefore = updatedBefore

		projectsSettings[int64(project.ID)] = ps
	}

	c.projectsSettings.Set(projectsSettings)
	return c.projectsSettings.Len(), nil
}

func isUpdated(p types.Project, after *time.Time, before *time.Time) bool {
	var updated bool = true

	if p.UpdatedAt != nil {
		if after != nil && p.UpdatedAt.Before(*after) {
			updated = false
		} else if before != nil && p.UpdatedAt.After(*before) {
			updated = false
		}
	}
	if p.LastActivityAt != nil {
		if after != nil && p.LastActivityAt.Before(*after) {
			updated = false
		} else if before != nil && p.LastActivityAt.After(*before) {
			updated = false
		}
	}

	return updated
}

func parseTimeISO8601(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05+07:00",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.In(time.UTC), nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse time from %q", s)
}
