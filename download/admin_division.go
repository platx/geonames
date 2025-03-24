package download

import (
	"context"
)

// AdminDivisionFirst parses names in English for admin divisions from the admin1CodesASCII.txt file.
func (c *Client) AdminDivisionFirst(ctx context.Context) (Iterator[AdminDivision], error) {
	return c.adminDivision(ctx, "admin1CodesASCII.txt")
}

// AdminDivisionSecond parses names for administrative subdivision 'admin2 code' (UTF8) from the admin2Codes.txt file.
func (c *Client) AdminDivisionSecond(ctx context.Context) (Iterator[AdminDivision], error) {
	return c.adminDivision(ctx, "admin2Codes.txt")
}

// AdminDivisionFifth parses the new adm5 column which is not present in other files due of backward compatibility.
func (c *Client) AdminDivisionFifth(ctx context.Context) (Iterator[AdminCode5], error) {
	res, err := c.downloadAndParseZIPFile(ctx, "adminCode5.zip")
	if err != nil {
		return nil, err
	}

	return withUnmarshalRows[AdminCode5](res), nil
}

func (c *Client) adminDivision(ctx context.Context, fileName string) (Iterator[AdminDivision], error) {
	res, err := c.downloadAndParseFile(ctx, fileName)
	if err != nil {
		return nil, err
	}

	return withUnmarshalRows[AdminDivision](res), nil
}
