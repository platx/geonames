package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathGTOPO30 = "/gtopo30JSON"

// GTOPO30 returns elevation in meters according to gtopo30, ocean areas have been masked as "no data"
// and have been assigned a value of -9999.
func (c *Client) GTOPO30(ctx context.Context, position value.Position) (int32, error) {
	var res struct {
		Elevation int32 `json:"gtopo30"`
	}

	err := c.apiRequest(
		ctx,
		pathGTOPO30,
		position,
		&res,
	)

	return res.Elevation, err
}
