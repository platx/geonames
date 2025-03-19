package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathSRTM3 = "/srtm3JSON"

// SRTM3 returns elevation in meters according to srtm3, ocean areas have been masked as "no data"
// and have been assigned a value of -32768.
func (c *Client) SRTM3(ctx context.Context, position value.Position) (int32, error) {
	var res struct {
		SRTM3 int32 `json:"srtm3"`
	}

	err := c.apiRequest(
		ctx,
		pathSRTM3,
		position,
		&res,
	)

	return res.SRTM3, err
}
