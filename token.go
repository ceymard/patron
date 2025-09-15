package main

import (
	"unicode"
)

type TokenType int

const (
	TokenTypeText       TokenType = iota // Text content
	TokenTypeCode                        // 2
	TokenTypeOutputCode                  // 3
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
	Content     []rune
	Pos         int // position inside the input
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Fmt         string
}

func (t *Token) isText() bool {
	return t.Type == TokenTypeText
}

func (t *Token) trimEnd() {
	pos := len(t.Content)
	for pos > 0 {
		c := t.Content[pos-1]
		if unicode.IsSpace(c) && c != '\n' {
			pos--
		} else {
			break
		}
	}

	if pos >= 0 && pos < len(t.Content) {
		t.Content = t.Content[:pos]
	}
}

func (t *Token) trimLeadingEmptyLines() {
	last_line_start := -1
	i := 0
	for i < len(t.Content) {
		c := t.Content[i]
		if c == '\n' {
			last_line_start = i
			break
		}
		if !unicode.IsSpace(c) {
			break
		}
		i++
	}
	if last_line_start > -1 {
		t.Content = t.Content[last_line_start+1:]
	} else if i > 0 {
		t.Content = t.Content[i:]
	}
}

// Measure the indentation level of the token
func (t *Token) measureIndent() int {
	indent := -1
	i := 0
	for i < len(t.Content) {
		c := t.Content[i]
		if c == '\n' {
			indent = 0
		} else if unicode.IsSpace(c) {
			if indent >= 0 {
				indent++
			}
		} else {
			break
		}
		i++
	}
	return indent
}

func (t *Token) applyIndent(indent int, start_of_line bool) {
	if len(t.Content) == 0 {
		return
	}

	var result []rune
	i := 0
	at_line_start := start_of_line

	for i < len(t.Content) {
		c := t.Content[i]

		if c == '\n' {
			result = append(result, c)
			at_line_start = true
			i++
		} else if at_line_start && unicode.IsSpace(c) {
			// Count consecutive spaces at the start of a line
			space_count := 0
			j := i
			for j < len(t.Content) && unicode.IsSpace(t.Content[j]) {
				space_count++
				j++
			}

			// Remove up to 'indent' spaces, but don't remove more than available
			spaces_to_remove := indent
			if spaces_to_remove > space_count {
				spaces_to_remove = space_count
			}

			// Add remaining spaces to result
			for k := 0; k < space_count-spaces_to_remove; k++ {
				result = append(result, t.Content[i+k])
			}

			i = j
			at_line_start = false
		} else {
			result = append(result, c)
			at_line_start = false
			i++
		}
	}

	t.Content = result
}
