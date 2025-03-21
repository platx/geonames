package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathPostalCodeLookup = "/postalCodeLookupJSON"

type PostalCodeLookupRequest struct {
	// PostalCode postal code
	PostalCode string `url:"postalcode"`
	// Country default is all countries
	Country []value.CountryCode `url:"country"`
	// MaxRows the maximal number of rows in the document returned by the service. Default is 20.
	MaxRows uint32 `url:"maxRows"`
}

// PostalCodeLookup Placename lookup with postalcode.
// [More info]: https://www.geonames.org/export/web-services.html#postalCodeLookupJSON
func (c *Client) PostalCodeLookup(ctx context.Context, req PostalCodeLookupRequest) ([]PostalCode, error) {
	var res struct {
		Items []PostalCode `json:"postalCodes"`
	}

	if err := c.apiRequest(
		ctx,
		pathPostalCodeLookup,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
