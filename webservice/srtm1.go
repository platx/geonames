package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathSRTM1 = "/srtm1JSON"

// SRTM1 returns elevation in meters according to srtm1, ocean areas have been masked as "no data"
// and have been assigned a value of -32768.
// [More info]: https://www.geonames.org/export/web-services.html#srtm1
func (c *Client) SRTM1(ctx context.Context, position value.Position) (int32, error) {
	var res struct {
		Elevation int32 `json:"srtm1"`
	}

	err := c.apiRequest(
		ctx,
		pathSRTM1,
		position,
		&res,
	)

	return res.Elevation, err
}
