package rest_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/rest"
)

func TestParseJobLog_OnlySections(t *testing.T) {
	trace := []byte(`
    [0KRunning with gitlab-runner 16.6.0~beta.105.gd2263193 (d2263193)[0;m
    [0K  on blue-4.saas-linux-small-amd64.runners-manager.gitlab.com/default J2nyww-s, system ID: s_cf1798852952[0;m
    [0K  feature flags: FF_USE_IMPROVED_URL_MASKING:true[0;m
    section_start:1700819846:prepare_executor
    [0K[0K[36;1mPreparing the "docker+machine" executor[0;m[0;m
    [0KUsing Docker executor with image golang:1.20 ...[0;m
    [0KPulling docker image golang:1.20 ...[0;m
    [0KUsing docker image sha256:f0a2018ec55c82e1734258f77fe3afcd430f1ccfe351a1caa7d61bf5de595247 for golang:1.20 with digest golang@sha256:4e4a34f7940eddba81c1f6df88057411bb2d822df087c317f8532cc169f2725a ...[0;m
    section_end:1700819865:prepare_executor
    [0Ksection_start:1700819865:prepare_script
    [0K[0K[36;1mPreparing environment[0;m[0;m
    Running on runner-j2nyww-s-project-50817395-concurrent-0 via runner-j2nyww-s-s-l-s-amd64-1700819808-26e7b50a...
    section_end:1700819868:prepare_script
    [0Ksection_start:1700819868:get_sources
    [0K[0K[36;1mGetting source from Git repository[0;m[0;m
    [32;1mFetching changes with git depth set to 20...[0;m
    Initialized empty Git repository in /builds/cluttrdev/gitlab-exporter/.git/
    [32;1mCreated fresh repository.[0;m
    [32;1mChecking out 3c530002 as detached HEAD (ref is main)...[0;m
    
    [32;1mSkipping Git submodules setup[0;m
    [32;1m$ git remote set-url origin "${CI_REPOSITORY_URL}"[0;m
    section_end:1700819869:get_sources
    [0Ksection_start:1700819869:restore_cache
    [0K[0K[36;1mRestoring cache[0;m[0;m
    [32;1mChecking cache for main-protected...[0;m
    Downloading cache from https://storage.googleapis.com/gitlab-com-runners-cache/project/50817395/main-protected[0;m 
    [32;1mSuccessfully extracted cache[0;m
    section_end:1700819875:restore_cache
    [0Ksection_start:1700819875:step_script
    [0K[0K[36;1mExecuting "step_script" stage of the job script[0;m[0;m
    [0KUsing docker image sha256:f0a2018ec55c82e1734258f77fe3afcd430f1ccfe351a1caa7d61bf5de595247 for golang:1.20 with digest golang@sha256:4e4a34f7940eddba81c1f6df88057411bb2d822df087c317f8532cc169f2725a ...[0;m
    [32;1m$ go build .[0;m
    section_end:1700819916:step_script
    [0Ksection_start:1700819916:archive_cache
    [0K[0K[36;1mSaving cache for successful job[0;m[0;m
    [32;1mNot uploading cache main-protected due to policy[0;m
    section_end:1700819916:archive_cache
    [0Ksection_start:1700819916:cleanup_file_variables
    [0K[0K[36;1mCleaning up project directory and file based variables[0;m[0;m
    section_end:1700819917:cleanup_file_variables
    [0K[32;1mJob succeeded[0;m
    `)

	data, err := rest.ParseJobLog(bytes.NewReader(trace))
	if err != nil {
		t.Errorf("%v", err)
	}

	expected := rest.JobLogData{
		Sections: []rest.SectionData{
			{Name: "prepare_executor", Start: 1700819846, End: 1700819865},
			{Name: "prepare_script", Start: 1700819865, End: 1700819868},
			{Name: "get_sources", Start: 1700819868, End: 1700819869},
			{Name: "restore_cache", Start: 1700819869, End: 1700819875},
			{Name: "step_script", Start: 1700819875, End: 1700819916},
			{Name: "archive_cache", Start: 1700819916, End: 1700819916},
			{Name: "cleanup_file_variables", Start: 1700819916, End: 1700819917},
		},
		Metrics: nil,
	}

	if diff := cmp.Diff(expected, data); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}
}

func TestParseJobLog_OnlySections_Unfinished(t *testing.T) {
	trace := []byte(`
    [0KRunning with gitlab-runner 16.6.0~beta.105.gd2263193 (d2263193)[0;m
    [0K  on blue-4.saas-linux-small-amd64.runners-manager.gitlab.com/default J2nyww-s, system ID: s_cf1798852952[0;m
    [0K  feature flags: FF_USE_IMPROVED_URL_MASKING:true[0;m
    section_start:1700819846:prepare_executor
    [0K[0K[36;1mPreparing the "docker+machine" executor[0;m[0;m
    [0KUsing Docker executor with image golang:1.20 ...[0;m
    [0KPulling docker image golang:1.20 ...[0;m
    [0KUsing docker image sha256:f0a2018ec55c82e1734258f77fe3afcd430f1ccfe351a1caa7d61bf5de595247 for golang:1.20 with digest golang@sha256:4e4a34f7940eddba81c1f6df88057411bb2d822df087c317f8532cc169f2725a ...[0;m
    section_end:1700819865:prepare_executor
    [0Ksection_start:1700819865:prepare_script
    [0K[0K[36;1mPreparing environment[0;m[0;m
    Running on runner-j2nyww-s-project-50817395-concurrent-0 via runner-j2nyww-s-s-l-s-amd64-1700819808-26e7b50a...
    section_end:1700819868:prepare_script
    [0Ksection_start:1700819868:get_sources
    [0K[0K[36;1mGetting source from Git repository[0;m[0;m
    [32;1mFetching changes with git depth set to 20...[0;m
    Initialized empty Git repository in /builds/cluttrdev/gitlab-exporter/.git/
    [32;1mCreated fresh repository.[0;m
    [32;1mChecking out 3c530002 as detached HEAD (ref is main)...[0;m

    [32;1mSkipping Git submodules setup[0;m
    [32;1m$ git remote set-url origin "${CI_REPOSITORY_URL}"[0;m
    section_end:1700819869:get_sources
    [0Ksection_start:1700819869:restore_cache
    [0K[0K[36;1mRestoring cache[0;m[0;m
    [32;1mChecking cache for main-protected...[0;m
    Downloading cache from https://storage.googleapis.com/gitlab-com-runners-cache/project/50817395/main-protected[0;m
    [32;1mSuccessfully extracted cache[0;m
    section_end:1700819875:restore_cache
    [0Ksection_start:1700819875:step_script
    [0K[0K[36;1mExecuting "step_script" stage of the job script[0;m[0;m
    [0KUsing docker image sha256:f0a2018ec55c82e1734258f77fe3afcd430f1ccfe351a1caa7d61bf5de595247 for golang:1.20 with digest golang@sha256:4e4a34f7940eddba81c1f6df88057411bb2d822df087c317f8532cc169f2725a ...[0;m
    [0Ksection_start:1700819875:script_step_build
    [32;1m$ go build .[0;m
    `)

	data, err := rest.ParseJobLog(bytes.NewReader(trace))
	if err != nil {
		t.Errorf("%v", err)
	}

	expected := rest.JobLogData{
		Sections: []rest.SectionData{
			{Name: "prepare_executor", Start: 1700819846, End: 1700819865},
			{Name: "prepare_script", Start: 1700819865, End: 1700819868},
			{Name: "get_sources", Start: 1700819868, End: 1700819869},
			{Name: "restore_cache", Start: 1700819869, End: 1700819875},
			{Name: "script_step_build", Start: 1700819875, End: 0},
			{Name: "step_script", Start: 1700819875, End: 0},
		},
		Metrics: nil,
	}

	if diff := cmp.Diff(expected, data); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}
}

func TestParseJobLog_OnlyMetrics(t *testing.T) {
	trace := []byte(`
METRIC_minimal_metric 42.1337
METRIC_with_timestamp 3.14 1234567890
METRIC_no_labels{} 3
METRIC_single_label{label_name="label_value"} 42
METRIC_multi_labels{label_name1="label_value1",label_name2="label_value2"} 1
    `)

	data, err := rest.ParseJobLog(bytes.NewReader(trace))
	if err != nil {
		t.Errorf("%v", err)
	}

	expected := rest.JobLogData{
		Sections: nil,
		Metrics: []rest.MetricData{
			{
				Name:      "minimal_metric",
				Value:     42.1337,
				Labels:    nil,
				Timestamp: 0,
			},
			{
				Name:      "with_timestamp",
				Labels:    nil,
				Value:     3.14,
				Timestamp: 1234567890,
			},
			{
				Name:      "no_labels",
				Labels:    nil,
				Value:     3,
				Timestamp: 0,
			},
			{
				Name: "single_label",
				Labels: map[string]string{
					"label_name": "label_value",
				},
				Value:     42,
				Timestamp: 0,
			},
			{
				Name: "multi_labels",
				Labels: map[string]string{
					"label_name1": "label_value1",
					"label_name2": "label_value2",
				},
				Value:     1,
				Timestamp: 0,
			},
		},
	}

	if diff := cmp.Diff(expected, data); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}
}
