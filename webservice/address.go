package webservice

import (
	"context"
	"encoding/json"

	"github.com/platx/geonames/value"
)

const pathAddress = "/addressJSON"

type AddressRequest struct {
	// Position the latitude and longitude of the search location
	Position value.Position `url:",dive"`
	// Radius buffer in km (default=0.2)
	Radius int32 `url:"radius"`
	// MaxRows the maximal number of rows returned by the service. Default is 1.
	MaxRows uint32 `url:"maxRows"`
}

// Address returns the nearest address for the given latitude/longitude.
// Supported countries: AT,AU,AX,CC,CH,CL,CX,CZ,DK,EE,ES,FI,FR,GF,GP,HK,IS,LU,MQ,NF,NL,NO,PL,PR,PT,RE,SG,SI,SJ,SK,US,YT.
// [More info]: https://www.geonames.org/maps/addresses.html#address
func (c *Client) Address(ctx context.Context, req AddressRequest) ([]AddressNearby, error) {
	var res struct {
		Result addressResult `json:"address"`
	}

	err := c.apiRequest(
		ctx,
		pathAddress,
		req,
		&res,
	)

	return res.Result.Address, err
}

type addressResult struct {
	Address []AddressNearby `json:"address"`
}

func (a *addressResult) UnmarshalJSON(data []byte) error {
	var single AddressNearby
	if err := json.Unmarshal(data, &single); err == nil {
		a.Address = []AddressNearby{single}
		return nil
	}

	var multiple []AddressNearby
	if err := json.Unmarshal(data, &multiple); err != nil {
		return err
	}

	a.Address = multiple
	return nil
}
