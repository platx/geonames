package download

import (
	"context"
)

// Languages parses ISO-639 language codes, as used for alternate names in the file alternateNames.zip.
func (c *Client) Languages(ctx context.Context) (Iterator[Language], error) {
	res, err := c.downloadAndParseFile(ctx, "iso-languagecodes.txt")
	if err != nil {
		return nil, err
	}

	return withUnmarshalRows[Language](withSkipHeader(res)), nil
}
