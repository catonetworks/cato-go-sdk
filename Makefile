.DEFAULT_GOAL := help

##@ Generator
.PHONY: generate

generate: ## Generate client and models
	go tool gqlgenc

##@ Help
.PHONY: help
help: ## Help screen
	@awk -F ':.*##' 'BEGIN{printf "Usage: make <target>"} /^##@/ {  printf "\n\033[1m%s\033[0m\n", substr($$0, 5); next } \
	/^[a-zA-Z0-9_ -]+:.*?## .*$$/ { split($$1, targets, " "); \
	for (target in targets) {  printf "  \033[36m%-30s\033[0m %s\n", targets[target], substr($$0, index($$0,"##")+3) }}' $(MAKEFILE_LIST)
