package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathChildren = "/childrenJSON"

type ChildrenRequest struct {
	// ID the geonameId of the parent
	ID uint64 `url:"geonameId"`
	// MaxRows number of rows returned, default is 200.
	MaxRows uint32 `url:"maxRows"`
	// Hierarchy allows to use other hiearchies then the default administrative hierarchy.
	Hierarchy value.Hierarchy `url:"hierarchy"`
}

// Children Returns the children (admin divisions and populated places) for a given geonameId.
// The children are the administrative divisions within an other administrative division,
// like the counties (ADM2) in a state (ADM1) or also the countries in a continent.
// The leafs are populated places, other feature classes like spots, mountains etc are not included in this service.
// Use the Search service if you need other feature classes as well.
// The top hierarchy (continent) is the first element in the list.
// [More info]: https://www.geonames.org/export/place-hierarchy.html#children
func (c *Client) Children(ctx context.Context, req ChildrenRequest) ([]GeoName, error) {
	var res struct {
		Items []GeoName `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathChildren,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
