package download

import (
	"context"
)

// AdminDivisionFirst parses names in English for admin divisions from the admin1CodesASCII.txt file.
func (c *Client) AdminDivisionFirst(ctx context.Context, callback func(parsed AdminDivision) error) error {
	return c.adminDivision(ctx, "admin1CodesASCII.txt", callback)
}

// AdminDivisionSecond parses names for administrative subdivision 'admin2 code' (UTF8) from the admin2Codes.txt file.
func (c *Client) AdminDivisionSecond(ctx context.Context, callback func(parsed AdminDivision) error) error {
	return c.adminDivision(ctx, "admin2Codes.txt", callback)
}

// AdminDivisionFifth parses the new adm5 column which is not present in other files due of backward compatibility.
func (c *Client) AdminDivisionFifth(ctx context.Context, callback func(parsed AdminCode5) error) error {
	return c.downloadAndParseZIPFile(ctx, "adminCode5.zip", func(row []string) error {
		var parsed AdminCode5

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
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
