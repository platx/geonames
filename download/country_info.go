package download

import (
	"context"
)

// CountryInfo parses country information from the countryInfo.txt file.
func (c *Client) CountryInfo(ctx context.Context, callback func(parsed Country) error) error {
	return c.downloadAndParseFile(ctx, "countryInfo.txt", func(row []string) error {
		var parsed Country

		if err := parsed.UnmarshalRow(row); err != nil {
			return err
		}

		return callback(parsed)
	})
}
