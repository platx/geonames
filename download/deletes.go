package download

import (
	"context"
	"fmt"
	"time"
)

// Deletes parses all records deleted on the previous day from the deletes-{date}.txt file.
func (c *Client) Deletes(ctx context.Context, callback func(parsed GeoNameDeleted) error) error {
	fileName := fmt.Sprintf("deletes-%s.txt", yesterday().Format(time.DateOnly))

	return c.downloadAndParseFile(ctx, fileName, func(row []string) error {
		var parsed GeoNameDeleted

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
