package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathWeather = "/weatherJSON"

type WeatherRequest struct {
	// BoundingBox only entries within the box are returned.
	BoundingBox value.BoundingBox `url:",dive"`
	// Language language code, supported languages are de,en,es,fr,it,nl,pl,pt,ru,zh (default = en)
	Language string `url:"lang"`
	// MaxRows maximal number of rows returned (default = 10)
	MaxRows uint32 `url:"maxRows"`
}

// Weather returns a list of weather stations with the most recent weather observation.
func (c *Client) Weather(ctx context.Context, req WeatherRequest) ([]WeatherObservation, error) {
	var res struct {
		Items []WeatherObservation `json:"weatherObservations"`
	}

	if err := c.apiRequest(
		ctx,
		pathWeather,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
