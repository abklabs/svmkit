package apt

import (
	"bytes"
	"strconv"
	"strings"
)

func NewDeb822Builder(b []byte) *Deb822Builder {
	return &Deb822Builder{buf: bytes.NewBuffer(b)}
}

type Deb822Builder struct {
	buf *bytes.Buffer
}

func (s *Deb822Builder) Append(v string) {
	s.buf.WriteString(v)
	s.buf.WriteString("\n")
}

func (s *Deb822Builder) AppendString(k string, v string) {
	s.buf.WriteString(k)
	s.buf.WriteString(": ")
	s.Append(v)
}

func (s *Deb822Builder) AppendMultilineString(k string, v string) {
	s.buf.WriteString(k)
	s.buf.WriteString(":\n")

	lines := strings.Split(v, "\n")

	for _, v := range lines {
		if len(v) == 0 {
			s.Append(" .")
		} else {
			s.buf.WriteString(" ")
			s.Append(v)
		}
	}
}

func (s *Deb822Builder) AppendStringP(k string, v *string) {
	if v == nil {
		return
	}

	s.AppendString(k, *v)
}

func (s *Deb822Builder) AppendBoolP(k string, v *bool) {
	if v == nil {
		return
	}

	s.AppendString(k, strconv.FormatBool(*v))
}

func (s *Deb822Builder) AppendArrayString(k string, v []string) {
	s.AppendString(k, strings.Join(v, " "))
}

func (s *Deb822Builder) AppendArrayStringP(k string, v *[]string) {
	if v == nil {
		return
	}

	s.AppendArrayString(k, *v)
}

func (s *Deb822Builder) AppendIntP(k string, v *int) {
	if v == nil {
		return
	}

	s.AppendString(k, strconv.Itoa(*v))
}

func (s *Deb822Builder) Bytes() []byte {
	return s.buf.Bytes()
}
