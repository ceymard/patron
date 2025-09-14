# patron

A simple templating language that generates fast Go code that priorises readability and ease of inclusion in your code.

# Why

Because most template engines want to have `{{` surrounding tags `}}` or have funky whitespace control adding characters to an already visually heavy syntax.

And because I have never seen any try to tackle having just `@thatkind.Of("Syntax")` with no end tag, which I find more appealing when escaping variables.

# Installation

go install github.com/ceymard/patron

# Usage

```
patron <path-to-pat-file>
```

Will generate a `.go` file in the same directory. The package will be the directory name.

The generated code does not depend on this package. It is up to you to implement Writer pooling and caching. (https://github.com/valyala/bytebufferpool is a good choice.)

# Example

See the `examples/` folder and run them with go test.

# Vscode extension



# Syntax

- Declare functions with `@func <function_name>() {` `}`. Either give the first argument an `io.Writer` or have the function return `string`. A corresponding function ending with `String(` or `Stream(` will be generated.
- Use `@for`, `@if`, `@else`, `@elseif`, `@switch`, `@case`, `@default` and `@continue` to control the flow as you would in go.
- Output go code directly between `@{` and `}`
- Commenting is achieved by using code blocks and putting regular go comments in them. These comments will be part of the go file. `@{/* Comments are done like so */}`.

# Displaying values of types other than string

By default, it is assumed that the displayed values will be of type string. When using other types, use percent modifiers as suffix to variable like they would be used by [fmt](https://pkg.go.dev/fmt).

```
@my.Variable%i will format it as an int.
@(my.Variable)%i will do the same.
```

# Control constructs

- `@{` go code that will be put right there `}`
- `@func functionName(w io.Writer, ...arguments) {` ... `}`
- or `@func functionNameString(arguments...) string {` ... `}`
- `@if condition {` ... `} @elseif condition2 {` ... `} @else {` ... `}`
- `@for i, v := range somevariable {` ... `}`
- `@switch variable { @case value {`... `} @default {` ... `}}`

# Whitespace control

Whitespace control is achieved through placement of the tags in the template.

String output with `@()` or `@inline.Variables` or `@"inline strings"` do not touch whitespace around it.

- Spaces before `}` are removed, up to but not including `\n`
- Spaces after `{` are removed, up to and including one `\n`
- If a `\n` is eaten right after a `{`, then indentation measured and applied to subsequent lines.

In indentation, the *amount* of space characters is the one being considered. Mixing tabs and space for indentation will produce inconsistent results.

If you with to explicitely keep some whitespace around, use @"" or any variation of inline strings to play with the way it works.