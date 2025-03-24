package download

import (
	"context"
	"fmt"

	"github.com/platx/geonames/value"
)

// ByCountry parses toponyms for country with iso code XX.
func (c *Client) ByCountry(ctx context.Context, code value.CountryCode) (Iterator[GeoName], error) {
	return c.geoNames(ctx, fmt.Sprintf("%s.zip", code))
}
