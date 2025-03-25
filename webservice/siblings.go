package webservice

import (
	"context"
)

const pathSiblings = "/siblingsJSON"

type SiblingsRequest struct {
	// ID the geonameId for the siblings
	ID uint64 `url:"geonameId"`
}

// Siblings returns GeoName records (feature class A) that have the same administrative level and the same father.
// The top hierarchy (continent) is the first element in the list.
// [More info]: https://www.geonames.org/export/place-hierarchy.html#siblings
func (c *Client) Siblings(ctx context.Context, req SiblingsRequest) ([]GeoName, error) {
	var res struct {
		Items []GeoName `json:"geonames"`
	}

	err := c.apiRequest(
		ctx,
		pathSiblings,
		req,
		&res,
	)

	return res.Items, err
}
