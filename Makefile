.DEFAULT_GOAL := help

SCHEMA_CURL_URL ?= https://system.cc.catonetworks.com/api/schema?with_undocumented=true
SCHEMA_FILE ?= cato_api.graphqls
PATCH_DIR ?= schema-patches
PATCH_FILES := $(sort $(wildcard $(PATCH_DIR)/*.patch))
CLI_ROOT ?= ../cato-cli
EXPECTED_OPERATIONS ?= 544

##@ Generator
.PHONY: generate operations-import operations-check generate-check

generate: operations-check ## Validate operations, then generate client and models
	go tool gqlgenc

operations-import: ## Import validated GraphQL operations from cato-cli
	go run ./cmd/gqlops import --cli-root "$(CLI_ROOT)" --expected "$(EXPECTED_OPERATIONS)"

operations-check: ## Validate canonical GraphQL operations and manifest
	go run ./cmd/gqlops check --expected "$(EXPECTED_OPERATIONS)"

generate-check: operations-check ## Require generated client and models to be current and deterministic
	@tmp_dir="$$(mktemp -d)"; \
	set -e; \
	trap 'rm -rf "$$tmp_dir"' EXIT; \
	mkdir -p "$$tmp_dir/models"; \
	cp client.go "$$tmp_dir/client.go"; \
	cp models/models.go "$$tmp_dir/models/models.go"; \
	go tool gqlgenc; \
	cmp client.go "$$tmp_dir/client.go"; \
	cmp models/models.go "$$tmp_dir/models/models.go"

.PHONY: schema-update
schema-update: ## Update cato_api schema using curl source + normalize
	@tmp_file="$$(mktemp)"; \
	set -e; \
	echo "Fetching schema from: $(SCHEMA_CURL_URL)"; \
	curl -fsSL "$(SCHEMA_CURL_URL)" -o "$$tmp_file"; \
	go run ./cmd/gqlschema normalize -f "$$tmp_file" -o "$(SCHEMA_FILE)"; \
	rm -f "$$tmp_file"; \
	echo "Updated $(SCHEMA_FILE)"

.PHONY: apply-patches
apply-patches: ## Apply schema patches from $(PATCH_DIR)
	@set -e; \
	if [ -z "$(PATCH_FILES)" ]; then \
		echo "No patch files found in $(PATCH_DIR)"; \
		exit 0; \
	fi; \
	for p in $(PATCH_FILES); do \
		if git apply --reverse --check "$$p" >/dev/null 2>&1; then \
			echo "Skipping already applied patch: $$p"; \
			continue; \
		fi; \
		echo "Applying patch: $$p"; \
		git apply "$$p"; \
	done

##@ Help
.PHONY: help
help: ## Help screen
	@awk -F ':.*##' 'BEGIN{printf "Usage: make <target>"} /^##@/ {  printf "\n\033[1m%s\033[0m\n", substr($$0, 5); next } \
	/^[a-zA-Z0-9_ -]+:.*?## .*$$/ { split($$1, targets, " "); \
	for (target in targets) {  printf "  \033[36m%-30s\033[0m %s\n", targets[target], substr($$0, index($$0,"##")+3) }}' $(MAKEFILE_LIST)
