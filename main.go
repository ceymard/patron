package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var st fmt.Stringer

func main() {

	for i, arg := range os.Args[1:] {
		file, err := os.Open(arg)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		lexer := NewLexer(file)

		// for _, tk := range lexer.Tokens {
		// 	os.Stdout.WriteString(tk.Type.DebugString() + " ")
		// 	pp.Print(string(tk.Content))
		// 	os.Stdout.WriteString(" " + strconv.Itoa(tk.StartLine) + "\n")
		// }

		outfile_path := strings.Replace(arg, path.Ext(arg), ".go", 1)
		outfile, err := os.Create(outfile_path)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer outfile.Close()

		package_name := path.Base(path.Dir(arg))
		outfile.WriteString("package " + package_name + "\n\n")
		outfile.WriteString("import (\n")
		outfile.WriteString(`  "io"
  "fmt"
	"bytes"
	"testing"
)

`)

		GenerateGoCode(lexer.Tokens, outfile)
		outfile.WriteString("\n\nfunc TestInclude" + strconv.Itoa(i) + "(t *testing.T) {\n")
		outfile.WriteString("  var buf bytes.Buffer\n" +
			"  var st fmt.Stringer = &buf\n" +
			"  var w io.Writer = &buf\n" +
			"  _, _ = w.Write([]byte(st.String()))\n" +
			"}\n")

		log.Printf("Generated %s", outfile_path)
	}

}
