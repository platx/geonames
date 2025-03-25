package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathCountryInfo = "/countryInfoJSON"

type CountryInfoRequest struct {
	// Country default is all countries
	Country []value.CountryCode `url:"country"`
	// Language ISO-639-1 language code (en,de,fr,it,es,...) (default = english)
	Language string `url:"lang"`
}

// CountryInfo Country information : Capital, Population, Area in square km,
// Bounding Box of mainland (excluding offshore islands).
// [More info]: https://www.geonames.org/export/web-services.html#countryInfo
func (c *Client) CountryInfo(ctx context.Context, req CountryInfoRequest) ([]CountryDetailed, error) {
	var res struct {
		Items []CountryDetailed `json:"geonames"`
	}

	err := c.apiRequest(
		ctx,
		pathCountryInfo,
		req,
		&res,
	)

	return res.Items, err
}
