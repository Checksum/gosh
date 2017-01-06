install:
	@echo "Building and installing gosh"
	@go install -ldflags "-s -w"
	@echo "Done"

release:
	@echo "Building release"
	@gox -ldflags "-s -w" -os="linux darwin windows openbsd" -arch="amd64 386" -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
	@echo "Done"

.PHONY: install release
