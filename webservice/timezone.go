package webservice

import (
	"context"
	"time"

	"github.com/platx/geonames/value"
)

const pathTimezone = "/timezoneJSON"

type TimezoneRequest struct {
	// Position the latitude and longitude of the search location
	Position value.Position `url:",dive"`
	// Radius is a buffer in km for closest timezone in coastal areas.
	Radius uint32 `url:"radius"`
	// Language for countryName
	Language string `url:"lang"`
	// Date for sunrise/sunset.
	Date time.Time `url:"date"`
}

// Timezone returns the closest timezone information for lat/lng.
func (c *Client) Timezone(ctx context.Context, req TimezoneRequest) (Timezone, error) {
	var res Timezone

	if err := c.apiRequest(
		ctx,
		pathTimezone,
		req,
		&res,
	); err != nil {
		return Timezone{}, err
	}

	return res, nil
}
