@{
  import (
    "testing"
  )

  // A regular golang code block.
  type Locals struct {
    Field string
    SliceOfInts []int
  }
}

@{/* Single line comment */}
@func ShowStuff(w io.Writer, locals Locals) {
  Display field : @locals.Field

  @{/* ShowStuff(w, locals) */}

  @{/*
    Multiline comment.
  */}

  Slice of ints: @for i, v := range locals.SliceOfInts {
    @if i > 0 { @", " } @v%d
  }
}

@func checkIndent(w io.Writer, ints []int) {
  @if true {
    no-space

  }
  same-indent
    indented


  @for _, v := range ints {
    @v%d
  }


  @for i, v := range ints { @if i > 0 { @", " } @v%d }

}

@{
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

}

@{
  func TestFmt(t *testing.T) {
    expect(t, "1.1, 2.2, 3.0\n", checkFmt([]float64{1.1, 2.22, 3.0}))
  }
}
@func checkFmt(flt []float64) string {
  @for i, v := range flt { @if i > 0 { @", " } @v%.1f }
}

@{
  //
  func TestInlineTagOpener(t *testing.T) {
    expect(t, "check stuff\n", checkInlineTagOpener())
  }
}
@func checkInlineTagOpener() string {
  check @if true { stuff }
}

@func checkTagOpener() string {
  check @if true { stuff }
  @if true {
    more stuff
  }
}

@{
  func TestTagOpener(t *testing.T) {
    expect(t, "check stuff\nmore stuff\n", checkTagOpener())
  }
}