package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathFindNearby = "/findNearbyJSON"

type FindNearbyRequest struct {
	// Position the latitude and longitude of the search location
	Position value.Position `url:",dive"`
	// FeatureClass default is all feature classes
	FeatureClass []string `url:"featureClass"`
	// FeatureCode default is all feature codes
	FeatureCode []string `url:"featureCode"`
	// Radius in km the maximal distance in km from the point specified via lat and lng that a result should be found
	Radius uint32 `url:"radius"`
	// MaxRows the maximal number of rows in the document returned by the service. Default is 10.
	MaxRows uint32 `url:"maxRows"`
	// LocalCountry in border areas this parameter will restrict the search on the local country, value=true
	LocalCountry bool `url:"localCountry"`
}

// FindNearby returns the closest toponym for the lat/lng query.
func (c *Client) FindNearby(ctx context.Context, req FindNearbyRequest) ([]GeoNameNearby, error) {
	var res struct {
		Items []GeoNameNearby `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathFindNearby,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
