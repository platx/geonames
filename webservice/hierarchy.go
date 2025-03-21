package webservice

import (
	"context"
)

const pathHierarchy = "/hierarchyJSON"

type HierarchyRequest struct {
	// ID the geonameId for the hierarchy
	ID uint64 `url:"geonameId"`
}

// Hierarchy returns a list of GeoName records, ordered by hierarchy level.
// The top hierarchy (continent) is the first element in the list.
// [More info]: https://www.geonames.org/export/place-hierarchy.html#hierarchy
func (c *Client) Hierarchy(ctx context.Context, req HierarchyRequest) ([]GeoName, error) {
	var res struct {
		Items []GeoName `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathHierarchy,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
