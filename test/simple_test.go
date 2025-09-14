package test

import (
  "io"
  "fmt"
	"bytes"
)


import (
  "testing"
)

// A regular golang code block.
type Locals struct {
  Field string
  SliceOfInts []int
}

/* Single line comment */

func ShowStuff(w io.Writer, locals Locals) {
  w.Write([]byte(`Display field : `))
  w.Write([]byte(locals.Field))
  w.Write([]byte(`

`))
  /* ShowStuff(w, locals) */
  w.Write([]byte(`
`))
  /*
    Multiline comment.
  */
  w.Write([]byte(`
Slice of ints: `))

  for i, v := range locals.SliceOfInts {

    if i > 0 {
      w.Write([]byte(", "))
    }
    w.Write([]byte(fmt.Sprintf("%d", v)))
    w.Write([]byte(`
`))
  }
}

func checkIndent(w io.Writer, ints []int) {

  if true {
    w.Write([]byte(`no-space

`))
  }
  w.Write([]byte(`same-indent
indented


`))

  for _, v := range ints {
    w.Write([]byte(fmt.Sprintf("%d", v)))
    w.Write([]byte(`
`))
  }
  w.Write([]byte(`

`))

  for i, v := range ints {

    if i > 0 {
      w.Write([]byte(", "))
    }
    w.Write([]byte(fmt.Sprintf("%d", v)))
  }
  w.Write([]byte(`

`))
}

func TestIndent(t *testing.T) {
  var buf bytes.Buffer
  checkIndent(&buf, []int{1, 2, 3})
  expect(t, `no-space

same-indent
indented


1
2
3


1, 2, 3

`, buf.String())
}



func TestFmt(t *testing.T) {
  expect(t, "1.1, 2.2, 3.0\n", checkFmt([]float64{1.1, 2.22, 3.0}))
}


func checkFmt(flt []float64) string {
  var ø bytes.Buffer

  for i, v := range flt {

    if i > 0 {
      ø.Write([]byte(", "))
    }
    ø.Write([]byte(fmt.Sprintf("%.1f", v)))
  }
  ø.Write([]byte(`
`))
  return ø.String()
}

//
func TestInlineTagOpener(t *testing.T) {
  expect(t, "check stuff\n", checkInlineTagOpener())
}


func checkInlineTagOpener() string {
  var ø bytes.Buffer
  ø.Write([]byte(`check `))

  if true {
    ø.Write([]byte(`stuff`))
  }
  ø.Write([]byte(`
`))
  return ø.String()
}

func checkTagOpener() string {
  var ø bytes.Buffer
  ø.Write([]byte(`check `))

  if true {
    ø.Write([]byte(`stuff`))
  }
  ø.Write([]byte(`
`))

  if true {
    ø.Write([]byte(`more stuff
`))
  }
  return ø.String()
}

func TestTagOpener(t *testing.T) {
  expect(t, "check stuff\nmore stuff\n", checkTagOpener())
}



func TestInclude0(t *testing.T) {
  var buf bytes.Buffer
  var st fmt.Stringer = &buf
  st.String()
}
