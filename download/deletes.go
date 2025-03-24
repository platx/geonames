package download

import (
	"context"
	"fmt"
	"time"
)

// Deletes parses all records deleted on the previous day from the deletes-{date}.txt file.
func (c *Client) Deletes(ctx context.Context) (Iterator[GeoNameDeleted], error) {
	fileName := fmt.Sprintf("deletes-%s.txt", yesterday().Format(time.DateOnly))

	res, err := c.downloadAndParseFile(ctx, fileName)

	return withUnmarshalRows[GeoNameDeleted](res), err
}
