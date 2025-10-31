.PHONY: test-all test-ffprobe test-ffmpeg-thumbnail test-ffmpeg-split test-ffmpeg-concat test-ffmpeg test-shaka-packager help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

test-all: test-ffprobe test-ffmpeg-thumbnail test-ffmpeg-split test-ffmpeg-concat test-ffmpeg test-shaka-packager ## Run all E2E tests

test-ffprobe: ## Run ffprobe E2E tests
	@echo "Running ffprobe tests..."
	go test -v -timeout 5m ./ffprobe/...

test-ffmpeg-thumbnail: ## Run ffmpeg-thumbnail E2E tests
	@echo "Running ffmpeg-thumbnail tests..."
	go test -v -timeout 5m ./ffmpeg-thumbnail/...

test-ffmpeg-split: ## Run ffmpeg-split E2E tests
	@echo "Running ffmpeg-split tests..."
	go test -v -timeout 5m ./ffmpeg-split/...

test-ffmpeg-concat: ## Run ffmpeg-concat E2E tests
	@echo "Running ffmpeg-concat tests..."
	go test -v -timeout 5m ./ffmpeg-concat/...

test-ffmpeg: ## Run ffmpeg E2E tests
	@echo "Running ffmpeg tests..."
	go test -v -timeout 10m ./ffmpeg/...

test-shaka-packager: ## Run shaka-packager E2E tests
	@echo "Running shaka-packager tests..."
	go test -v -timeout 15m ./shaka-packager/...

# Clean test artifacts
clean-test: ## Clean all test output directories
	@echo "Cleaning test artifacts..."
	find testdata -type d -name '[0-9a-f]*' -exec rm -rf {} + 2>/dev/null || true
	@echo "Clean complete"

# Initialize test environment
test-setup: ## Setup test environment (download sample video if needed)
	@echo "Setting up test environment..."
	@mkdir -p testdata
	@if [ ! -f testdata/sample.mp4 ]; then \
		echo "Downloading sample video..."; \
		curl -L "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4" \
			-o testdata/sample.mp4 --create-dirs; \
	fi
	@echo "Test environment ready"
