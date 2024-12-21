package runner

import (
	"strconv"
	"strings"
)

type FlagBuilder []string

func (f *FlagBuilder) AppendRaw(s ...string) {
	*f = append(*f, s...)
}

func (f *FlagBuilder) Args() []string {
	var res []string

	return append(res, *f...)
}

// XXX - We need to work backwards from this, and make sure we're
// quoting spaces.
func (f *FlagBuilder) String() string {
	return strings.Join(f.Args(), " ")
}

func (f *FlagBuilder) AppendP(k string, v *string) {
	if v == nil {
		return
	}

	f.Append(k, *v)
}

func (f *FlagBuilder) Append(k string, v string) {
	f.AppendRaw("--"+k, v)
}

func (f *FlagBuilder) AppendArrayP(k string, array *[]string) {
	if array == nil {
		return
	}

	f.AppendArray(k, *array)
}

func (f *FlagBuilder) AppendArray(k string, array []string) {
	for _, v := range array {
		f.Append(k, v)
	}
}

func (f *FlagBuilder) AppendInt64P(k string, v *int64) {
	if v == nil {
		return
	}

	f.AppendRaw("--"+k, strconv.FormatInt(*v, 10))
}

func (f *FlagBuilder) AppendIntP(k string, v *int) {
	if v == nil {
		return
	}

	f.AppendRaw("--"+k, strconv.FormatInt(int64(*v), 10))
}

func (f *FlagBuilder) AppendFloat64(k string, v float64) {
	f.AppendRaw(k, strconv.FormatFloat(v, 'f', -1, 64))
}

func (f *FlagBuilder) AppendFloat64P(k string, v *float64) {
	if v == nil {
		return
	}

	f.AppendFloat64("--"+k, *v)
}

func (f *FlagBuilder) AppendBoolP(k string, v *bool) {
	if v == nil {
		return
	}

	if !*v {
		return
	}

	f.AppendRaw("--" + k)
}
