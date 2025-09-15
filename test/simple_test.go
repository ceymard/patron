package test

import (
  "io"
  "fmt"
	"bytes"
	"testing"
)


// A regular golang code block.
type Locals struct {
  Field string
  SliceOfInts []int
}

/* Single line comment */

func ShowStuff(w io.Writer, locals Locals) {
  w.Write([]byte("Display field : "))
  w.Write([]byte(locals.Field))
  w.Write([]byte("\n\n"))
  /* ShowStuff(w, locals) */
  w.Write([]byte("\n\n"))
  /*
    Multiline comment.
  */
  w.Write([]byte("\n\nSlice of ints: "))

  for i, v := range locals.SliceOfInts {
    w.Write([]byte("  "))

    if i > 0 {
      w.Write([]byte(""))
      w.Write([]byte(", "))
      w.Write([]byte(""))
    }
    w.Write([]byte(""))
    w.Write([]byte(fmt.Sprintf("%d", v)))
    w.Write([]byte("\n"))
  }
  w.Write([]byte(""))
}

func checkIndent(w io.Writer, ints []int) {
  w.Write([]byte(""))

  if true {
    w.Write([]byte("no-space\n\n"))
  }
  w.Write([]byte("same-indent\n  indented\n\n"))

  for _, v := range ints {
    w.Write([]byte(""))
    w.Write([]byte(fmt.Sprintf("%d", v)))
    w.Write([]byte("\n"))
  }
  w.Write([]byte("\n"))

  for i, v := range ints {
    w.Write([]byte(""))

    if i > 0 {
      w.Write([]byte(""))
      w.Write([]byte(", "))
      w.Write([]byte(""))
    }
    w.Write([]byte(""))
    w.Write([]byte(fmt.Sprintf("%d", v)))
    w.Write([]byte(""))
  }
  w.Write([]byte("\n\n"))
}

func TestIndent(t *testing.T) {
  var buf bytes.Buffer
  checkIndent(&buf, []int{1, 2, 3})
  expect(t, "no-space\n\nsame-indent\n  indented\n\n1\n2\n3\n\n1, 2, 3\n\n", buf.String())
}



func TestFmt(t *testing.T) {
  expect(t, "1.1, 2.2, 3.0\n", checkFmt([]float64{1.1, 2.22, 3.0}))
}


func checkFmt(flt []float64) string {
  var ø bytes.Buffer
  ø.Write([]byte(""))

  for i, v := range flt {
    ø.Write([]byte(""))

    if i > 0 {
      ø.Write([]byte(""))
      ø.Write([]byte(", "))
      ø.Write([]byte(""))
    }
    ø.Write([]byte(""))
    ø.Write([]byte(fmt.Sprintf("%.1f", v)))
    ø.Write([]byte(""))
  }
  ø.Write([]byte("\n"))
  return ø.String()
}

//
func TestInlineTagOpener(t *testing.T) {
  expect(t, "check stuff\n", checkInlineTagOpener())
}


func checkInlineTagOpener() string {
  var ø bytes.Buffer
  ø.Write([]byte("check "))

  if true {
    ø.Write([]byte("stuff"))
  }
  ø.Write([]byte("\n"))
  return ø.String()
}

func checkTagOpener() string {
  var ø bytes.Buffer
  ø.Write([]byte("check "))

  if true {
    ø.Write([]byte("stuff"))
  }
  ø.Write([]byte("\n"))

  if true {
    ø.Write([]byte("more stuff\n"))
  }
  ø.Write([]byte(""))
  return ø.String()
}

func TestTagOpener(t *testing.T) {
  expect(t, "check stuff\nmore stuff\n", checkTagOpener())
}


func TestCodeOutput(t *testing.T) {
  expect(t, "hello world how's it going?\n", checkCodeOutput("world"))
}


func checkCodeOutput(str string) string {
  var ø bytes.Buffer
  ø.Write([]byte("hello "))
  ø.Write([]byte(str))
  ø.Write([]byte(" how's it going?\n"))
  return ø.String()
}

func TestIndent2(t *testing.T) {
  expect(t, "in\n  dent\nin\n  dent\n", checkIndent2())
}


func checkIndent2() string {
  var ø bytes.Buffer
  ø.Write([]byte(""))

  if true {
    ø.Write([]byte("in\n  dent\n"))
  }
  ø.Write([]byte("in\n  dent\n"))
  return ø.String()
}

func TestIndent3(t *testing.T) {
    expect(t, "in\n  dent\nin\n  dent\n  dent\n", checkIndent3())
  }


func checkIndent3() string {
  var ø bytes.Buffer
  ø.Write([]byte(""))

  if true {
    ø.Write([]byte("in\n  dent\n"))
  }
  ø.Write([]byte("in\n"))

  if true {
    ø.Write([]byte("  "))

    if true {
      ø.Write([]byte("d"))
    }
    ø.Write([]byte("ent\n  "))

    if true {
      ø.Write([]byte("d"))
    }
    ø.Write([]byte("ent\n"))
  }
  ø.Write([]byte(""))
  return ø.String()
}


func TestInclude0(t *testing.T) {
  var buf bytes.Buffer
  var st fmt.Stringer = &buf
  st.String()
}
