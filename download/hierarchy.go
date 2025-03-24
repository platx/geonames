package download

import (
	"context"
)

// Hierarchy parses hierarchy of toponyms from the hierarchy.zip file.
func (c *Client) Hierarchy(ctx context.Context) (Iterator[HierarchyItem], error) {
	res, err := c.downloadAndParseZIPFile(ctx, "hierarchy.zip")
	if err != nil {
		return nil, err
	}

	return withUnmarshalRows[HierarchyItem](res), nil
}
