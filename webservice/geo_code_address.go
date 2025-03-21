package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathGeoCodeAddress = "/geoCodeAddressJSON"

type GeoCodeAddressRequest struct {
	// Query query term
	Query string `url:"q"`
	// Country default is all countries
	Country value.CountryCode `url:"country"`
	// PostalCode postal code
	PostalCode string `url:"postalcode"`
}

// GeoCodeAddress returns the nearest address for the given latitude/longitude.
// Supported countries: AT,AU,AX,CC,CH,CL,CX,CZ,DK,EE,ES,FI,FR,GF,GP,HK,IS,LU,MQ,NF,NL,NO,PL,PR,PT,RE,SG,SI,SJ,SK,US,YT.
// [More info]: https://www.geonames.org/maps/addresses.html#geoCodeAddress
func (c *Client) GeoCodeAddress(ctx context.Context, req GeoCodeAddressRequest) (Address, error) {
	var res struct {
		Address Address `json:"address"`
	}

	err := c.apiRequest(
		ctx,
		pathGeoCodeAddress,
		req,
		&res,
	)

	return res.Address, err
}
