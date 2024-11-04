package utils

import (
	"strconv"
)

type FlagBuilder []string

func (f *FlagBuilder) Append(s ...string) {
	*f = append(*f, s...)
}

func (f *FlagBuilder) ToArgs() []string {
	var res []string

	return append(res, *f...)
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
