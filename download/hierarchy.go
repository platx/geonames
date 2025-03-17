package download

import (
	"context"
)

// Hierarchy parses hierarchy of toponyms from the hierarchy.zip file.
func (c *Client) Hierarchy(ctx context.Context, callback func(parsed HierarchyItem) error) error {
	return c.downloadAndParseZIPFile(ctx, "hierarchy.zip", func(row []string) error {
		var parsed HierarchyItem

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
