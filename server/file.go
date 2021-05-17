package tadl

import "tadl/protocol"

// File is a file that is located at an Uri and has Content.
type File struct {
	Uri     protocol.DocumentURI
	Content string
}
