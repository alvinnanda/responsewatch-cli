REPO     = alvinnanda/responsewatch-cli
BINARY   = rwcli
DIST_DIR = ./dist

# Read version from git tag or default to 'dev'
VERSION  ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")

.PHONY: help build release upload tag clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build binaries for all platforms
	@echo "Building $(BINARY) version $(VERSION)..."
	@rm -rf $(DIST_DIR) && mkdir -p $(DIST_DIR)
	GOOS=darwin  GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY)_darwin_amd64  .
	GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY)_darwin_arm64  .
	GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY)_linux_amd64   .
	GOOS=linux   GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY)_linux_arm64   .
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(DIST_DIR)/$(BINARY)_windows_amd64.exe .
	@echo "✅ Build complete:"
	@ls -lh $(DIST_DIR)

tag: ## Create and push a new git tag. Usage: make tag VERSION=v1.2.0
ifndef VERSION
	$(error VERSION is not set. Usage: make tag VERSION=v1.2.0)
endif
	@echo "Creating tag $(VERSION)..."
	git tag $(VERSION)
	git push origin $(VERSION)
	@echo "✅ Tag $(VERSION) pushed"

release: build ## Build, create GitHub Release, and upload all binaries. Usage: make release VERSION=v1.2.0
ifndef VERSION
	$(error VERSION is not set. Usage: make release VERSION=v1.2.0)
endif
	@echo "Creating GitHub Release $(VERSION)..."
	@RELEASE_ID=$$(curl -s -X POST \
		-H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer $$(git credential fill <<< $$'protocol=https\nhost=github.com' 2>/dev/null | grep password | cut -d= -f2)" \
		https://api.github.com/repos/$(REPO)/releases \
		-d "{\"tag_name\":\"$(VERSION)\",\"name\":\"$(VERSION)\",\"body\":\"Release $(VERSION)\",\"draft\":false,\"prerelease\":false}" \
		| python3 -c "import sys,json; print(json.load(sys.stdin)['id'])"); \
	echo "Release ID: $$RELEASE_ID"; \
	for BINARY_FILE in $(DIST_DIR)/*; do \
		NAME=$$(basename $$BINARY_FILE); \
		echo "Uploading $$NAME..."; \
		curl -s -X POST \
			-H "Accept: application/vnd.github+json" \
			-H "Authorization: Bearer $$(git credential fill <<< $$'protocol=https\nhost=github.com' 2>/dev/null | grep password | cut -d= -f2)" \
			-H "Content-Type: application/octet-stream" \
			"https://uploads.github.com/repos/$(REPO)/releases/$$RELEASE_ID/assets?name=$$NAME" \
			--data-binary "@$$BINARY_FILE" > /dev/null && echo "  ✅ $$NAME"; \
	done
	@echo ""
	@echo "🚀 Release $(VERSION) published!"
	@echo "   https://github.com/$(REPO)/releases/tag/$(VERSION)"

upload: ## Upload binaries to an existing release. Usage: make upload VERSION=v1.2.0
ifndef VERSION
	$(error VERSION is not set. Usage: make upload VERSION=v1.2.0)
endif
	@echo "Fetching release ID for $(VERSION)..."
	@RELEASE_ID=$$(curl -s \
		-H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer $$(git credential fill <<< $$'protocol=https\nhost=github.com' 2>/dev/null | grep password | cut -d= -f2)" \
		https://api.github.com/repos/$(REPO)/releases/tags/$(VERSION) \
		| python3 -c "import sys,json; print(json.load(sys.stdin)['id'])"); \
	echo "Release ID: $$RELEASE_ID"; \
	for BINARY_FILE in $(DIST_DIR)/*; do \
		NAME=$$(basename $$BINARY_FILE); \
		echo "Uploading $$NAME..."; \
		curl -s -X POST \
			-H "Accept: application/vnd.github+json" \
			-H "Authorization: Bearer $$(git credential fill <<< $$'protocol=https\nhost=github.com' 2>/dev/null | grep password | cut -d= -f2)" \
			-H "Content-Type: application/octet-stream" \
			"https://uploads.github.com/repos/$(REPO)/releases/$$RELEASE_ID/assets?name=$$NAME" \
			--data-binary "@$$BINARY_FILE" > /dev/null && echo "  ✅ $$NAME"; \
	done
	@echo "✅ Upload complete"

clean: ## Remove build artifacts
	rm -rf $(DIST_DIR)
	@echo "✅ Cleaned"
