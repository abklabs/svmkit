package utils

import (
	"strconv"
	"strings"
)

type FlagBuilder []string

func (f *FlagBuilder) Append(s ...string) {
	*f = append(*f, s...)
}

func (f *FlagBuilder) ToArgs() []string {
	var res []string

	return append(res, *f...)
}

// XXX - We need to work backwards from this, and make sure we're
// quoting spaces.
func (f *FlagBuilder) String() string {
	return strings.Join(f.ToArgs(), " ")
}

func (f *FlagBuilder) S(k string, v *string) {
	if v == nil {
		return
	}

	f.Append("--"+k, *v)
}

func (f *FlagBuilder) I64(k string, v *int64) {
	if v == nil {
		return
	}

	f.Append("--"+k, strconv.FormatInt(*v, 10))
}

func (f *FlagBuilder) I(k string, v *int) {
	if v == nil {
		return
	}

	f.Append("--"+k, strconv.FormatInt(int64(*v), 10))
}

func (f *FlagBuilder) B(k string, v *bool) {
	if v == nil {
		return
	}

	if !*v {
		return
	}

	f.Append("--" + k)
}
