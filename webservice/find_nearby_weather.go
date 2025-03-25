package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathFindNearbyWeather = "/findNearByWeatherJSON"

type FindNearbyWeatherRequest struct {
	// Position the service will return the station closest to this given point (reverse geocoding).
	Position value.Position `url:",dive"`
	// Radius search radius, only weather stations within this radius are considered. Default is about 100km.
	Radius int32 `url:"radius"`
}

// FindNearbyWeather returns a weather station with the most recent weather observation.
// [More info]: https://www.geonames.org/export/JSON-webservices.html#findNearByWeatherJSON
func (c *Client) FindNearbyWeather(
	ctx context.Context,
	req FindNearbyWeatherRequest,
) (WeatherObservationNearby, error) {
	var res struct {
		Data WeatherObservationNearby `json:"weatherObservation"`
	}

	err := c.apiRequest(
		ctx,
		pathFindNearbyWeather,
		req,
		&res,
	)

	return res.Data, err
}
