package rest_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/rest"
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

func TestParseJobLog_OnlyProperties(t *testing.T) {
	trace := []byte(`
PROPERTY_property_name="property_value"
PROPERTY_empty_value=""
PROPERTY_escaped_quotes="\"42\""
PROPERTY_escaped_newlines="This is\\na multi-line\\nvalue."
PROPERTY_invalid_only_name
PROPERTY_invalid_missing_value=
    `)

	data, err := rest.ParseJobLog(bytes.NewReader(trace))
	if err != nil {
		t.Errorf("%v", err)
	}

	expected := rest.JobLogData{
		Sections: nil,
		Metrics:  nil,
		Properties: []rest.PropertyData{
			{
				Name:  "property_name",
				Value: "property_value",
			},
			{
				Name:  "empty_value",
				Value: "",
			},
			{
				Name:  "escaped_quotes",
				Value: "\"42\"",
			},
			{
				Name:  "escaped_newlines",
				Value: "This is\\na multi-line\\nvalue.",
			},
		},
	}

	if diff := cmp.Diff(expected, data); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}
}

func TestParseJobLog_WithTimestamps(t *testing.T) {
	trace := []byte(`
2026-01-21T09:47:05.094219Z 00O [0KRunning with gitlab-runner 18.7.0~pre.433.g3a5f2314 (3a5f2314)[0;m
2026-01-21T09:47:05.094234Z 00O [0K  on green-5.saas-linux-small-amd64.runners-manager.gitlab.com/default xS6Vzpvoq, system ID: s_6b1e4f06fcfd[0;m
2026-01-21T09:47:05.094269Z 00O section_start:1768988825:prepare_executor
2026-01-21T09:47:05.094270Z 00O+[0K[0K[36;1mPreparing the "docker+machine" executor[0;m[0;m
2026-01-21T09:47:05.242666Z 00O [0KUsing Docker executor with image golang:1.24 ...[0;m
2026-01-21T09:47:09.874960Z 00O [0KUsing effective pull policy of [always] for container golang:1.24[0;m
2026-01-21T09:47:09.876722Z 00O [0KPulling docker image golang:1.24 ...[0;m
2026-01-21T09:47:26.582932Z 00O [0KUsing docker image sha256:f61098245d584f6b0b12f87e34986ff7618ff620d3a0165b2d28a9d885597888 for golang:1.24 with digest golang@sha256:c2131140c7c29ff277b1c412d524b7f56289513f49672c57a3d992247dd146f8 ...[0;m
2026-01-21T09:47:26.583005Z 00O section_end:1768988846:prepare_executor
2026-01-21T09:47:26.583007Z 00O+[0Ksection_start:1768988846:prepare_script
2026-01-21T09:47:26.583103Z 00O+[0K[0K[36;1mPreparing environment[0;m[0;m
2026-01-21T09:47:26.584705Z 00O [0KUsing effective pull policy of [always] for container sha256:025f4f2c1a0f27e2c4ea82cff9631d6980aac6571718761f8675ce67ebdac11d[0;m
2026-01-21T09:47:30.067573Z 01O Running on runner-xs6vzpvoq-project-50817395-concurrent-0 via runner-xs6vzpvoq-s-l-s-amd64-1768988782-af0dfd37...
2026-01-21T09:47:30.249548Z 00O section_end:1768988850:prepare_script
2026-01-21T09:47:30.249553Z 00O+[0Ksection_start:1768988850:get_sources
2026-01-21T09:47:30.250075Z 00O+[0K[0K[36;1mGetting source from Git repository[0;m[0;m
2026-01-21T09:47:30.639008Z 01O [32;1m$ printf 'PROPERTY_%s="%s"\n' "ci_runner_id" "${CI_RUNNER_ID}" # collapsed multi-line command[0;m
2026-01-21T09:47:30.639030Z 01O PROPERTY_ci_runner_id="12270859"
2026-01-21T09:47:30.639031Z 01O PROPERTY_ci_runner_version="18.7.0~pre.433.g3a5f2314"
2026-01-21T09:47:30.639031Z 01O PROPERTY_ci_runner_revision="3a5f2314"
2026-01-21T09:47:30.639032Z 01O PROPERTY_ci_runner_hostname="runner-xs6vzpvoq-project-50817395-concurrent-0"
2026-01-21T09:47:30.659538Z 01O [32;1mGitaly correlation ID: a4a357cdb4fb4423a1c8c1bd9b96016d[0;m
2026-01-21T09:47:30.664750Z 01O [32;1mFetching changes...[0;m
2026-01-21T09:47:30.668241Z 01O Initialized empty Git repository in /builds/gitlab-exporter/gitlab-exporter/.git/
2026-01-21T09:47:30.670771Z 01O [32;1mCreated fresh repository.[0;m
2026-01-21T09:47:31.407321Z 01O [32;1mChecking out 9833c929 as detached HEAD (ref is v0.23.0)...[0;m
2026-01-21T09:47:31.518956Z 01O 
2026-01-21T09:47:31.519031Z 01O [32;1mSkipping Git submodules setup[0;m
2026-01-21T09:47:31.519155Z 01O [32;1m$ git remote set-url origin "${CI_REPOSITORY_URL}" || echo 'Not a git repository; skipping'[0;m
2026-01-21T09:47:31.725607Z 00O section_end:1768988851:get_sources
2026-01-21T09:47:31.725612Z 00O+[0Ksection_start:1768988851:restore_cache
2026-01-21T09:47:31.727645Z 00O+[0K[0K[36;1mRestoring cache[0;m[0;m
2026-01-21T09:47:32.188078Z 01O [32;1mChecking cache for v0-23-0-protected...[0;m
2026-01-21T09:47:32.386480Z 01E Downloading cache from https://storage.googleapis.com/gitlab-com-runners-cache/project/50817395/v0-23-0-protected[0;m  ETag[0;m="701ae521fc7fd67943041f4fa42ecdcc"
2026-01-21T09:47:49.222685Z 01O [32;1mSuccessfully extracted cache[0;m
2026-01-21T09:47:58.161233Z 00O section_end:1768988878:restore_cache
2026-01-21T09:47:58.161237Z 00O+[0Ksection_start:1768988878:step_script
2026-01-21T09:47:58.161846Z 00O+[0K[0K[36;1mExecuting "step_script" stage of the job script[0;m[0;m
2026-01-21T09:47:58.161865Z 00O [0KUsing effective pull policy of [always] for container golang:1.24[0;m
2026-01-21T09:47:58.163101Z 00O [0KUsing docker image sha256:f61098245d584f6b0b12f87e34986ff7618ff620d3a0165b2d28a9d885597888 for golang:1.24 with digest golang@sha256:c2131140c7c29ff277b1c412d524b7f56289513f49672c57a3d992247dd146f8 ...[0;m
2026-01-21T09:47:58.676964Z 01O [32;1m$ make test REPORTS=1[0;m
2026-01-21T09:47:58.683888Z 01O Testing ./exporter/...
2026-01-21T09:49:04.655430Z 01O 
2026-01-21T09:49:04.655434Z 01O DONE 2 tests in 0.016s
2026-01-21T09:49:05.237211Z 00O section_end:1768988945:step_script
2026-01-21T09:49:05.286731Z 00O+[0Ksection_start:1768988945:after_script
2026-01-21T09:49:05.289529Z 00O+[0K[0K[36;1mRunning after_script[0;m[0;m
2026-01-21T09:49:05.787398Z 01O [32;1mRunning after script...[0;m
2026-01-21T09:49:05.787425Z 01O [32;1m$ echo "METRIC_go_test_success 1"[0;m
2026-01-21T09:49:05.787461Z 01O METRIC_go_test_success 1
2026-01-21T09:49:05.787462Z 01O [32;1m$ echo "METRIC_go_test_runtime{ref=\"$CI_COMMIT_REF_NAME\",short_sha=\"$CI_COMMIT_SHORT_SHA\"} $go_test_runtime"[0;m
2026-01-21T09:49:05.787465Z 01O METRIC_go_test_runtime{ref="main",short_sha="eb48b511"} 0.016
2026-01-21T09:49:05.787466Z 01O 
2026-01-21T09:49:05.787500Z 00O section_end:1768988945:after_script
	`)

	data, err := rest.ParseJobLog(bytes.NewReader(trace))
	if err != nil {
		t.Errorf("%v", err)
	}

	expected := rest.JobLogData{
		Sections: []rest.SectionData{
			{Name: "prepare_executor", Start: 1768988825, End: 1768988846},
			{Name: "prepare_script", Start: 1768988846, End: 1768988850},
			{Name: "get_sources", Start: 1768988850, End: 1768988851},
			{Name: "restore_cache", Start: 1768988851, End: 1768988878},
			{Name: "step_script", Start: 1768988878, End: 1768988945},
			{Name: "after_script", Start: 1768988945, End: 1768988945},
		},
		Metrics: []rest.MetricData{
			{
				Name:  "go_test_success",
				Value: 1,
			},
			{
				Name:   "go_test_runtime",
				Labels: map[string]string{"ref": "main", "short_sha": "eb48b511"},
				Value:  0.016,
			},
		},
		Properties: []rest.PropertyData{
			{
				Name:  "ci_runner_id",
				Value: "12270859",
			},
			{
				Name:  "ci_runner_version",
				Value: "18.7.0~pre.433.g3a5f2314",
			},
			{
				Name:  "ci_runner_revision",
				Value: "3a5f2314",
			},
			{
				Name:  "ci_runner_hostname",
				Value: "runner-xs6vzpvoq-project-50817395-concurrent-0",
			},
		},
	}

	if diff := cmp.Diff(expected, data); diff != "" {
		t.Errorf("Result mismatch (-want +got):\n%s", diff)
	}
}
