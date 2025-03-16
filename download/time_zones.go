package download

import "context"

// TimeZones parses timezone information from the timeZones.txt file.
func (c *Client) TimeZones(ctx context.Context, callback func(parsed TimeZone) error) error {
	header := true

	return c.downloadAndParseFile(ctx, "timeZones.txt", func(row []string) error {
		if header {
			header = false

			return nil
		}

		var parsed TimeZone

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
