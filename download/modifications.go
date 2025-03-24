package download

import (
	"context"
	"fmt"
	"time"
)

// Modifications parses all records modified on the previous day from the modifications-{date}.txt file.
func (c *Client) Modifications(ctx context.Context) (Iterator[GeoName], error) {
	fileName := fmt.Sprintf("modifications-%s.txt", yesterday().Format(time.DateOnly))

	res, err := c.downloadAndParseFile(ctx, fileName)

	return withUnmarshalRows[GeoName](res), err
}
