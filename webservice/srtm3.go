package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathSRTM3 = "/srtm3JSON"

// SRTM3 returns elevation in meters according to srtm3, ocean areas have been masked as "no data"
// and have been assigned a value of -32768.
// [More info]: https://www.geonames.org/export/web-services.html#srtm3
func (c *Client) SRTM3(ctx context.Context, position value.Position) (int32, error) {
	var res struct {
		Elevation int32 `json:"srtm3"`
	}

	err := c.apiRequest(
		ctx,
		pathSRTM3,
		position,
		&res,
	)

	return res.Elevation, err
}
