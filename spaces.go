package main

import (
	"slices"
)

// Handle spaces
func (l *Lexer) CollapseText() {

	i := 0
	for i < len(l.Tokens)-1 {
		tk := l.Tokens[i]
		next := l.Tokens[i+1]
		if tk.isText() && len(tk.Content) == 0 {
			l.Tokens = slices.Delete(l.Tokens, i, i+1)
		} else if tk.isText() && next.isText() {
			tk.Content = slices.Concat(tk.Content, next.Content)
			l.Tokens = slices.Delete(l.Tokens, i+1, i+2)
		} else {
			i++
		}
	}

}

/* HandleSpace removes space according to spacing rules */
func (l *Lexer) CollapseSpaces() int {
	var prev *Token = nil
	// var next *Token = nil
	var tk *Token = nil

	// indent := 0
	i := 0
	// indent_stack := make([]int, 0)

	for i < len(l.Tokens) {
		prev = tk
		tk = l.Tokens[i]

		switch tk.Type {
		case TokenTypeText:

		case TokenTypeControl:

			line_started := false
			if tk.StartsLine() {
				line_started = true
			}

			newline_eaten := false
			// At the end of a control construct, right after {, eat up spaces a newline if it exists

			for j := i + 1; j < len(l.Tokens) && (l.Tokens[j].Type == TokenTypeSpace || l.Tokens[j].Type == TokenTypeNewline); j++ {
				l.Tokens[j].Skip = true
				if l.Tokens[j].Type == TokenTypeNewline {
					newline_eaten = true
					break
				}
			}

			if line_started && newline_eaten {
				// own_indent := 0
				if prev.Type == TokenTypeSpace {
					prev.Skip = true
					// own_indent = len(prev.Content)
				}

				//
			}

		case TokenTypeEnd:

			// Remove spaces before the end token
			if prev != nil && prev.Type == TokenTypeSpace {
				prev.Skip = true
			}

			is_alone := true
			for j := i - 1; j >= 0; j-- {
				t := l.Tokens[j]
				if t.Type == TokenTypeNewline {
					break
				}
				if t.Type != TokenTypeSpace && t.Type != TokenTypeEnd {
					is_alone = false
					break
				}
			}

			for j := i + 1; j < len(l.Tokens) && (l.Tokens[j].Type == TokenTypeSpace || is_alone && l.Tokens[j].Type == TokenTypeNewline); j++ {
				l.Tokens[j].Skip = true
				if l.Tokens[j].Type == TokenTypeNewline {
					break
				}
			}
		}

		i++
	}

	// for i := 0; i < len(l.Tokens); i++ {
	// 	pp.Println(string(l.Tokens[i].Content), l.Tokens[i].Skip)
	// }

	return i
}
