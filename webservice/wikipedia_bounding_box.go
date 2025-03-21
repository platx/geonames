package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathWikipediaBoundingBox = "/wikipediaBoundingBoxJSON"

type WikipediaBoundingBoxRequest struct {
	// BoundingBox only entries within the box are returned.
	BoundingBox value.BoundingBox `url:",dive"`
	// Language language code, supported languages are de,en,es,fr,it,nl,pl,pt,ru,zh (default = en)
	Language string `url:"lang"`
	// MaxRows maximal number of rows returned (default = 10)
	MaxRows uint32 `url:"maxRows"`
}

// WikipediaBoundingBox returns the wikipedia entries within the bounding box.
// [More info]: https://www.geonames.org/export/wikipedia-webservice.html#wikipediaBoundingBox
func (c *Client) WikipediaBoundingBox(
	ctx context.Context,
	req WikipediaBoundingBoxRequest,
) ([]Wikipedia, error) {
	var res struct {
		Items []Wikipedia `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathWikipediaBoundingBox,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
