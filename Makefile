.DEFAULT_GOAL := help

SCHEMA_CURL_URL ?= https://system.cc.catonetworks.com/api/schema?with_undocumented=true
SCHEMA_FILE ?= cato_api.graphqls
PATCH_DIR ?= schema-patches
PATCH_FILES := $(sort $(wildcard $(PATCH_DIR)/*.patch))

##@ Generator
.PHONY: generate

generate: ## Generate client and models
	go tool gqlgenc

.PHONY: schema-update
schema-update: ## Update cato_api schema using curl source + normalize
	@tmp_file="$$(mktemp)"; \
	archive_dir="archives"; \
	archive_file="$$archive_dir/cato_api-$$(date +%Y%m%d).graphqls"; \
	set -e; \
	if [ -f "$(SCHEMA_FILE)" ]; then \
		mkdir -p "$$archive_dir"; \
		cp "$(SCHEMA_FILE)" "$$archive_file"; \
		echo "Archived $(SCHEMA_FILE) -> $$archive_file"; \
	fi; \
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
