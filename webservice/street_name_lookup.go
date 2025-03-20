package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathStreetNameLookup = "/streetNameLookupJSON"

type StreetNameLookupRequest struct {
	// Query query term
	Query string `url:"q"`
	// CountryCode default is all countries
	CountryCode value.CountryCode `url:"country"`
	// PostalCode postal code
	PostalCode string `url:"postalcode"`
	// AdminCode code of administrative subdivision
	AdminCode value.AdminCode `url:",dive"`
	// UniqueStreetName duplicate street names are filtered, the placename will be empty
	// when the same street name occurs at different places.
	UniqueStreetName bool `url:"isUniqueStreetName"`
}

// StreetNameLookup returns a list of street names starting with the query term.
// Supported countries: AT,AU,AX,CC,CH,CL,CX,CZ,DK,EE,ES,FI,FR,GF,GP,HK,IS,LU,MQ,NF,NL,NO,PL,PR,PT,RE,SG,SI,SJ,SK,US,YT.
func (c *Client) StreetNameLookup(ctx context.Context, req StreetNameLookupRequest) ([]Address, error) {
	var res struct {
		Address []Address `json:"address"`
	}

	err := c.apiRequest(
		ctx,
		pathStreetNameLookup,
		req,
		&res,
	)

	return res.Address, err
}
