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
func (c *Client) CountryInfo(ctx context.Context, req CountryInfoRequest) ([]CountryDetailed, error) {
	var res struct {
		Items []CountryDetailed `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathCountryInfo,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
