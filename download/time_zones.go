package download

import "context"

// TimeZones parses timezone information from the timeZones.txt file.
func (c *Client) TimeZones(ctx context.Context) (Iterator[TimeZone], error) {
	res, err := c.downloadAndParseFile(ctx, "timeZones.txt")

	return withUnmarshalRows[TimeZone](withSkipHeader(res)), err
}
