package download

import (
	"context"
	"fmt"

	"github.com/platx/geonames/value"
)

// Cities parses cities by size from the cities*.zip file.
func (c *Client) Cities(ctx context.Context, size value.Cities, callback func(parsed GeoName) error) error {
	return c.geoNames(ctx, fmt.Sprintf("%s.zip", size), callback)
}

// Cities500 parses cities with a population > 500 or seats of adm div down to PPLA4 (ca 185.000)
// from the cities500.zip file.
func (c *Client) Cities500(ctx context.Context, callback func(parsed GeoName) error) error {
	return c.Cities(ctx, value.Cities500, callback)
}

// Cities1000 parses with a population > 1000 or seats of adm div down to PPLA3 (ca 130.000)
// from the cities1000.zip file.
func (c *Client) Cities1000(ctx context.Context, callback func(parsed GeoName) error) error {
	return c.Cities(ctx, value.Cities1000, callback)
}

// Cities5000 parses with a population > 5000 or PPLA (ca 50.000)
// from the cities5000.zip file.
func (c *Client) Cities5000(ctx context.Context, callback func(parsed GeoName) error) error {
	return c.Cities(ctx, value.Cities5000, callback)
}

// Cities15000 parses with a population > 15000 or capitals (ca 25.000)
// from the cities15000.zip file.
func (c *Client) Cities15000(ctx context.Context, callback func(parsed GeoName) error) error {
	return c.Cities(ctx, value.Cities15000, callback)
}
