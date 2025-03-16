package download

import (
	"context"
	"fmt"
)

// FeatureCodes parses name and description for feature classes and feature codes from a featureCodes_xx.txt file.
func (c *Client) FeatureCodes(ctx context.Context, language string, callback func(parsed Feature) error) error {
	return c.downloadAndParseFile(ctx, fmt.Sprintf("featureCodes_%s.txt", language), func(row []string) error {
		var parsed Feature

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
