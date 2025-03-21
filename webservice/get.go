package webservice

import (
	"context"
)

const pathGet = "/getJSON"

type GetRequest struct {
	// ID geonameId
	ID uint64 `url:"geonameId"`
	// Language place name and country name will be returned in the specified language.
	// Default is English. With the pseudo language code 'local' the local language will be returned.
	// Feature classes and codes are only available in English and Bulgarian.
	Language string `url:"lang"`
}

// Get returns the attribute of the geoNames feature with the given geonameId.
// [More info]: https://www.geonames.org/export/web-services.html#get
func (c *Client) Get(ctx context.Context, req GetRequest) (GeoNameDetailed, error) {
	var res GeoNameDetailed

	err := c.apiRequest(
		ctx,
		pathGet,
		req,
		&res,
	)

	return res, err
}
