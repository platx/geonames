package download

import (
	"context"
	"fmt"
	"time"
)

// AlternateNamesModifications parses all alternate names modified on the previous day from
// the alternateNamesModifications-{date}.txt file.
func (c *Client) AlternateNamesModifications(ctx context.Context, callback func(parsed AlternateName) error) error {
	fileName := fmt.Sprintf("alternateNamesModifications-%s.txt", yesterday().Format(time.DateOnly))

	return c.downloadAndParseFile(ctx, fileName, func(row []string) error {
		var parsed AlternateName

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
