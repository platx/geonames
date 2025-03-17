package download

import (
	"context"
	"fmt"
	"time"
)

// AlternateNamesDeletes parses all alternate names deleted on the previous day from
// the alternateNamesDeletes-{date}.txt file.
func (c *Client) AlternateNamesDeletes(ctx context.Context, callback func(parsed AlternateNameDeleted) error) error {
	fileName := fmt.Sprintf("alternateNamesDeletes-%s.txt", yesterday().Format(time.DateOnly))

	return c.downloadAndParseFile(ctx, fileName, func(row []string) error {
		var parsed AlternateNameDeleted

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
