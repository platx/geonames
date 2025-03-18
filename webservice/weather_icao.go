package webservice

import (
	"context"
)

const pathWeatherICAO = "/weatherIcaoJSON"

type WeatherICAORequest struct {
	// ICAO International Civil Aviation Organization code
	ICAO string `url:"ICAO"`
	// Language language code, supported languages are de,en,es,fr,it,nl,pl,pt,ru,zh (default = en)
	Language string `url:"lang"`
}

// WeatherICAO returns the weather station and the most recent weather observation for the ICAO code.
func (c *Client) WeatherICAO(
	ctx context.Context,
	req WeatherICAORequest,
) (WeatherObservation, error) {
	var res WeatherObservation

	if err := c.apiRequest(
		ctx,
		pathWeatherICAO,
		req,
		&res,
	); err != nil {
		return WeatherObservation{}, err
	}

	return res, nil
}
