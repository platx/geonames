package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathFindNearbyPostalCodes = "/findNearbyPostalCodesJSON"

type FindNearbyPostalCodesRequest struct {
	// Position the latitude and longitude of the search location
	Position value.Position `url:",dive"`
	// Radius in km the maximal distance in km from the point specified via lat and lng that a result should be found
	Radius int32 `url:"radius"`
	// MaxRows the maximal number of rows returned by the service. Default is 5.
	MaxRows uint32 `url:"maxRows"`
	// Country default is all countries
	Country value.CountryCode `url:"country"`
	// LocalCountry in border areas this parameter will restrict the search on the local country
	LocalCountry bool `url:"localCountry"`
}

// FindNearbyPostalCodes returns a list of postalcodes and places for the lat/lng query.
// The result is sorted by distance. For Canada the FSA is returned (first 3 characters of full postal code).
// [More info]: https://www.geonames.org/export/web-services.html#findNearbyPostalCodes
func (c *Client) FindNearbyPostalCodes(
	ctx context.Context,
	req FindNearbyPostalCodesRequest,
) ([]PostalCodeNearby, error) {
	var res struct {
		Items []PostalCodeNearby `json:"postalCodes"`
	}

	err := c.apiRequest(
		ctx,
		pathFindNearbyPostalCodes,
		req,
		&res,
	)

	return res.Items, err
}
