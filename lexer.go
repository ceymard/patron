package main

import (
	"bytes"
	"io"
	"unicode"
)

type Lexer struct {
	Input  []rune
	Pos    int
	Line   int
	Column int
	Tokens []*Token
}

func NewLexer(input io.Reader) *Lexer {
	b, _ := io.ReadAll(input)

	l := &Lexer{
		Input:  bytes.Runes(b),
		Pos:    0,
		Line:   1,
		Column: 1,
		Tokens: make([]*Token, 0),
	}

	l.Lex(false)

	l.CollapseSpaces()

	return l
}

func (l *Lexer) Lex(waiting_for_bracket bool) {
	for l.Pos < len(l.Input) {
		c := l.Input[l.Pos]

		if c == '\n' {
			l.NewToken(TokenTypeNewline, l.Pos, l.Pos+1)
		} else if unicode.IsSpace(c) {
			start := l.Pos
			l.Pos++
			for l.Pos < len(l.Input) && unicode.IsSpace(l.Input[l.Pos]) {
				l.Pos++
			}
			l.NewToken(TokenTypeSpace, start, l.Pos)
		} else if c == '@' && l.Pos+1 < len(l.Input) {
			l.lexAt()
		} else if waiting_for_bracket && c == '}' {
			l.NewToken(TokenTypeEnd, l.Pos, l.Pos+1)
			return
		} else {
			l.lexText(waiting_for_bracket)
			// No space, no special character, it's text then
		}
	}
}

// lexText lexes text until it encounters a '@' or the end of the input. If it encounters a '@', it will backtrack to the last non-space character.
// It considers that the current character is part of the text.
func (l *Lexer) lexText(waiting_for_bracket bool) {
	start := l.Pos
	pos := l.Pos

	for pos < len(l.Input) {

		if unicode.IsSpace(l.Input[pos]) || l.Input[pos] == '@' || waiting_for_bracket && l.Input[pos] == '}' {
			break
		}

		pos++
	}

	l.NewToken(TokenTypeText, start, pos)
}

func (l *Lexer) lexAt() {
	start := l.Pos + 1 // skip the @
	pos := start

	c := l.Input[pos]

	if c == '(' || c == '{' {
		start_is_paren := c == '(' // to differentiate between @(go code) and @{go code}
		kind := TokenTypeCode      // {
		until := '}'
		if start_is_paren {
			kind = TokenTypeOutputCode
			until = ')'
		}

		// @(go code)

		pos = l.advanceGoCodeUntil(pos, until)
		tk := l.NewToken(kind, start, pos)
		tk.Content = Dedent(tk.Content[1 : len(tk.Content)-1])

		if kind == TokenTypeOutputCode {
			l.parseFmtSpecifier(tk)
		}

		// FIXME parse output specifier
		return
	} else if c == '_' || unicode.IsLetter(c) {
		// Start of an identifier
		pos++

		for pos < len(l.Input) {
			if !unicode.IsLetter(l.Input[pos]) && !unicode.IsDigit(l.Input[pos]) && l.Input[pos] != '_' {
				break
			}
			pos++
		}

		identifier := string(l.Input[start:pos])

		switch identifier {
		case "func", "if", "else", "elseif", "for", "switch", "case", "default":
			pos = l.advanceGoCodeUntil(pos, '{')
			l.NewToken(TokenTypeControl, start, pos)
			l.Lex(true)
			// should now be on the closing bracket
			// l.NewToken(TokenTypeEnd, pos, pos+1)
			// l.Pos++
		case "break", "continue", "return":
			l.NewToken(TokenTypeStatement, start, pos)
		default:
			for pos < len(l.Input) {
				c := l.Input[pos]
				if c == '.' {
					pos++
					// Advance on an identifier
					found := false
					for pos < len(l.Input) && (unicode.IsLetter(l.Input[pos]) ||
						unicode.IsDigit(l.Input[pos]) ||
						l.Input[pos] == '_') {
						found = true
						pos++
					}
					if !found {
						pos--
						break
					}
				} else if c == '(' {
					pos = l.advanceGoCodeUntil(pos, ')')
				} else if c == '[' {
					pos = l.advanceGoCodeUntil(pos, ']')
				} else {
					break
				}
			}
			tk := l.NewToken(TokenTypeOutputCode, start, pos)
			l.parseFmtSpecifier(tk)
			// inline expression
		}
	} else if c == '"' || c == '\'' || c == '`' {
		pos++
		for pos < len(l.Input) {
			if l.Input[pos] == '\\' && pos+1 < len(l.Input) {
				pos++
			} else if l.Input[pos] == c {
				pos++
				break
			}
			pos++
		}
		l.NewToken(TokenTypeString, start, pos)
		// tk.Content = []rune(strings.ReplaceAll(string(tk.Content[1:len(tk.Content)-1]), "\\\"", "\""))
	} else {
		l.NewToken(TokenTypeIllegal, start, pos+1)
	}

	// pos := l.Pos
}

func (l *Lexer) advanceGoCodeUntil(pos int, until rune) int {

	balance := 0
	in_quote := ' '

	for pos < len(l.Input) {
		c := l.Input[pos]

		if c == until && balance == 0 {
			pos++
			break
		}

		if c == in_quote {
			in_quote = ' ' // found the matching quote
		} else if c == '(' || c == '{' || c == '[' {
			balance++
		} else if c == ')' || c == '}' || c == ']' {
			balance--
		} else if c == '"' || c == '\'' || c == '`' {
			in_quote = c
		} else if c == '\\' && pos+1 < len(l.Input) {
			pos++
		} else if c == '/' {
			// Single line comment
			if pos+1 < len(l.Input) && l.Input[pos+1] == '/' {
				for pos < len(l.Input) && l.Input[pos] != '\n' {
					pos++
				}
			} else if pos+1 < len(l.Input) && l.Input[pos+1] == '*' {
				// Multiline comment
				for pos < len(l.Input) {
					if l.Input[pos] == '*' && pos+1 < len(l.Input) && l.Input[pos+1] == '/' {
						pos++
						break
					}
					pos++
				}
			}
		}

		if c == until && balance == 0 {
			pos++
			break
		}

		pos++
	}

	return pos
}

func (l *Lexer) parseFmtSpecifier(tk *Token) {
	start := tk.Pos + 1
	pos := tk.Pos + 1
	if pos >= len(l.Input) || l.Input[pos] != '%' {
		return
	}
	pos++
	if pos < len(l.Input) && l.Input[pos] == '#' {
		pos++
	}
	dot_found := false
	// parse float specifier
	for {
		if !dot_found && pos < len(l.Input) && l.Input[pos] == '.' {
			dot_found = true
		} else if !unicode.IsDigit(l.Input[pos]) {
			break
		}
		pos++
	}

	if pos >= len(l.Input) {
		// invalid specifier
		return
	}

	pos++

	spec := string(l.Input[start:pos])
	tk.Fmt = spec
	l.Pos = pos
}

// Only when tokens are actually create do we advance the lexer position
func (l *Lexer) NewToken(t TokenType, start int, end int) *Token {
	startLine := l.Line
	startColumn := l.Column

	for i := start; i < end; i++ {
		if l.Input[i] == '\n' {
			l.Line++
			l.Column = 1
		} else {
			l.Column++
		}

	}

	endLine := l.Line
	endColumn := l.Column

	tk := &Token{
		Lexer:       l,
		Type:        t,
		Content:     l.Input[start:end],
		Pos:         start,
		PosInSlice:  len(l.Tokens),
		StartLine:   startLine,
		StartColumn: startColumn,
		EndLine:     endLine,
		EndColumn:   endColumn,
	}

	l.Tokens = append(l.Tokens, tk)
	l.Pos = end

	return tk
}
