package download

import (
	"context"
)

// AlternateNames parses alternate names for toponyms from the alternateNamesV2.zip file.
func (c *Client) AlternateNames(ctx context.Context) (Iterator[AlternateName], error) {
	res, err := c.downloadAndParseZIPFile(ctx, "alternateNamesV2.zip")

	return withUnmarshalRows[AlternateName](res), err
}
