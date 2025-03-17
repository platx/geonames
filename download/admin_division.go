package download

import "context"

// AdminDivision1 parses names in English for admin divisions from the admin1CodesASCII.txt file.
func (c *Client) AdminDivision1(ctx context.Context, callback func(parsed AdminDivision) error) error {
	return c.adminDivision(ctx, "admin2Codes.txt", callback)
}

// AdminDivision2 parses names for administrative subdivision 'admin2 code' (UTF8) from the admin2Codes.txt file.
func (c *Client) AdminDivision2(ctx context.Context, callback func(parsed AdminDivision) error) error {
	return c.adminDivision(ctx, "admin2Codes.txt", callback)
}

func (c *Client) adminDivision(ctx context.Context, fileName string, callback func(parsed AdminDivision) error) error {
	return c.downloadAndParseFile(ctx, fileName, func(row []string) error {
		var parsed AdminDivision

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
