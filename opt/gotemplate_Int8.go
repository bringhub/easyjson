// generated by gotemplate

package opt

import (
	"fmt"

	"github.com/bringhub/easyjson/jlexer"
	"github.com/bringhub/easyjson/jwriter"
)

// template type Optional(A)

// A 'gotemplate'-based type for providing optional semantics without using pointers.
type Int8 struct {
	V       int8
	Defined bool
}

// Creates an optional type with a given value.
func OInt8(v int8) Int8 {
	return Int8{V: v, Defined: true}
}

// Get returns the value or given default in the case the value is undefined.
func (v Int8) Get(deflt int8) int8 {
	if !v.Defined {
		return deflt
	}
	return v.V
}

// MarshalEasyJSON does JSON marshaling using easyjson interface.
func (v Int8) MarshalEasyJSON(w *jwriter.Writer) {
	if v.Defined {
		w.Int8(v.V)
	} else {
		w.RawString("null")
	}
}

// UnmarshalEasyJSON does JSON unmarshaling using easyjson interface.
func (v *Int8) UnmarshalEasyJSON(l *jlexer.Lexer) {
	if l.IsNull() {
		l.Skip()
		*v = Int8{}
	} else {
		v.V = l.Int8()
		v.Defined = true
	}
}

// MarshalJSON implements a standard json marshaler interface.
func (v *Int8) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	v.MarshalEasyJSON(&w)
	return w.Buffer.BuildBytes(), w.Error
}

// UnmarshalJSON implements a standard json unmarshaler interface.
func (v *Int8) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{}
	v.UnmarshalEasyJSON(&l)
	return l.Error()
}

// IsDefined returns whether the value is defined, a function is required so that it can
// be used in an interface.
func (v Int8) IsDefined() bool {
	return v.Defined
}

// String implements a stringer interface using fmt.Sprint for the value.
func (v Int8) String() string {
	if !v.Defined {
		return "<undefined>"
	}
	return fmt.Sprint(v.V)
}
