package download

import (
	"context"
)

// AllCountries parses all toponyms from the allCountries.zip file.
func (c *Client) AllCountries(ctx context.Context, callback func(parsed GeoName) error) error {
	return c.geoNames(ctx, "allCountries.zip", callback)
}
