package apt

import (
	"errors"
)

type SignedBy struct {
	Paths     *[]string `pulumi:"paths,optional"`
	PublicKey *string   `pulumi:"publicKey,optional"`
}

type Source struct {
	Types      []string `pulumi:"types"`
	URIs       []string `pulumi:"URIs"`
	Suites     []string `pulumi:"suites"`
	Components []string `pulumi:"components"`

	AllowDowngradeToInsecure *bool     `pulumi:"allowDowngradeToInsecure,optional"`
	AllowInsecure            *bool     `pulumi:"allowInsecure,optional"`
	AllowWeak                *bool     `pulumi:"allowWeak,optional"`
	Architectures            *[]string `pulumi:"architectures,optional"`
	CheckDate                *bool     `pulumi:"checkDate,optional"`
	CheckValidUntil          *bool     `pulumi:"checkValidUntil,optional"`
	DateMaxFuture            *int      `pulumi:"dateMaxFuture,optional"`
	InReleasePath            *string   `pulumi:"inReleasePath,optional"`
	SignedBy                 *SignedBy `pulumi:"signedBy,optional"`
	Snapshot                 *string   `pulumi:"snapshot,optional"`
	Trusted                  *bool     `pulumi:"trusted,optional"`
	ValidUntilMax            *int      `pulumi:"validUntilMax,optional"`
	ValidUntilMin            *int      `pulumi:"validUntilMin,optional"`

	ExtraLines *[]string `pulumi:"extraLines,optional"`
}

func (s Source) AppendText(bytes []byte) ([]byte, error) {
	b := NewDeb822Builder(bytes)

	b.AppendArrayString("Types", s.Types)
	b.AppendArrayString("URIs", s.URIs)
	b.AppendArrayString("Suites", s.Suites)
	b.AppendArrayString("Components", s.Components)

	b.AppendBoolP("Allow-Downgrade-To-Insecure", s.AllowDowngradeToInsecure)
	b.AppendBoolP("Allow-Insecure", s.AllowInsecure)
	b.AppendBoolP("Allow-Weak", s.AllowWeak)
	b.AppendArrayStringP("Architectures", s.Architectures)
	b.AppendBoolP("Check-Date", s.CheckDate)
	b.AppendBoolP("Check-Valid-Until", s.CheckValidUntil)
	b.AppendIntP("Date-Max-Future", s.DateMaxFuture)
	b.AppendStringP("In-Release-Path", s.InReleasePath)

	if t := s.SignedBy; t != nil {
		found := 0

		if t.Paths != nil {
			b.AppendArrayStringP("Signed-By", t.Paths)
			found += 1
		}

		if t.PublicKey != nil {
			b.AppendMultilineString("Signed-By", *t.PublicKey)
			found += 1
		}

		if found > 1 {
			return nil, errors.New("only one type of signed-by is supported at a time")
		}
	}

	b.AppendStringP("Snapshot", s.Snapshot)
	b.AppendBoolP("Trusted", s.Trusted)
	b.AppendIntP("Valid-Until-Max", s.ValidUntilMax)
	b.AppendIntP("Valid-Until-Min", s.ValidUntilMin)

	if s.ExtraLines != nil {
		for _, v := range *s.ExtraLines {
			b.Append(v)
		}
	}

	return b.Bytes(), nil
}

func (s Source) MarshalText() ([]byte, error) {
	return s.AppendText(nil)
}

type Sources []Source

func (s Sources) AppendText(b []byte) ([]byte, error) {
	for _, v := range s {
		var err error

		b, err = v.AppendText(b)

		if err != nil {
			return nil, err
		}

		b = append(b, '\n')
	}

	return b, nil
}

func (s Sources) MarshalText() ([]byte, error) {
	return s.AppendText(nil)
}
