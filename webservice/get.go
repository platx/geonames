package webservice

import (
	"context"
)

const pathGet = "/getJSON"

type GeoNameGetRequest struct {
	// ID geonameId
	ID uint64 `url:"geonameId"`
	// Language place name and country name will be returned in the specified language.
	// Default is English. With the pseudo language code 'local' the local language will be returned.
	// Feature classes and codes are only available in English and Bulgarian.
	Language string `url:"lang"`
}

// Get returns the attribute of the geoNames feature with the given geonameId.
func (c *Client) Get(ctx context.Context, req GeoNameGetRequest) (GeoNameDetailed, error) {
	var res GeoNameDetailed

	err := c.apiRequest(
		ctx,
		pathGet,
		req,
		&res,
	)

	return res, err
}
