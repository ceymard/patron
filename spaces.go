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
	current_indent := 0
	indent_stack := make([]int, 0)
	i := 0
	// indent_stack := make([]int, 0)

	for i < len(l.Tokens) {
		prev = tk
		tk = l.Tokens[i]

		switch tk.Type {
		case TokenTypeSpace:
			if prev != nil && prev.Type == TokenTypeNewline && current_indent > 0 {
				tk.FixIndent = current_indent
				// tk.Content = tk.Content[min(current_indent, len(tk.Content)):]
			}

		case TokenTypeControl:
			indent_stack = append(indent_stack, current_indent)

			// line_started := false
			// if tk.StartsLine() {
			// 	line_started = true
			// }

			newline_eaten := false
			newline_pos := -1
			// At the end of a control construct, right after {, eat up spaces a newline if it exists

			for j := i + 1; j < len(l.Tokens) && (l.Tokens[j].Type == TokenTypeSpace || l.Tokens[j].Type == TokenTypeNewline); j++ {
				l.Tokens[j].Skip = true
				if l.Tokens[j].Type == TokenTypeNewline {
					newline_eaten = true
					newline_pos = j
					break
				}
			}

			if newline_eaten {
				// Content that follows the control construct will have to be put at the same indentation level the beginning of the line is
				own_indent := 0

				for j := i - 1; j >= 0; j-- {
					t := l.Tokens[j]
					if t.Type == TokenTypeSpace {
						// We will keep the last space
						own_indent = len(t.Content)
					} else if t.Type == TokenTypeNewline {
						break
					}
				}

				if prev.Type == TokenTypeSpace {
					prev.Skip = true
				}

				next_indent := -1
				for j := newline_pos + 1; j < len(l.Tokens); j++ {
					t := l.Tokens[j]
					if t.Type == TokenTypeSpace && l.Tokens[j-1].Type == TokenTypeNewline {
						if next_indent == -1 {
							next_indent = len(t.Content)
						} else {
							next_indent = min(len(t.Content), next_indent)
						}
					} else if t.Type != TokenTypeNewline && t.Type != TokenTypeSpace {
						break
					}
				}

				// pp.Println(current_indent, own_indent, next_indent, current_indent)
				current_indent = next_indent - max(own_indent-current_indent, 0)
				//
			}

		case TokenTypeEnd, TokenTypeCode:

			if tk.Type == TokenTypeEnd && len(indent_stack) > 0 {
				current_indent = indent_stack[len(indent_stack)-1]
				indent_stack = indent_stack[:len(indent_stack)-1]
			}

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
