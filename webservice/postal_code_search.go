package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathPostalCodeSearch = "/postalCodeSearchJSON"

type PostalCodeSearchRequest struct {
	// PostalCode postal code
	PostalCode string `url:"postalcode"`
	// PostalCodeStartsWith the first characters or letters of a postal code
	PostalCodeStartsWith string `url:"postalcode_startsWith"`
	// PlaceName all fields: placename, postal code, country, admin name
	PlaceName string `url:"placename"`
	// PlaceNameStartsWith the first characters of a place name
	PlaceNameStartsWith string `url:"placename_startsWith"`
	// Country default is all countries
	Country []value.CountryCode `url:"country"`
	// CountryBias records from the countryBias are listed first
	CountryBias value.CountryCode `url:"countryBias"`
	// MaxRows the maximal number of rows in the document returned by the service. Default is 10.
	MaxRows uint32 `url:"maxRows"`
	// Operator default is 'AND', with the operator 'OR' not all search terms need to be matched by the response
	Operator value.Operator `url:"operator"`
	// Reduced default is false, when set to true only the UK outer codes respectivel the NL 4-digits are returned.
	// Attention: the default value on the commercial servers is currently set to true.
	// It will be changed later to false.
	Reduced bool `url:"isReduced"`
	// BoundingBox only features within the box are returned.
	BoundingBox value.BoundingBox `url:",dive"`
}

// PostalCodeSearch returns a list of postal codes and places for the placename/postalcode query as xml document
// For the US the first returned zip code is determined using zip code area shapes, the following zip codes
// are based on the centroid. For all other supported countries all returned postal codes are based on centroids.
// [More info]: https://www.geonames.org/export/web-services.html#postalCodeSearch
func (c *Client) PostalCodeSearch(ctx context.Context, req PostalCodeSearchRequest) ([]PostalCode, error) {
	var res struct {
		Items []PostalCode `json:"postalCodes"`
	}

	if err := c.apiRequest(
		ctx,
		pathPostalCodeSearch,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
