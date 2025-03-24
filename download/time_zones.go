package download

import "context"

// TimeZones parses timezone information from the timeZones.txt file.
func (c *Client) TimeZones(ctx context.Context) (Iterator[TimeZone], error) {
	res, err := c.downloadAndParseFile(ctx, "timeZones.txt")
	if err != nil {
		return nil, err
	}

	return withUnmarshalRows[TimeZone](withSkipHeader(res)), nil
}
