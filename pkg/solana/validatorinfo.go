package solana

import (
	"fmt"
	"net/url"
	"unicode/utf8"
)

const (
	MaxShortFieldLength = 80
	MaxLongFieldLength  = 300
	MaxValidatorInfo    = 576
)

type ValidatorInfo struct {
	Name    string  `pulumi:"name"`
	Website *string `pulumi:"website,optional"`
	IconURL *string `pulumi:"iconURL,optional"`
	Details *string `pulumi:"details,optional"`
}

func (i *ValidatorInfo) Check() error {
	if utf8.RuneCountInString(i.Name) > MaxShortFieldLength {
		return fmt.Errorf("name exceeds maximum length of %d", MaxShortFieldLength)
	}

	if i.Website != nil {
		if utf8.RuneCountInString(*i.Website) > MaxShortFieldLength {
			return fmt.Errorf("website exceeds maximum length of %d", MaxShortFieldLength)
		}
		if _, err := url.ParseRequestURI(*i.Website); err != nil {
			return fmt.Errorf("invalid website URL: %w", err)
		}
	}

	if i.IconURL != nil {
		if utf8.RuneCountInString(*i.IconURL) > MaxShortFieldLength {
			return fmt.Errorf("icon URL exceeds maximum length of %d", MaxShortFieldLength)
		}
		if _, err := url.ParseRequestURI(*i.IconURL); err != nil {
			return fmt.Errorf("invalid icon URL: %w", err)
		}
	}

	if i.Details != nil {
		if utf8.RuneCountInString(*i.Details) > MaxLongFieldLength {
			return fmt.Errorf("description exceeds maximum length of %d", MaxLongFieldLength)
		}
	}

	totalLength := utf8.RuneCountInString(i.Name)
	if i.Website != nil {
		totalLength += utf8.RuneCountInString(*i.Website)
	}
	if i.IconURL != nil {
		totalLength += utf8.RuneCountInString(*i.IconURL)
	}
	if i.Details != nil {
		totalLength += utf8.RuneCountInString(*i.Details)
	}

	if totalLength > MaxValidatorInfo {
		return fmt.Errorf("total length of fields exceeds maximum allowed length of %d", MaxValidatorInfo)
	}

	return nil
}
