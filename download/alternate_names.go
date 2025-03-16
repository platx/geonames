package download

import (
	"context"
)

// AlternateNames parses alternate names for toponyms from the alternateNamesV2.zip file.
func (c *Client) AlternateNames(ctx context.Context, callback func(parsed AlternateName) error) error {
	return c.downloadAndParseZIPFile(ctx, "alternateNamesV2.zip", func(row []string) error {
		var parsed AlternateName

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
