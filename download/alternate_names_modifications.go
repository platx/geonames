package download

import (
	"context"
	"fmt"
	"time"
)

// AlternateNamesModifications parses all alternate names modified on the previous day from
// the alternateNamesModifications-{date}.txt file.
func (c *Client) AlternateNamesModifications(ctx context.Context) (Iterator[AlternateName], error) {
	fileName := fmt.Sprintf("alternateNamesModifications-%s.txt", yesterday().Format(time.DateOnly))

	res, err := c.downloadAndParseFile(ctx, fileName)
	if err != nil {
		return nil, err
	}

	return withUnmarshalRows[AlternateName](res), nil
}
