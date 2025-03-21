package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	gitlab_mock "go.cluttr.dev/gitlab-exporter/test/mock/gitlab"
)

const testSet string = "default"

func TestMain(m *testing.M) {
	fh := &gitlab_mock.FileHandler{Root: http.Dir("testdata/gitlab.com")}
	srv := httptest.NewServer(fh)
	defer srv.Close()

	testEnv := GitLabExporterTestEnvironment{
		URL: srv.URL,
	}

	SetTestEnvironment(testSet, testEnv)

	os.Exit(m.Run())
}

type GitLabExporterTestEnvironment struct {
	URL string
}

func SetTestEnvironment(name string, env GitLabExporterTestEnvironment) {
	bytes, err := json.Marshal(env)
	if err != nil {
		panic(err)
	}

	os.Setenv(fmt.Sprintf("GLE_TEST_ENV_%s", strings.ToUpper(testSet)), string(bytes))
}

func GetTestEnvironment(name string) (GitLabExporterTestEnvironment, error) {
	sEnv := os.Getenv(fmt.Sprintf("GLE_TEST_ENV_%s", strings.ToUpper(testSet)))
	if sEnv == "" {
		return GitLabExporterTestEnvironment{}, fmt.Errorf("environment not found: %s", name)
	}

	var env GitLabExporterTestEnvironment
	if err := json.Unmarshal([]byte(sEnv), &env); err != nil {
		return GitLabExporterTestEnvironment{}, err
	}

	return env, nil
}
