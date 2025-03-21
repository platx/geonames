package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathAstergdem = "/astergdemJSON"

// Astergdem returns elevation in meters according to aster gdem, ocean areas have been masked as "no data"
// and have been assigned a value of -32768.
// [More info]: https://www.geonames.org/export/web-services.html#astergdem
func (c *Client) Astergdem(ctx context.Context, position value.Position) (int32, error) {
	var res struct {
		Elevation int32 `json:"astergdem"`
	}

	err := c.apiRequest(
		ctx,
		pathAstergdem,
		position,
		&res,
	)

	return res.Elevation, err
}
