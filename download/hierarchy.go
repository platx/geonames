package download

import (
	"context"
)

// Hierarchy parses hierarchy of toponyms from the hierarchy.zip file.
func (c *Client) Hierarchy(ctx context.Context) (Iterator[HierarchyItem], error) {
	res, err := c.downloadAndParseZIPFile(ctx, "hierarchy.zip")

	return withUnmarshalRows[HierarchyItem](res), err
}
