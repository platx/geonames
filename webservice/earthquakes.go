package webservice

import (
	"context"
	"time"

	"github.com/platx/geonames/value"
)

const pathEarthquakes = "/earthquakesJSON"

type EarthquakesRequest struct {
	// BoundingBox only entries within the box are returned.
	BoundingBox value.BoundingBox `url:",dive"`
	// Date earthquakes older or equal the given date sorted by date,magnitude.
	Date time.Time `url:"date"`
	// MinMagnitude minimal magnitude.
	MinMagnitude float64 `url:"minMagnitude"`
	// MaxRows maximal number of rows returned (default = 10).
	MaxRows uint32 `url:"maxRows"`
}

// Earthquakes returns a list of earthquakes, ordered by magnitude.
func (c *Client) Earthquakes(ctx context.Context, req EarthquakesRequest) ([]Earthquake, error) {
	var res struct {
		Items []Earthquake `json:"earthquakes"`
	}

	if err := c.apiRequest(
		ctx,
		pathEarthquakes,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
