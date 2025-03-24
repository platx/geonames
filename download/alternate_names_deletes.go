package download

import (
	"context"
	"fmt"
	"time"
)

// AlternateNamesDeletes parses all alternate names deleted on the previous day from
// the alternateNamesDeletes-{date}.txt file.
func (c *Client) AlternateNamesDeletes(ctx context.Context) (Iterator[AlternateNameDeleted], error) {
	fileName := fmt.Sprintf("alternateNamesDeletes-%s.txt", yesterday().Format(time.DateOnly))

	res, err := c.downloadAndParseFile(ctx, fileName)
	if err != nil {
		return nil, err
	}

	return withUnmarshalRows[AlternateNameDeleted](res), nil
}
