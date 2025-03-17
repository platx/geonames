package download

import "context"

// AdminSubdivision1 parses names in English for admin divisions from the admin1CodesASCII.txt file.
func (c *Client) AdminSubdivision1(ctx context.Context, callback func(parsed AdminSubdivision) error) error {
	return c.downloadAndParseFile(ctx, "admin1CodesASCII.txt", func(row []string) error {
		var parsed AdminSubdivision

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
