.PHONY: release
release:
	GOOS=linux GOARCH=amd64 go build -o ~/Downloads/janus-linux-amd64 github.com/dcb9/janus/cli/janus
	GOOS=darwin GOARCH=amd64 go build -o ~/Downloads/janus-darwin-amd64 github.com/dcb9/janus/cli/janus
