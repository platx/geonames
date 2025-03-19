package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathCountrySubdivision = "/countrySubdivisionJSON"

type CountrySubdivisionRequest struct {
	// Position the latitude and longitude of the search location.
	Position value.Position `url:",dive"`
	// Radius buffer in km for closest country in coastal areas, a positive buffer expands the positive area
	// whereas a negative buffer reduces it.
	Radius int32 `url:"radius"`
	// Language default= names in local language.
	Language string `url:"lang"`
	// Level administrative level (1-5).
	Level uint8 `url:"level"`
}

// CountrySubdivision returns the country and the administrative subdivison (state, province,...)
// for the given latitude/longitude.
func (c *Client) CountrySubdivision(ctx context.Context, req CountrySubdivisionRequest) (CountrySubdivision, error) {
	var res CountrySubdivision

	err := c.apiRequest(
		ctx,
		pathCountrySubdivision,
		req,
		&res,
	)

	return res, err
}
