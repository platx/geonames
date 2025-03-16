package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathFindNearbyWikipedia = "/findNearbyWikipediaJSON"

type FindNearbyWikipediaRequest struct {
	// Position the latitude and longitude of the search location
	Position value.Position `url:",dive"`
	// Language ISO language code of article text
	Language string `url:"lang"`
	// Radius in km the maximal distance in km from the point specified via lat and lng that a result should be found
	Radius uint32 `url:"radius"`
	// MaxRows the maximal number of rows in the document returned by the service. Default is 5.
	MaxRows uint32 `url:"maxRows"`
	// Country default is all countries
	Country []value.CountryCode `url:"country"`
}

// FindNearbyWikipedia returns the closest toponym for the lat/lng query.
func (c *Client) FindNearbyWikipedia(
	ctx context.Context,
	req FindNearbyWikipediaRequest,
) ([]WikipediaNearby, error) {
	var res struct {
		Items []WikipediaNearby `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathFindNearbyWikipedia,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
