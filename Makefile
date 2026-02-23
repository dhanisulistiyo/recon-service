
.Phony: test build

all: test build

test:
	@echo "=================================================================================="
	@echo "Coverage Test"
	@echo "=================================================================================="
	go test ./... -coverprofile coverage.cov
	@echo "\n"
	@echo "=================================================================================="
	@echo "All Package Coverage"
	@echo "=================================================================================="
	go tool cover -func coverage.cov

build:
	docker build --no-cache .