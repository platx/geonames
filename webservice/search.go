package webservice

import (
	"context"

	"github.com/platx/geonames/value"
)

const pathGeoNameSearch = "/searchJSON"

type SearchRequest struct {
	// Query search over all attributes of a place : place name, country name, continent, admin codes,...
	Query string `url:"q"`
	// Name place name
	Name string `url:"name"`
	// NameEquals exact place name
	NameEquals string `url:"name_equals"`
	// NameStartsWith place name starts with given characters
	NameStartsWith string `url:"name_startsWith"`
	// MaxRows the maximal number of rows in the document returned by the service.
	// Default is 100, the maximal allowed value is 1000.
	MaxRows uint32 `url:"maxRows"`
	// StartRow Used for paging results. If you want to get results 30 to 40, use startRow=30 and maxRows=10.
	// Default is 0, the maximal allowed value is 5000 for the free services and 25000 for the premium services.
	StartRow uint32 `url:"startRow"`
	// Country default is all countries
	Country []value.CountryCode `url:"country"`
	// CountryBias records from the countryBias are listed first
	CountryBias value.CountryCode `url:"countryBias"`
	// ContinentCode restricts the search for toponym of the given continent.
	ContinentCode value.ContinentCode `url:"continentCode"`
	// AdminCode code of administrative subdivision
	AdminCode value.AdminCode `url:",dive"`
	// FeatureClass default is all feature classes
	FeatureClass []string `url:"featureClass"`
	// FeatureCode default is all feature codes
	FeatureCode []string `url:"featureCode"`
	// Cities used to categorize the populated places into three groups according to size/relevance
	Cities value.Cities `url:"cities"`
	// Language place name and country name will be returned in the specified language. Default is English.
	// With the pseudo language code 'local' the local language will be returned.
	// Feature classes and codes are only available in English and Bulgarian.
	Language string `url:"lang"`
	// SearchLanguage in combination with the name parameter, the search will only consider names
	// in the specified language. Used for instance to query for IATA airport codes.
	SearchLanguage string `url:"searchLanguage"`
	// NameRequired At least one of the search term needs to be part of the place name.
	// Example : A normal search for Berlin will return all places within the state of Berlin.
	// If we only want to find places with 'Berlin' in the name we set the parameter isNameRequired to 'true'.
	// The difference to the name_equals parameter is that this will allow searches for 'Berlin, Germany' as
	// only one search term needs to be part of the name.
	NameRequired bool `url:"isNameRequired"`
	// Tag search for toponyms tagged with the specified tag
	Tag string `url:"tag"`
	// Operator default is 'AND', with the operator 'OR' not all search terms need to be matched by the response
	Operator value.Operator `url:"operator"`
	// Fuzzy default is '1', defines the fuzziness of the search terms. float between 0 and 1.
	// The search term is only applied to the name attribute.
	Fuzzy float32 `url:"fuzzy"`
	// BoundingBox only features within the box are returned.
	BoundingBox value.BoundingBox `url:",dive"`
	// OrderBy in combination with the name_startsWith, if set to 'relevance' than the result is sorted by relevance.
	OrderBy value.OrderBy `url:"orderby"`
}

// Search returns the names found for the searchterm, the search is using an AND operator.
// [More info]: https://www.geonames.org/export/geonames-search.html
func (c *Client) Search(ctx context.Context, req SearchRequest) ([]GeoName, error) {
	var res struct {
		Items []GeoName `json:"geonames"`
	}

	if err := c.apiRequest(
		ctx,
		pathGeoNameSearch,
		req,
		&res,
	); err != nil {
		return nil, err
	}

	return res.Items, nil
}
