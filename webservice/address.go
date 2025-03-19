package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathAddress = "/addressJSON"

type AddressRequest struct {
	// Position the latitude and longitude of the search location
	Position value.Position `url:",dive"`
	// Radius buffer in km (default=0.2)
	Radius int32 `url:"radius"`
	// Radius buffer in km (default=1)
	MaxRows int32 `url:"maxRows"`
}

// Address returns the nearest address for the given latitude/longitude.
// Supported countries: AT,AU,AX,CC,CH,CL,CX,CZ,DK,EE,ES,FI,FR,GF,GP,HK,IS,LU,MQ,NF,NL,NO,PL,PR,PT,RE,SG,SI,SJ,SK,US,YT.
func (c *Client) Address(ctx context.Context, req AddressRequest) ([]AddressNearby, error) {
	var res struct {
		Address []AddressNearby `json:"address"`
	}

	err := c.apiRequest(
		ctx,
		pathAddress,
		req,
		&res,
	)

	return res.Address, err
}
