package main

type TokenType int

const (
	TokenTypeText TokenType = iota // Text content
	TokenTypeSpace
	TokenTypeNewline
	TokenTypeCode       // 2
	TokenTypeOutputCode // 3
	TokenTypeControl
	TokenTypeStatement
	TokenTypeEnd
	TokenTypeIllegal
	TokenTypeString
)

var mp = map[TokenType]string{
	TokenTypeText:       "Text",
	TokenTypeCode:       "Code",
	TokenTypeOutputCode: "OutputCode",
	TokenTypeControl:    "Control",
	TokenTypeStatement:  "Statement",
	TokenTypeEnd:        "End",
	TokenTypeIllegal:    "Illegal",
	TokenTypeString:     "String",
}

func (t TokenType) DebugString() string {
	return mp[t]
}

type Token struct {
	Type        TokenType
	Lexer       *Lexer
	Content     []rune
	PosInSlice  int
	Pos         int // position inside the input
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Fmt         string
	Skip        bool
	FixIndent   int
}

func (t *Token) isText() bool {
	return t.Type == TokenTypeText || t.Type == TokenTypeSpace || t.Type == TokenTypeNewline
}

// StartsLine returns true if this token is the first non space token of the line.
func (t *Token) StartsLine() bool {
	for i := t.PosInSlice - 1; i >= 0; i-- {
		prev := t.Lexer.Tokens[i]
		if prev.Type == TokenTypeNewline {
			return true
		}
		if prev.Type == TokenTypeSpace {
			continue
		}
		break
	}
	return false
}

// func (t *Token) EndsLine() bool {

// }
