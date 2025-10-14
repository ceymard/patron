package main

import (
	"io"
	"regexp"
	"strings"
	"unicode"
)

var id = `\p{L}[\p{L}\p{N}_]*`
var re_func = regexp.MustCompile(`^func\s+(` + id + `)\s*\((` + id + `)?\s*([\w.]*).*?\)\s*(string)?\s*\{`)

func GenerateGoCode(tokens []*Token, output io.Writer) int {

	i := 0
	indent_stack := make([]string, 0)
	indent := ""

	writer_var := "ø"
	output_string := false
	balance := -1

	for i < len(tokens) {
		token := tokens[i]
		content := string(token.Content)

		switch token.Type {
		case TokenTypeText, TokenTypeSpace, TokenTypeNewline:
			str := ""
			for i < len(tokens) {
				tk := tokens[i]
				if !tk.isText() {
					break
				}
				if !tk.Skip {
					if tk.FixIndent > 0 {
						str += string(tk.Content[min(tk.FixIndent, len(tk.Content)):])
					} else {
						str += string(tk.Content)
					}
				}
				i++
			}

			if balance >= 0 && str != "" {
				str = strings.ReplaceAll(str, "\\", "\\\\")
				str = strings.ReplaceAll(str, "\"", "\\\"")
				str = strings.ReplaceAll(str, "\n", "\\n")
				output.Write([]byte(indent + writer_var + ".Write([]byte(\"" + str + "\"))\n"))
			}
			continue
		case TokenTypeCode:
			output.Write([]byte(indent + content + "\n"))
		case TokenTypeOutputCode:
			if token.Fmt != "" {
				output.Write([]byte(indent + writer_var + ".Write([]byte(fmt.Sprintf(\"" + token.Fmt + "\", " + content + ")))\n"))
			} else {
				output.Write([]byte(indent + writer_var + ".Write([]byte(" + content + "))\n"))
			}
		case TokenTypeControl:
			balance++

			matches := re_func.FindStringSubmatch(content)
			if len(matches) > 0 {
				// name := matches[1]
				write_var := matches[2]
				write_type := matches[3]
				rettype := matches[4]

				if write_type == "io.Writer" {
					writer_var = write_var
					output_string = false
				} else if rettype == "string" {
					writer_var = "ø"
					output_string = true
				}
			}

			if !strings.HasPrefix(content, "else") {
				output.Write([]byte("\n" + indent))
			} else {
				content = strings.Replace(content, "elseif", "else if", 1)
			}
			output.Write([]byte(content + "\n"))
			indent_stack = append(indent_stack, indent)
			indent = indent + "  "

			if output_string && len(matches) > 0 {
				output.Write([]byte(indent + "var ø bytes.Buffer\n"))
			}

		case TokenTypeEnd:
			balance--
			if len(indent_stack) > 0 {
				indent = indent_stack[len(indent_stack)-1]
				indent_stack = indent_stack[:len(indent_stack)-1]
			} else {
				indent = ""
			}

			if indent == "" && output_string {
				output.Write([]byte(indent + "  return ø.String()\n"))
			}
			output.Write([]byte(indent + "}"))

			if i < len(tokens)-1 {
				next := tokens[i+1]

				else_position := -1
				for j := i + 1; j < len(tokens); j++ {
					tk := tokens[j]
					if tk.Type == TokenTypeControl && strings.HasPrefix(string(tk.Content), "else") {
						else_position = j
						break
					} else if tk.Type != TokenTypeSpace && tk.Type != TokenTypeNewline {
						break
					}
				}

				if next != nil && else_position == -1 {
					output.Write([]byte("\n"))
				} else {
					output.Write([]byte(" "))
					i = else_position
					continue
				}
			}

		case TokenTypeStatement:
			output.Write([]byte(indent + content + "\n"))

		case TokenTypeString:
			output.Write([]byte(indent + writer_var + ".Write([]byte(" + content + "))\n"))
		}

		i++
	}

	return i
}

// Dedent detects indentation from the first non-space after a newline,
// and removes that amount of leading spaces from all lines.
// If the first non-space appears before any newline, the text is unchanged.
func Dedent(runes []rune) []rune {
	// Find position of first non-space rune
	firstNonSpaceIdx := -1
	spaceCount := 0
	lastNewlineIdx := -1

	for i, r := range runes {
		if r == '\n' {
			lastNewlineIdx = i
			spaceCount = 0
			continue
		}
		if unicode.IsSpace(r) && r != '\n' {
			spaceCount++
			continue
		}
		// Found first non-space
		firstNonSpaceIdx = i
		break
	}

	// If nothing found, return as-is
	if firstNonSpaceIdx == -1 {
		return runes
	}

	// If the first non-space is before any newline -> do not modify
	if lastNewlineIdx == -1 || firstNonSpaceIdx < lastNewlineIdx {
		return []rune(strings.TrimSpace(string(runes)))
		// return runes
	}

	indent := spaceCount
	if indent == 0 {
		return runes
	}

	// Walk line by line, remove up to `indent` leading spaces
	var out []rune
	start := 0
	for i := 0; i <= len(runes); i++ {
		if i == len(runes) || runes[i] == '\n' {
			line := runes[start:i]
			// Trim leading spaces
			trim := 0
			j := 0
			for ; j < len(line) && trim < indent; j++ {
				if line[j] == ' ' {
					trim++
				} else {
					break
				}
			}
			out = append(out, line[j:]...)
			if i < len(runes) { // preserve newline
				out = append(out, '\n')
			}
			start = i + 1
		}
	}

	return out
}
