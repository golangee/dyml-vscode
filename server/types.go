package dyml

import (
	"dyml-support/protocol"

	"github.com/golangee/dyml/token"
)

// TokenTypes is a list of our supported types from the LSP spec.
// This array is sent once to the editor and after that only integers are used to refer
// to this array.
var TokenTypes = []string{"type", "string", "comment", "keyword"}

// These are indices into the TokenTypes array.
const (
	TokenType = iota
	TokenString
	TokenComment
	TokenKeyword
)

// File is a file that is located at an Uri and has Content.
type File struct {
	Uri     protocol.DocumentURI
	Content string
}

// See https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_semanticTokens
// for an explanation of how this array is built.
// In short: every 5 elements form a tuple (line, col, length, type, modifiers),
// where line, col are relative and type being an index into the array we
// sent to the client in initialize.
// Here the token positions are absolute, they will need to be made relative later.
func SerializeToken(tok token.Token) []uint32 {
	data := make([]uint32, 5)

	// token Package handles tokens with 1-based positions, we want 0-based.
	data[0] = uint32(tok.Pos().BeginPos.Line - 1)
	data[1] = uint32(tok.Pos().BeginPos.Col - 1)
	data[2] = uint32(tok.Pos().End().Offset - tok.Pos().Begin().Offset)

	switch tok.Type() {
	case token.TokenIdentifier:
		data[3] = TokenKeyword
	case token.TokenCharData:
		data[3] = TokenString
	case token.TokenG1Comment, token.TokenG2Comment:
		data[3] = TokenComment
	case token.TokenDefineElement:
		data[3] = TokenType
	default:
		data[3] = TokenType
	}

	return data
}
