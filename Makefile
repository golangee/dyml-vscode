vsix: server
	vsce package

server:
	make -C server

.PHONY: server
