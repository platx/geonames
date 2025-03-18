package webservice

import (
	"context"
)

const pathContains = "/containsJSON"

type ContainsRequest struct {
	// ID the geonameId for enclosing feature
	ID uint64 `url:"geonameId"`
	// FeatureClass filter by featureClass
	FeatureClass string `url:"featureClass"`
	// FeatureCode filter by featureCode
	FeatureCode string `url:"featureCode"`
}

// Contains returns all features within the GeoName feature for the given geoNameId.
// It only returns contained features when a polygon boundary for the input feature is defined.
func (c *Client) Contains(ctx context.Context, req ContainsRequest) ([]GeoName, error) {
	var res struct {
		Items []GeoName `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathContains,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
