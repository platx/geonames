package webservice

import (
	"context"
)

const pathSiblings = "/siblingsJSON"

type SiblingsRequest struct {
	// ID the geonameId for the hierarchy
	ID uint64 `url:"geonameId"`
}

// Siblings returns GeoName records (feature class A) that have the same administrative level and the same father.
// The top hierarchy (continent) is the first element in the list.
func (c *Client) Siblings(ctx context.Context, req SiblingsRequest) ([]GeoName, error) {
	var res struct {
		Items []GeoName `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathSiblings,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
