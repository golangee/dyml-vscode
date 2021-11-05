# Build binaries to the out directory.
# The filenames should be chosen so that Node's OS API (https://nodejs.org/api/os.html#os_os_arch)
# can be used to select a binary.
build: test
	GOOS=linux GOARCH=amd64 go build -o ../out/bin/dyml-linux-x64 cmd/dyml.go
	GOOS=darwin GOARCH=amd64 go build -o ../out/bin/dyml-darwin-x64 cmd/dyml.go
	# Skip building of darwin-arm64 until github actions can do that.
	# GOOS=darwin GOARCH=arm64 go build -o ../out/bin/dyml-darwin-arm64 cmd/dyml.go

test:
	golangci-lint run || true

# Download the newest LSP types from https://github.com/golang/tools.
# License and code will be copied to the protocol directory.
protocol:
	rm -rf /tmp/dyml/go-tools
	git clone https://github.com/golang/tools.git /tmp/dyml/go-tools
	mkdir protocol || true
	cp /tmp/dyml/go-tools/LICENSE protocol
	cp /tmp/dyml/go-tools/internal/lsp/protocol/tsprotocol.go protocol

.PHONY: protocol
