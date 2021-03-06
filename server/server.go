package dyml

import (
	"dyml-support/protocol"
	"path/filepath"
	"strings"

	"github.com/golangee/dyml/encoder"
	"github.com/golangee/dyml/parser"
	"github.com/golangee/dyml/token"
)

// DYML language server.
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
					TokenTypes: TokenTypes,
				},
				Full: true,
			},
		},
	}
}

// Initialized tells us, that the client is ready.
func (s *Server) Initialized() {
}

// Handle a hover event.
func (t *Server) Hover(params *protocol.HoverParams) protocol.Hover {
	return protocol.Hover{} // Don't forget to enable hover capabilities when using this.
}

// A document was saved.
func (s *Server) DidSaveTextDocument(params *protocol.DidSaveTextDocumentParams) {
	s.sendDiagnostics()
}

// A document was opened.
func (s *Server) DidOpenTextDocument(params *protocol.DidOpenTextDocumentParams) {
	s.files[params.TextDocument.URI] = File{
		Uri:     params.TextDocument.URI,
		Content: params.TextDocument.Text,
	}
	s.sendDiagnostics()
}

// A document was close.
func (s *Server) DidCloseTextDocument(params *protocol.DidCloseTextDocumentParams) {
	delete(s.files, params.TextDocument.URI)
}

// A document was changed.
func (s *Server) DidChangeTextDocument(params *protocol.DidChangeTextDocumentParams) {
	// There is only a ever single full content change, as we requested.
	s.files[params.TextDocument.URI] = File{
		Uri:     params.TextDocument.URI,
		Content: params.ContentChanges[0].Text,
	}
	s.sendDiagnostics()
}

func (s *Server) FullSemanticTokens(params *protocol.SemanticTokensParams) protocol.SemanticTokens {
	// Mark "let" as a keyword for testing purposes.

	file := s.files[params.TextDocument.URI]

	var data []uint32

	lexer := token.NewLexer(string(file.Uri), strings.NewReader(file.Content))

	// When a comment occurs the lexer emits a comment token and a chardata token.
	// We want to change the type of the chardata to be shown as a comment.
	nextCharIsComment := false

	for {
		tok, err := lexer.Token()
		if err != nil {
			// TODO What should we do here?
			break
		}

		part := SerializeToken(tok, nextCharIsComment)
		nextCharIsComment = false

		switch tok.Type() {
		case token.TokenG1Comment, token.TokenG2Comment:
			nextCharIsComment = true
		}

		data = append(data, part...)
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

func (s *Server) EncodeXML(filename protocol.DocumentURI) string {
	var out strings.Builder
	in := strings.NewReader(s.files[filename].Content)

	enc := encoder.NewXMLEncoder(filepath.Base(string(filename)), in, &out)
	err := enc.Encode()
	if err != nil {
		return ""
	}

	return out.String()
}

// sendDiagnostics sends any parser errors.
func (s *Server) sendDiagnostics() {
	for _, file := range s.files {

		fileContent := file.Content
		fileContent = strings.ToLower(fileContent)
		fileName := filepath.Base(string(file.Uri))

		// Parse file for any errors. Ideally we would be able to catch multiple errors and then recover.
		// Currently only the first error will be reported.
		diagnostics := []protocol.Diagnostic{}
		parser := parser.NewParser(fileName, strings.NewReader(fileContent))
		_, err := parser.Parse()
		if err != nil {
			switch e := err.(type) {
			case *token.PosError:
				for _, detail := range e.Details {
					diagnostics = append(diagnostics, protocol.Diagnostic{
						Range: protocol.Range{
							Start: protocol.Position{
								// Subtract 1 since dyml has 1 based lines and columns, but LSP wants 0 based
								Line:      uint32(detail.Node.Begin().Line) - 1,
								Character: uint32(detail.Node.Begin().Col) - 1,
							},
							End: protocol.Position{
								Line:      uint32(detail.Node.End().Line) - 1,
								Character: uint32(detail.Node.End().Col) - 1,
							},
						},
						Severity: protocol.SeverityError,
						Message:  e.Error(),
					})
				}
			default:
				diagnostics = append(diagnostics, protocol.Diagnostic{
					Severity: protocol.SeverityError,
					Message:  e.Error(),
				})
			}
		}

		_ = SendNotification("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
			URI:         protocol.DocumentURI(file.Uri),
			Diagnostics: diagnostics,
		})
	}
}
