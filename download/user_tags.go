package download

import (
	"context"
)

// UserTags parses toponyms not belonging to a country.
func (c *Client) UserTags(ctx context.Context, callback func(parsed UserTag) error) error {
	return c.downloadAndParseZIPFile(ctx, "userTags.zip", func(row []string) error {
		var parsed UserTag

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
