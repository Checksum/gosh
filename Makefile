install:
	@echo "Building and installing gosh"
	@go install -ldflags "-s -w"
	@echo "Done"

release:
	@echo "Building release"
	@gox -ldflags "-s -w" -osarch="darwin/amd64 linux/amd64 windows/amd64" -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
	@echo "Done"

.PHONY: install release
