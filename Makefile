help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

lint: ## Run linter
	go mod tidy
	golangci-lint run

test: ## Run go tests
	go test ./... -vet=off -count=1
