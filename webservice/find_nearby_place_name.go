package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathFindNearbyPlaceName = "/findNearbyPlaceNameJSON"

type FindNearbyPlaceNameRequest struct {
	// Position the latitude and longitude of the search location
	Position value.Position `url:",dive"`
	// Language of returned 'name' element (the pseudo language code 'local' will return it in local language)
	Language string `url:"lang"`
	// Radius in km the maximal distance in km from the point specified via lat and lng that a result should be found
	Radius int32 `url:"radius"`
	// MaxRows the maximal number of rows returned by the service. Default is 10.
	MaxRows uint32 `url:"maxRows"`
	// LocalCountry in border areas this parameter will restrict the search on the local country, value=true
	LocalCountry bool `url:"localCountry"`
	// Cities used to categorize the populated places into three groups according to size/relevance
	Cities value.Cities `url:"cities"`
}

// FindNearbyPlaceName returns the closest populated place (feature class=P) for the lat/lng query.
// The unit of the distance element is 'km'.
// [More info]: https://www.geonames.org/export/web-services.html#findNearbyPlaceName
func (c *Client) FindNearbyPlaceName(
	ctx context.Context,
	req FindNearbyPlaceNameRequest,
) ([]GeoNameNearby, error) {
	var res struct {
		Items []GeoNameNearby `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathFindNearbyPlaceName,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
