// Package easyjson contains marshaler/unmarshaler interfaces and helper functions.
package easyjson

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/bringhub/easyjson/jlexer"
	"github.com/bringhub/easyjson/jwriter"
)

// Marshaler is an easyjson-compatible marshaler interface.
type Marshaler interface {
	MarshalEasyJSON(w *jwriter.Writer, usingTagName string)
}

// Marshaler is an easyjson-compatible unmarshaler interface.
type Unmarshaler interface {
	UnmarshalEasyJSON(w *jlexer.Lexer, usingTagName string)
}

// Optional defines an undefined-test method for a type to integrate with 'omitempty' logic.
type Optional interface {
	IsDefined() bool
}

// Marshal returns data as a single byte slice. Method is suboptimal as the data is likely to be copied
// from a chain of smaller chunks.
func Marshal(v Marshaler) ([]byte, error) {
	return MarshalCustom(v, "json")
}
func MarshalCustom(v Marshaler, usingTagName string) ([]byte, error) {
	w := jwriter.Writer{}
	w.Flags = jwriter.NilSliceAsEmpty
	v.MarshalEasyJSON(&w, usingTagName)
	return w.BuildBytes()
}

// MarshalToWriter marshals the data to an io.Writer.
func MarshalToWriter(v Marshaler, w io.Writer) (written int, err error) {
	return MarshalToWriterCustom(v, w, "json")
}
func MarshalToWriterCustom(v Marshaler, w io.Writer, usingTagName string) (written int, err error) {
	jw := jwriter.Writer{}
	v.MarshalEasyJSON(&jw, usingTagName)
	return jw.DumpTo(w)
}

// MarshalToHTTPResponseWriter sets Content-Length and Content-Type headers for the
// http.ResponseWriter, and send the data to the writer. started will be equal to
// false if an error occurred before any http.ResponseWriter methods were actually
// invoked (in this case a 500 reply is possible).
func MarshalToHTTPResponseWriter(v Marshaler, w http.ResponseWriter) (started bool, written int, err error) {
	return MarshalToHTTPResponseWriterCustom(v, w, "json")
}
func MarshalToHTTPResponseWriterCustom(v Marshaler, w http.ResponseWriter, usingTagName string) (started bool, written int, err error) {
	jw := jwriter.Writer{}
	v.MarshalEasyJSON(&jw, usingTagName)
	if jw.Error != nil {
		return false, 0, jw.Error
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(jw.Size()))

	started = true
	written, err = jw.DumpTo(w)
	return
}

// Unmarshal decodes the JSON in data into the object.
func Unmarshal(data []byte, v Unmarshaler) error {
	return UnmarshalCustom(data, v, "json")
}
func UnmarshalCustom(data []byte, v Unmarshaler, usingTagName string) error {
	l := jlexer.Lexer{Data: data}
	v.UnmarshalEasyJSON(&l, usingTagName)
	return l.Error()
}

// UnmarshalFromReader reads all the data in the reader and decodes as JSON into the object.
func UnmarshalFromReader(r io.Reader, v Unmarshaler) error {
	return UnmarshalFromReaderCustom(r, v, "json")
}
func UnmarshalFromReaderCustom(r io.Reader, v Unmarshaler, usingTagName string) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	// fmt.Println("found: ", string(data))
	l := jlexer.Lexer{Data: data}
	v.UnmarshalEasyJSON(&l, usingTagName)
	return l.Error()
}
