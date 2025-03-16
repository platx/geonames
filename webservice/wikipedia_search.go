package webservice

import (
	"context"
)

const pathWikipediaSearch = "/wikipediaSearchJSON"

type WikipediaSearchRequest struct {
	// Query place name
	Query string `url:"q"`
	// Title search in the wikipedia title
	Title string `url:"title"`
	// Language language code, supported languages are de,en,es,fr,it,nl,pl,pt,ru,zh (default = en)
	Language string `url:"lang"`
	// MaxRows maximal number of rows returned (default = 10)
	MaxRows uint32 `url:"maxRows"`
}

// WikipediaSearch returns the wikipedia entries found for the searchterm.
func (c *Client) WikipediaSearch(ctx context.Context, req WikipediaSearchRequest) ([]Wikipedia, error) {
	var res struct {
		Items []Wikipedia `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathWikipediaSearch,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
