package solana

import (
	"fmt"
	"net/url"
)

const (
	MAX_SHORT_FIELD_LENGTH = 80
	MAX_LONG_FIELD_LENGTH  = 300
	MAX_VALIDATOR_INFO     = 576
)

type ValidatorInfo struct {
	Name    string  `pulumi:"name"`
	Website *string `pulumi:"website,optional"`
	IconURL *string `pulumi:"iconURL,optional"`
	Details *string `pulumi:"details,optional"`
}

func (i *ValidatorInfo) Check() error {
	if len(i.Name) > MAX_SHORT_FIELD_LENGTH {
		return fmt.Errorf("name exceeds maximum length of %d", MAX_SHORT_FIELD_LENGTH)
	}

	if i.Website != nil {
		if len(*i.Website) > MAX_SHORT_FIELD_LENGTH {
			return fmt.Errorf("website exceeds maximum length of %d", MAX_SHORT_FIELD_LENGTH)
		}
		if _, err := url.ParseRequestURI(*i.Website); err != nil {
			return fmt.Errorf("invalid website URL: %w", err)
		}
	}

	if i.IconURL != nil {
		if len(*i.IconURL) > MAX_SHORT_FIELD_LENGTH {
			return fmt.Errorf("icon URL exceeds maximum length of %d", MAX_SHORT_FIELD_LENGTH)
		}
		if _, err := url.ParseRequestURI(*i.IconURL); err != nil {
			return fmt.Errorf("invalid icon URL: %w", err)
		}
	}

	if i.Details != nil {
		if len(*i.Details) > MAX_LONG_FIELD_LENGTH {
			return fmt.Errorf("description exceeds maximum length of %d", MAX_LONG_FIELD_LENGTH)
		}
	}

	totalLength := len(i.Name)
	if i.Website != nil {
		totalLength += len(*i.Website)
	}
	if i.IconURL != nil {
		totalLength += len(*i.IconURL)
	}
	if i.Details != nil {
		totalLength += len(*i.Details)
	}

	if totalLength > MAX_VALIDATOR_INFO {
		return fmt.Errorf("total length of fields exceeds maximum allowed length of %d", MAX_VALIDATOR_INFO)
	}

	return nil
}
