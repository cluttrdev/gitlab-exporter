ifneq "${VERBOSE}" "1"
.SILENT:
endif

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help page
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: graphql-schema
graphql-schema: ## Generate the GitLab GraphQL schema
	go tool -modfile go.tool.mod gqlfetch --endpoint https://gitlab.com/api/graphql > internal/gitlab/graphql/spec/schema.graphql

.PHONY: graphql
graphql: ## Generate the GitLab GraphQL API client code
	go tool -modfile go.tool.mod genqlient internal/gitlab/graphql/genqlient.yaml
