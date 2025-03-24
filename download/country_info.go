package download

import (
	"context"
)

// CountryInfo parses country information from the countryInfo.txt file.
func (c *Client) CountryInfo(ctx context.Context) (Iterator[Country], error) {
	res, err := c.downloadAndParseFile(ctx, "countryInfo.txt")

	return withUnmarshalRows[Country](res), err
}
