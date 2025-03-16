package download

import (
	"context"
)

// NoCountry parses toponyms not belonging to a country.
func (c *Client) NoCountry(ctx context.Context, callback func(parsed GeoName) error) error {
	return c.geoNames(ctx, "no-country.zip", callback)
}
