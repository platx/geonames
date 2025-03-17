package download

import (
	"context"
	"fmt"
	"time"
)

// Modifications parses all records modified on the previous day from the modifications-{date}.txt file.
func (c *Client) Modifications(ctx context.Context, callback func(parsed GeoName) error) error {
	fileName := fmt.Sprintf("modifications-%s.txt", yesterday().Format(time.DateOnly))

	return c.downloadAndParseFile(ctx, fileName, func(row []string) error {
		var parsed GeoName

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
