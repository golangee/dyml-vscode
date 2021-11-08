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
// charIsComment can be set to true to set the type of CharData to comment.
func SerializeToken(tok token.Token, charIsComment bool) []uint32 {
	// The resulting serialized form we will build in this method.
	var data []uint32

	// Some tokens might span multiple lines and need to be serialized per line.
	// This list contains tokens per line.
	var toks []token.Token
	switch t := tok.(type) {
	case *token.CharData:
		for _, cd := range t.SplitLines() {
			toks = append(toks, cd)
		}
	default:
		toks = append(toks, tok)
	}

	for _, tokPart := range toks {
		// Collect data for this token here and append it to data later.
		tokPartData := make([]uint32, 5)

		// token package handles tokens with 1-based positions, we want 0-based.
		tokPartData[0] = uint32(tokPart.Pos().BeginPos.Line - 1)
		tokPartData[1] = uint32(tokPart.Pos().BeginPos.Col - 1)
		tokPartData[2] = uint32(tokPart.Pos().End().Offset - tokPart.Pos().Begin().Offset)

		switch tokPart.Type() {
		case token.TokenIdentifier:
			tokPartData[3] = TokenKeyword
		case token.TokenCharData:
			tokPartData[3] = TokenString
		case token.TokenG1Comment, token.TokenG2Comment:
			tokPartData[3] = TokenComment
		case token.TokenDefineElement:
			tokPartData[3] = TokenType
		default:
			tokPartData[3] = TokenType
		}

		if tokPart.Type() == token.TokenCharData && charIsComment {
			tokPartData[3] = TokenComment
		}

		data = append(data, tokPartData...)
	}

	return data
}
