package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathNeighbours = "/neighboursJSON"

type NeighboursRequest struct {
	// ID the geonameId for the neighbours (country or ADM)
	ID uint64 `url:"geonameId"`
	// Country alternative parameter instead of ID
	Country value.CountryCode `url:"country"`
}

// Neighbours returns the neighbours of a toponym, currently only implemented for countries.
// [More info]: https://www.geonames.org/export/place-hierarchy.html#neighbours
func (c *Client) Neighbours(ctx context.Context, req NeighboursRequest) ([]GeoName, error) {
	var res struct {
		Items []GeoName `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathNeighbours,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
