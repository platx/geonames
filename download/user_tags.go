package download

import (
	"context"
)

// UserTags parses toponyms not belonging to a country.
func (c *Client) UserTags(ctx context.Context) (Iterator[UserTag], error) {
	res, err := c.downloadAndParseZIPFile(ctx, "userTags.zip")

	return withUnmarshalRows[UserTag](res), err
}
