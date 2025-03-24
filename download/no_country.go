package download

import (
	"context"
)

// NoCountry parses toponyms not belonging to a country.
func (c *Client) NoCountry(ctx context.Context) (Iterator[GeoName], error) {
	return c.geoNames(ctx, "no-country.zip")
}
