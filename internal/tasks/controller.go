package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/graphql"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/rest"
	"github.com/cluttrdev/gitlab-exporter/internal/types"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
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
	return cfg.Export.Sections.Enabled || cfg.Export.Metrics.Enabled
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
}

type ProjectAccessLevel string

const ( // https://gitlab.com/gitlab-org/gitlab/-/blob/master/app/models/concerns/featurable.rb
	ProjectAccessLevelUnspecified ProjectAccessLevel = ""
	ProjectAccessLevelDisabled                       = "disabled"
	ProjectAccessLevelPrivate                        = "private"
	ProjectAccessLevelEnabled                        = "enabled"
	ProjectAccessLevelPublic                         = "public"
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
						slog.Error("[RUN] error processing projects", "error", err,
							"iteration", iteration, "after", after.Format(time.RFC3339), "before", before.Format(time.RFC3339),
							"batchIndex", i, "batchCount", len(projectBatches),
						)
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
					slog.Error("[CATCHUP] error processing projects", "error", err,
						"iteration", iteration, "after", after.Format(time.RFC3339), "before", before.Format(time.RFC3339),
						"batchIndex", i, "batchCount", len(projectBatches),
					)
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

		if err := c.exportProjects(ctx, result.UpdatedProjects); err != nil {
			errChan <- fmt.Errorf("export projects: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := c.processProjectsCiData(ctx, result.ProjectsWithUpdatedPipelines, updatedAfter, updatedBefore); err != nil {
			errChan <- fmt.Errorf("process ci data: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := c.processProjectsMrData(ctx, result.ProjectsWithUpdatedMergeRequests, updatedAfter, updatedBefore); err != nil {
			errChan <- fmt.Errorf("process mr data: %w", err)
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

func (s *Controller) exportProjects(ctx context.Context, projects []types.Project) error {
	ps := make([]*typespb.Project, 0, len(projects))
	for _, p := range projects {
		ps = append(ps, types.ConvertProject(p))
	}

	return s.Exporter.ExportProjects(ctx, ps)
}

type projectsCiData struct {
	Pipelines        []types.Pipeline
	Jobs             []types.Job
	Sections         []types.Section
	Metrics          []types.Metric
	Traces           []*typespb.Trace
	TestReports      []types.TestReport
	TestSuites       []types.TestSuite
	TestCases        []types.TestCase
	CoverageReports  []types.CoverageReport
	CoveragePackages []types.CoveragePackage
	CoverageClasses  []types.CoverageClass
	CoverageMethods  []types.CoverageMethod
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

	pbDeployments := make([]*typespb.Deployment, 0, len(deployments))
	for _, deployment := range deployments {
		pbDeployments = append(pbDeployments, types.ConvertDeployment(deployment))
	}

	if err := c.Exporter.ExportDeployments(ctx, pbDeployments); err != nil {
		return fmt.Errorf("export deployments: %w", err)
	}

	return nil
}

func (c *Controller) processProjectsCiData(ctx context.Context, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) error {
	data, err := c.fetchProjectsCiData(ctx, projectIds, updatedAfter, updatedBefore)
	if err != nil {
		return fmt.Errorf("fetch projects ci data: %w", err)
	}

	if err := c.exportProjectsCiData(ctx, data); err != nil {
		return fmt.Errorf("export projects ci data: %w", err)
	}

	return nil
}

func (c *Controller) fetchProjectsCiData(ctx context.Context, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) (projectsCiData, error) {
	var errs error

	handleError := func(err error, msg string) error {
		if err == nil {
			return nil
		}

		if errors.Is(err, context.Canceled) {
			return err
		}

		errs = errors.Join(errs, fmt.Errorf("%s: %w", msg, err))
		return nil
	}

	pipelines, err := FetchProjectsPipelines(ctx, c.GitLab, projectIds, updatedAfter, updatedBefore)
	if err := handleError(err, "fetch projects pipelines"); err != nil {
		return projectsCiData{}, err
	}

	jobs, err := FetchProjectsPipelinesJobs(ctx, c.GitLab, projectIds, updatedAfter, updatedBefore)
	if err := handleError(err, "fetch projects pipelines jobs"); err != nil {
		return projectsCiData{}, err
	}

	var logDataProjectJobs []types.Job
	for _, job := range jobs {
		if job.Kind == types.JobKindBridge { // bridges don't have logs
			continue
		}
		if !c.projectsSettings.ExportLogData(job.Pipeline.Project.Id) { // duh
			continue
		}

		logDataProjectJobs = append(logDataProjectJobs, job)
	}
	sections, metrics, err := FetchProjectsJobsLogData(ctx, c.GitLab, logDataProjectJobs)
	if err := handleError(err, "fetch job log data"); err != nil {
		return projectsCiData{}, err
	}

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
	testReports, testSuites, testCases, err := FetchProjectsPipelinesTestReports(ctx, c.GitLab, testReportProjectPipelines)
	if err := handleError(err, "fetch test reports"); err != nil {
		return projectsCiData{}, err
	}
	junitReports, junitSuites, junitCases, err := FetchProjectsPipelinesJunitReports(ctx, c.GitLab, junitReportProjectPipelines, junitReportProjectArtifactPaths)
	if err := handleError(err, "fetch junit reports"); err != nil {
		return projectsCiData{}, err
	}
	covReports, covPackages, covClasses, covMethods, err := FetchProjectsPipelinesCoberturaReports(ctx, c.GitLab, coberturaReportProjectPipelines, coberturaReportProjectArtifactPaths)
	if err := handleError(err, "fetch coverage reports"); err != nil {
		return projectsCiData{}, err
	}

	traceData, err := c.convertTraceData(pipelines, jobs, sections)
	if err := handleError(err, "convert trace spans"); err != nil {
		return projectsCiData{}, err
	}
	var traces []*typespb.Trace
	if len(traceData.ResourceSpans) > 0 {
		traces = append(traces, &typespb.Trace{
			Data: traceData,
		})
	}

	return projectsCiData{
		Pipelines:        pipelines,
		Jobs:             jobs,
		Sections:         sections,
		Metrics:          metrics,
		Traces:           traces,
		TestReports:      append(testReports, junitReports...),
		TestSuites:       append(testSuites, junitSuites...),
		TestCases:        append(testCases, junitCases...),
		CoverageReports:  covReports,
		CoveragePackages: covPackages,
		CoverageClasses:  covClasses,
		CoverageMethods:  covMethods,
	}, nil
}

func (c *Controller) exportProjectsCiData(ctx context.Context, data projectsCiData) error {
	var (
		err, errs error
	)

	handleError := func(err error, entity string) error {
		if err == nil {
			return nil
		}

		if errors.Is(err, context.Canceled) {
			return err
		}

		errs = errors.Join(errs, fmt.Errorf("export %s: %w", entity, err))
		return nil
	}

	// Pipelines
	pbPipelines := make([]*typespb.Pipeline, 0, len(data.Pipelines))
	for _, p := range data.Pipelines {
		pbPipelines = append(pbPipelines, types.ConvertPipeline(p))
	}
	err = c.Exporter.ExportPipelines(ctx, pbPipelines)
	if herr := handleError(err, "pipelines"); herr != nil {
		return herr
	}

	// Jobs
	pbJobs := make([]*typespb.Job, 0, len(data.Jobs))
	for _, j := range data.Jobs {
		pbJobs = append(pbJobs, types.ConvertJob(j))
	}
	err = c.Exporter.ExportJobs(ctx, pbJobs)
	if herr := handleError(err, "jobs"); herr != nil {
		return herr
	}

	// Sections
	pbSections := make([]*typespb.Section, 0, len(data.Sections))
	for _, s := range data.Sections {
		pbSections = append(pbSections, types.ConvertSection(s))
	}
	err = c.Exporter.ExportSections(ctx, pbSections)
	if herr := handleError(err, "sections"); herr != nil {
		return herr
	}

	// Metrics
	pbMetrics := make([]*typespb.Metric, 0, len(data.Metrics))
	for _, m := range data.Metrics {
		pbMetrics = append(pbMetrics, types.ConvertMetric(m))
	}
	err = c.Exporter.ExportMetrics(ctx, pbMetrics)
	if herr := handleError(err, "metrics"); herr != nil {
		return herr
	}

	// Traces
	pbTraces := data.Traces
	err = c.Exporter.ExportTraces(ctx, pbTraces)
	if herr := handleError(err, "traces"); herr != nil {
		return herr
	}

	// Test Reports
	pbTestReports := make([]*typespb.TestReport, 0, len(data.TestReports))
	for _, tr := range data.TestReports {
		pbTestReports = append(pbTestReports, types.ConvertTestReport(tr))
	}
	err = c.Exporter.ExportTestReports(ctx, pbTestReports)
	if herr := handleError(err, "test reports"); herr != nil {
		return herr
	}

	// Test Suites
	pbTestSuites := make([]*typespb.TestSuite, 0, len(data.TestSuites))
	for _, ts := range data.TestSuites {
		pbTestSuites = append(pbTestSuites, types.ConvertTestSuite(ts))
	}
	err = c.Exporter.ExportTestSuites(ctx, pbTestSuites)
	if herr := handleError(err, "test suites"); herr != nil {
		return herr
	}

	// Test Cases
	pbTestCases := make([]*typespb.TestCase, 0, len(data.TestCases))
	for _, tc := range data.TestCases {
		pbTestCases = append(pbTestCases, types.ConvertTestCase(tc))
	}
	err = c.Exporter.ExportTestCases(ctx, pbTestCases)
	if herr := handleError(err, "test cases"); herr != nil {
		return herr
	}

	// Coverage Reports
	pbCoverageReports := make([]*typespb.CoverageReport, 0, len(data.CoverageReports))
	for _, cr := range data.CoverageReports {
		pbCoverageReports = append(pbCoverageReports, types.ConvertCoverageReport(cr))
	}
	err = c.Exporter.ExportCoverageReports(ctx, pbCoverageReports)
	if herr := handleError(err, "coverage reports"); herr != nil {
		return herr
	}

	// Coverage Packages
	pbCoveragePackages := make([]*typespb.CoveragePackage, 0, len(data.CoveragePackages))
	for _, cp := range data.CoveragePackages {
		pbCoveragePackages = append(pbCoveragePackages, types.ConvertCoveragePackage(cp))
	}
	err = c.Exporter.ExportCoveragePackages(ctx, pbCoveragePackages)
	if herr := handleError(err, "coverage packages"); herr != nil {
		return herr
	}

	// Coverage Classes
	pbCoverageClasses := make([]*typespb.CoverageClass, 0, len(data.CoverageClasses))
	for _, cc := range data.CoverageClasses {
		pbCoverageClasses = append(pbCoverageClasses, types.ConvertCoverageClass(cc))
	}
	err = c.Exporter.ExportCoverageClasses(ctx, pbCoverageClasses)
	if herr := handleError(err, "coverage packages"); herr != nil {
		return herr
	}

	// Coverage Methods
	pbCoverageMethods := make([]*typespb.CoverageMethod, 0, len(data.CoverageMethods))
	for _, cm := range data.CoverageMethods {
		pbCoverageMethods = append(pbCoverageMethods, types.ConvertCoverageMethod(cm))
	}
	err = c.Exporter.ExportCoverageMethods(ctx, pbCoverageMethods)
	if herr := handleError(err, "coverage reports"); herr != nil {
		return herr
	}

	return errs
}

func (c *Controller) convertTraceData(pipelines []types.Pipeline, jobs []types.Job, sections []types.Section) (*tracepb.TracesData, error) {
	var (
		pipelineSpans  []*tracepb.Span
		buildJobSpans  []*tracepb.Span
		bridgeJobSpans []*tracepb.Span
		sectionSpans   []*tracepb.Span

		resourceSpans []*tracepb.ResourceSpans
	)

	for _, p := range pipelines {
		if c.projectsSettings.ExportTraces(p.Project.Id) {
			pipelineSpans = append(pipelineSpans, types.PipelineSpan(p))
		}
	}
	if len(pipelineSpans) > 0 {
		resourceSpans = append(resourceSpans,
			types.NewResourceSpan(
				map[string]string{
					"service.name": "gitlab_ci.pipeline",
				},
				pipelineSpans,
			),
		)
	}

	for _, j := range jobs {
		if c.projectsSettings.ExportTraces(j.Pipeline.Project.Id) {
			if j.Kind == types.JobKindBuild {
				buildJobSpans = append(buildJobSpans, types.JobSpan(j))
			} else if j.Kind == types.JobKindBridge {
				bridgeJobSpans = append(bridgeJobSpans, types.JobSpan(j))
			}
		}
	}
	if len(buildJobSpans) > 0 {
		resourceSpans = append(resourceSpans,
			types.NewResourceSpan(
				map[string]string{
					"service.name": "gitlab_ci.job",
				},
				buildJobSpans,
			),
		)
	}
	if len(bridgeJobSpans) > 0 {
		resourceSpans = append(resourceSpans,
			types.NewResourceSpan(
				map[string]string{
					"service.name": "gitlab_ci.bridge",
				},
				bridgeJobSpans,
			),
		)
	}

	for _, s := range sections {
		if c.projectsSettings.ExportTraces(s.Job.Pipeline.Project.Id) {
			sectionSpans = append(sectionSpans, types.SectionSpan(s))
		}
	}
	if len(sectionSpans) > 0 {
		resourceSpans = append(resourceSpans,
			types.NewResourceSpan(
				map[string]string{
					"service.name": "gitlab_ci.section",
				},
				sectionSpans,
			),
		)
	}

	return &tracepb.TracesData{
		ResourceSpans: resourceSpans,
	}, nil
}

type projectsMrData struct {
	MergeRequests          []types.MergeRequest
	MergeRequestNoteEvents []types.MergeRequestNoteEvent
}

func (c *Controller) processProjectsMrData(ctx context.Context, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) error {
	data, err := c.fetchProjectsMrData(ctx, projectIds, updatedAfter, updatedBefore)
	if err != nil {
		return fmt.Errorf("fetch projects mr data: %w", err)
	}

	if err := c.exportProjectsMrData(ctx, data); err != nil {
		return fmt.Errorf("export projects mr data: %w", err)
	}

	return nil
}

func (c *Controller) fetchProjectsMrData(ctx context.Context, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) (projectsMrData, error) {
	mergeRequests, err := FetchProjectsMergeRequests(ctx, c.GitLab, projectIds, updatedAfter, updatedBefore)
	if err != nil {
		return projectsMrData{}, err
	}

	mergeRequestNoteEvents, err := FetchProjectsMergeRequestsNotes(ctx, c.GitLab, projectIds, updatedAfter, updatedBefore)
	if err != nil {
		return projectsMrData{}, err
	}

	return projectsMrData{
		MergeRequests:          mergeRequests,
		MergeRequestNoteEvents: mergeRequestNoteEvents,
	}, nil
}

func (c *Controller) exportProjectsMrData(ctx context.Context, data projectsMrData) error {
	pbMergeRequests := make([]*typespb.MergeRequest, 0, len(data.MergeRequests))
	for _, mr := range data.MergeRequests {
		pbMergeRequests = append(pbMergeRequests, types.ConvertMergeRequest(mr))
	}

	pbMergeRequestNoteEvents := make([]*typespb.MergeRequestNoteEvent, 0, len(data.MergeRequestNoteEvents))
	for _, ne := range data.MergeRequestNoteEvents {
		if ne.Type != "" {
			pbMergeRequestNoteEvents = append(pbMergeRequestNoteEvents, types.ConvertMergeRequestNoteEvent(ne))
		}
	}

	if err := c.Exporter.ExportMergeRequests(ctx, pbMergeRequests); err != nil {
		return fmt.Errorf("export merge requests: %w", err)
	}

	if err := c.Exporter.ExportMergeRequestNoteEvents(ctx, pbMergeRequestNoteEvents); err != nil {
		return fmt.Errorf("export merge request note events: %w", err)
	}

	return nil
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
