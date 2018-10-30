# easyjson [![Go Report Card](https://goreportcard.com/badge/github.com/stpetejess/easyjson)](https://goreportcard.com/report/github.com/stpetejess/easyjson)

This was forked from https://github.com/mailru/easyjson for a personal project
where I needed more control to specify custom (and multiple custom) tags to use for generation.

Package easyjson provides a fast and easy way to marshal/unmarshal Go structs
to/from JSON without the use of reflection.

easyjson aims to keep generated Go code simple enough so that it can be easily
optimized or fixed. Another goal is to provide users with the ability to
customize the generated code by providing options not available with the
standard `encoding/json` package, such as generating "snake_case" names or
enabling `omitempty` behavior by default.

## Usage
```sh
# install
go get -u github.com/stpetejess/easyjson/...

# run
easyjson -all <file>.go
```

The above will generate `<file>_easyjson.go` containing the appropriate marshaler and
unmarshaler funcs for all structs contained in `<file>.go`.

Please note that easyjson requires a full Go build environment and the `GOPATH`
environment variable to be set. This is because easyjson code generation
invokes `go run` on a temporary file (an approach to code generation borrowed
from [ffjson](https://github.com/pquerna/ffjson)).

## Options
```txt
Usage of easyjson:
  -all
    	generate marshaler/unmarshalers for all structs in a file
  -build_tags string
    	build tags to add to generated file
  -leave_temps
    	do not delete temporary files
  -no_std_marshalers
    	don't generate MarshalJSON/UnmarshalJSON funcs
  -noformat
    	do not run 'gofmt -w' on output file
  -omit_empty
    	omit empty fields by default
  -output_filename string
    	specify the filename of the output
  -pkg
    	process the whole package instead of just the given file
  -snake_case
    	use snake_case names instead of CamelCase by default
  -lower_camel_case
        use lowerCamelCase instead of CamelCase by default
  -stubs
    	only generate stubs for marshaler/unmarshaler funcs
  -tags
      generate functions for marshaling and unmarshaling based on custom tags. Comma delimited list.
      to specify a default use a colon eg: elastic:json will generate marshal and unmarshal fn for
      tag elastic and, if missing, will default to tag json. If no default specified it
      uses camel case of the attribute name. If - is specified as the default it will
      ignore (de)/serialization for the specified attribute.
```

Using `-tags` will generate functions for marashaling and unmarshaling based on custom tags.
For example:
```easyjson -all -pkg -no_std_marshalers -tags=elastic:json,bson,render:- ./models```
will generate marshal and unmarshal fn for tag elastic (default to json) bson (no default)
and render (default to ignore if no render tag) for all structs found in the package models

Using `-all` will generate marshalers/unmarshalers for all Go structs in the
file. If `-all` is not provided, then only those structs whose preceding
comment starts with `easyjson:json` will have marshalers/unmarshalers
generated. For example:

```go
//easyjson:json
type A struct {}
```

Additional option notes:

* `-snake_case` tells easyjson to generate snake\_case field names by default
  (unless overridden by a field tag). The CamelCase to snake\_case conversion
  algorithm should work in most cases (ie, HTTPVersion will be converted to
  "http_version").

* `-build_tags` will add the specified build tags to generated Go sources.

## Generated Marshaler/Unmarshaler Funcs

For Go struct types, easyjson generates the funcs `MarshalEasyJSON` /
`UnmarshalEasyJSON` for marshaling/unmarshaling JSON. In turn, these satisify
the `easyjson.Marshaler` and `easyjson.Unmarshaler` interfaces and when used in
conjunction with `easyjson.Marshal` / `easyjson.Unmarshal` avoid unnecessary
reflection / type assertions during marshaling/unmarshaling to/from JSON for Go
structs.

easyjson also generates `MarshalJSON` and `UnmarshalJSON` funcs for Go struct
types compatible with the standard `json.Marshaler` and `json.Unmarshaler`
interfaces. Please be aware that using the standard `json.Marshal` /
`json.Unmarshal` for marshaling/unmarshaling will incur a significant
performance penalty when compared to using `easyjson.Marshal` /
`easyjson.Unmarshal`.

Additionally, easyjson exposes utility funcs that use the `MarshalEasyJSON` and
`UnmarshalEasyJSON` for marshaling/unmarshaling to and from standard readers
and writers. For example, easyjson provides `easyjson.MarshalToHTTPResponseWriter`
which marshals to the standard `http.ResponseWriter`. Please see the [GoDoc
listing](https://godoc.org/github.com/stpetejess/easyjson) for the full listing of
utility funcs that are available.

## Controlling easyjson Marshaling and Unmarshaling Behavior

Go types can provide their own `MarshalEasyJSON` and `UnmarshalEasyJSON` funcs
that satisify the `easyjson.Marshaler` / `easyjson.Unmarshaler` interfaces.
These will be used by `easyjson.Marshal` and `easyjson.Unmarshal` when defined
for a Go type.

Go types can also satisify the `easyjson.Optional` interface, which allows the
type to define its own `omitempty` logic.

## Type Wrappers

easyjson provides additional type wrappers defined in the `easyjson/opt`
package. These wrap the standard Go primitives and in turn satisify the
easyjson interfaces.

The `easyjson/opt` type wrappers are useful when needing to distinguish between
a missing value and/or when needing to specifying a default value. Type
wrappers allow easyjson to avoid additional pointers and heap allocations and
can significantly increase performance when used properly.

## Memory Pooling

easyjson uses a buffer pool that allocates data in increasing chunks from 128
to 32768 bytes. Chunks of 512 bytes and larger will be reused with the help of
`sync.Pool`. The maximum size of a chunk is bounded to reduce redundant memory
allocation and to allow larger reusable buffers.

easyjson's custom allocation buffer pool is defined in the `easyjson/buffer`
package, and the default behavior pool behavior can be modified (if necessary)
through a call to `buffer.Init()` prior to any marshaling or unmarshaling.
Please see the [GoDoc listing](https://godoc.org/github.com/stpetejess/easyjson/buffer)
for more information.

## Issues, Notes, and Limitations

* easyjson is still early in its development. As such, there are likely to be
  bugs and missing features when compared to `encoding/json`. In the case of a
  missing feature or bug, please create a GitHub issue. Pull requests are
  welcome!

* Unlike `encoding/json`, object keys are case-sensitive. Case-insensitive
  matching is not currently provided due to the significant performance hit
  when doing case-insensitive key matching. In the future, case-insensitive
  object key matching may be provided via an option to the generator.

* easyjson makes use of `unsafe`, which simplifies the code and
  provides significant performance benefits by allowing no-copy
  conversion from `[]byte` to `string`. That said, `unsafe` is used
  only when unmarshaling and parsing JSON, and any `unsafe` operations
  / memory allocations done will be safely deallocated by
  easyjson. Set the build tag `easyjson_nounsafe` to compile it
  without `unsafe`.

* easyjson is compatible with Google App Engine. The `appengine` build
  tag (set by App Engine's environment) will automatically disable the
  use of `unsafe`, which is not allowed in App Engine's Standard
  Environment. Note that the use with App Engine is still experimental.

* Floats are formatted using the default precision from Go's `strconv` package.
  As such, easyjson will not correctly handle high precision floats when
  marshaling/unmarshaling JSON. Note, however, that there are very few/limited
  uses where this behavior is not sufficient for general use. That said, a
  different package may be needed if precise marshaling/unmarshaling of high
  precision floats to/from JSON is required.

* While unmarshaling, the JSON parser does the minimal amount of work needed to
  skip over unmatching parens, and as such full validation is not done for the
  entire JSON value being unmarshaled/parsed.

* Currently there is no true streaming support for encoding/decoding as
  typically for many uses/protocols the final, marshaled length of the JSON
  needs to be known prior to sending the data. Currently this is not possible
  with easyjson's architecture.
