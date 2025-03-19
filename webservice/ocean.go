package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathOcean = "/oceanJSON"

type OceanRequest struct {
	// Position the latitude and longitude of the search location
	Position value.Position `url:",dive"`
	// Radius buffer in km
	Radius int32 `url:"radius"`
}

// Ocean returns the ocean or sea for the given latitude/longitude.
func (c *Client) Ocean(ctx context.Context, req OceanRequest) (Ocean, error) {
	var res struct {
		Ocean Ocean `json:"ocean"`
	}

	err := c.apiRequest(
		ctx,
		pathOcean,
		req,
		&res,
	)

	return res.Ocean, err
}
