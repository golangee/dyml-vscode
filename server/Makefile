# Build binaries to the out directory.
# The filenames should be chosen so that Node's OS API (https://nodejs.org/api/os.html#os_os_arch)
# can be used to select a binary.
build: test
	GOOS=linux GOARCH=amd64 go build -o ../out/bin/tadl-linux-x64 cmd/tadl.go
	GOOS=darwin GOARCH=amd64 go build -o ../out/bin/tadl-darwin-x64 cmd/tadl.go
	GOOS=darwin GOARCH=arm64 go build -o ../out/bin/tadl-darwin-arm64 cmd/tadl.go

test:
	golangci-lint run

# Download the newest LSP types from https://github.com/golang/tools.
# License and code will be copied to the protocol directory.
protocol:
	rm -rf /tmp/tadl/go-tools
	git clone https://github.com/golang/tools.git /tmp/tadl/go-tools
	mkdir protocol || true
	cp /tmp/tadl/go-tools/LICENSE protocol
	cp /tmp/tadl/go-tools/internal/lsp/protocol/tsprotocol.go protocol

.PHONY: protocol