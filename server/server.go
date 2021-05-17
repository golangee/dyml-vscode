package tadl

import (
	"fmt"
	"strings"
	"tadl/protocol"
)

// Tadl language server.
type Server struct {
	// Map from Uri's to files.
	files map[protocol.DocumentURI]File
}

func NewServer() Server {
	return Server{
		files: make(map[protocol.DocumentURI]File),
	}
}

// Handle a client's request to initialize and respond with our capabilities.
func (s *Server) Initialize(params *protocol.InitializeParams) protocol.InitializeResult {
	// For a perfect server we would need to check params.Capabilities to know
	// what information the client can handle.
	return protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.Full,
			SemanticTokensProvider: protocol.SemanticTokensOptions{
				Legend: protocol.SemanticTokensLegend{
					TokenTypes: []string{"keyword", "variable", "assignment", "number"},
				},
				Full: true,
			},
		},
	}
}

// Initialized tells us, that the client is ready.
func (s *Server) Initialized() {
	s.sendPreview()
}

// Handle a hover event.
func (t *Server) Hover(params *protocol.HoverParams) protocol.Hover {
	return protocol.Hover{} // Don't forget to enable hover capabilities when using this.
}

// A document was saved.
func (s *Server) DidSaveTextDocument(params *protocol.DidSaveTextDocumentParams) {
	s.sendFunnyDiagnostics()
	s.sendPreview()
}

// A document was opened.
func (s *Server) DidOpenTextDocument(params *protocol.DidOpenTextDocumentParams) {
	s.files[params.TextDocument.URI] = File{
		Uri:     params.TextDocument.URI,
		Content: params.TextDocument.Text,
	}
	s.sendFunnyDiagnostics()
	s.sendPreview()
}

// A document was close.
func (s *Server) DidCloseTextDocument(params *protocol.DidCloseTextDocumentParams) {
	delete(s.files, params.TextDocument.URI)
	s.sendPreview()
}

// A document was changed.
func (s *Server) DidChangeTextDocument(params *protocol.DidChangeTextDocumentParams) {
	// There is only a ever single full content change, as we requested.
	s.files[params.TextDocument.URI] = File{
		Uri:     params.TextDocument.URI,
		Content: params.ContentChanges[0].Text,
	}
	s.sendFunnyDiagnostics()
}

func (s *Server) FullSemanticTokens(params *protocol.SemanticTokensParams) protocol.SemanticTokens {
	// Mark "let" as a keyword for testing purposes.

	file := s.files[params.TextDocument.URI]
	var data []uint32

	// Find absolute positions of "let"
	for lineIdx, line := range strings.Split(file.Content, "\n") {
		for i := 0; i < len(line); i++ {
			if i < len(line)-2 {
				if line[i:i+3] == "let" {
					data = append(data, uint32(lineIdx), uint32(i), 3, 0, 0)
				}
			}
		}
	}

	// Make token positions relative.
	// Tokens are always 5 ints, first entry is line, second is char.
	for i := len(data) - 5; i >= 5; i -= 5 {
		// Make line difference relativ to previous
		data[i] -= data[i-5]
		// If item is in the same line, make char difference relative to previous
		if data[i] == 0 {
			data[i+1] -= data[i-5+1]
		}
	}

	return protocol.SemanticTokens{
		Data: data,
	}
}

// Send some kind of diagnostics to test it out.
func (s *Server) sendFunnyDiagnostics() {

	for _, file := range s.files {

		fileContent := file.Content
		fileContent = strings.ToLower(fileContent)

		diagnostics := []protocol.Diagnostic{}
		for lineIdx, line := range strings.Split(fileContent, "\n") {
			i := strings.Index(line, "servus")
			if i >= 0 {
				diagnostics = append(diagnostics, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{Line: uint32(lineIdx), Character: uint32(i)},
						End:   protocol.Position{Line: uint32(lineIdx), Character: uint32(i + 6)},
					},
					Severity: protocol.SeverityError,
					Message:  "Unzulässige Anrede",
					Source:   "bayern-lint",
				})
			}
		}

		_ = SendNotification("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
			URI:         protocol.DocumentURI(file.Uri),
			Diagnostics: diagnostics,
		})
	}
}

func (s *Server) sendPreview() {
	text := fmt.Sprintf(`
<html>
	<head>
		<style>
		</style>
	</head>
	<body>
		<h1>Tadl Preview</h1>
		<hr>
		Du hast %d Dateien geöffnet.<br>
	</body>
</html>`, len(s.files))

	_ = SendNotification("custom/preview", text)
}
