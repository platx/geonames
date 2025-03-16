package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathCountryCode = "/countryCodeJSON"

type CountryCodeRequest struct {
	// Position the latitude and longitude of the search location
	Position value.Position `url:",dive"`
	// Radius buffer in km for closest country in coastal areas, a positive buffer expands the positive area
	// whereas a negative buffer reduces it
	Radius int32 `url:"radius"`
	// Language of returned 'name' element (the pseudo language code 'local' will return it in local language)
	Language string `url:"lang"`
}

// CountryCode returns the iso country code for the given latitude/longitude.
func (c *Client) CountryCode(ctx context.Context, req CountryCodeRequest) (CountryNearby, error) {
	var res CountryNearby

	err := c.apiRequest(
		ctx,
		pathCountryCode,
		req,
		&res,
	)

	return res, err
}
