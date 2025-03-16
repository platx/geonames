package webservice

import (
	"encoding/json"
	"fmt"

	"github.com/platx/geonames/value"
)

type GeoName struct {
	// ID of record in geonames database
	ID               uint64
	CountryID        uint64
	CountryCode      value.CountryCode
	CountryName      string
	AdminSubdivision value.AdminSubdivisions
	FeatureClass     string
	FeatureClassName string
	FeatureCode      string
	FeatureCodeName  string
	// Name is a localized name of geographical point, the preferred name in the language passed
	// in the optional 'lang' parameter or the name that triggered the response in a 'startWith' search.
	Name string
	// ToponymName is the main name of the toponym as displayed on the google maps interface page
	// or in the geoname file in the download. The 'name' attribute is derived from the alternate names.
	ToponymName string
	Position    value.Position
	Population  uint64
}

func (v *GeoName) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		GeoNameID        uint64            `json:"geonameId"`
		Name             string            `json:"name"`
		ToponymName      string            `json:"toponymName"`
		FeatureClass     string            `json:"fcl"`
		FeatureClassName string            `json:"fclName"`
		FeatureCode      string            `json:"fcode"`
		FeatureCodeName  string            `json:"fcodeName"`
		Latitude         string            `json:"lat"`
		Longitude        string            `json:"lng"`
		Population       uint64            `json:"population"`
		CountryID        string            `json:"countryId"`
		CountryCode      value.CountryCode `json:"countryCode"`
		CountryName      string            `json:"countryName"`
		AdminCode1       string            `json:"adminCode1"`
		AdminName1       string            `json:"adminName1"`
		AdminCode2       string            `json:"adminCode2"`
		AdminName2       string            `json:"adminName2"`
		AdminCode3       string            `json:"adminCode3"`
		AdminName3       string            `json:"adminName3"`
		AdminCode4       string            `json:"adminCode4"`
		AdminName4       string            `json:"adminName4"`
		AdminCode5       string            `json:"adminCode5"`
		AdminName5       string            `json:"adminName5"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.ID = raw.GeoNameID
	v.CountryCode = raw.CountryCode
	v.CountryName = raw.CountryName
	v.AdminSubdivision = value.AdminSubdivisions{
		First:  value.AdminSubdivision{Code: raw.AdminCode1, Name: raw.AdminName1},
		Second: value.AdminSubdivision{Code: raw.AdminCode2, Name: raw.AdminName2},
		Third:  value.AdminSubdivision{Code: raw.AdminCode3, Name: raw.AdminName3},
		Fourth: value.AdminSubdivision{Code: raw.AdminCode4, Name: raw.AdminName4},
		Fifth:  value.AdminSubdivision{Code: raw.AdminCode5, Name: raw.AdminName5},
	}
	v.FeatureClass = raw.FeatureClass
	v.FeatureClassName = raw.FeatureClassName
	v.FeatureCode = raw.FeatureCode
	v.FeatureCodeName = raw.FeatureCodeName
	v.Name = raw.Name
	v.ToponymName = raw.ToponymName
	v.Population = raw.Population

	if v.CountryID, err = value.ParseUint64(raw.CountryID); err != nil {
		return fmt.Errorf("parse CountryID => %w", err)
	}

	if v.Position, err = value.ParsePosition(raw.Latitude, raw.Longitude); err != nil {
		return fmt.Errorf("parse Position => %w", err)
	}

	return nil
}

type GeoNameNearby struct {
	GeoName

	// Distance in km from the point specified via lat and lng that a result was found
	Distance float64
}

func (v *GeoNameNearby) UnmarshalJSON(data []byte) error {
	var err error

	var geoName GeoName

	if err = json.Unmarshal(data, &geoName); err != nil {
		return err
	}

	var raw struct {
		Distance string `json:"distance"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.GeoName = geoName

	if v.Distance, err = value.ParseFloat64(raw.Distance); err != nil {
		return fmt.Errorf("parse Distance => %w", err)
	}

	return nil
}

type GeoNameDetailed struct {
	GeoName

	ContinentCode value.ContinentCode
	// ASCIIName name of geographical point in plain ascii characters, varchar(200)
	ASCIIName      string
	AlternateNames []value.AlternateName
	// Timezone the iana timezone id
	Timezone value.Timezone
	// Elevation in meters
	Elevation   int32
	SRTM3       uint64
	Astergdem   uint64
	BoundingBox value.BoundingBox
}

func (v *GeoNameDetailed) UnmarshalJSON(data []byte) error {
	var err error

	var parent GeoName
	if err = json.Unmarshal(data, &parent); err != nil {
		return err
	}

	var raw struct {
		ContinentCode  value.ContinentCode `json:"continentCode"`
		ASCIIName      string              `json:"asciiName"`
		AlternateNames []struct {
			Language string `json:"lang"`
			Value    string `json:"name"`
		} `json:"alternateNames"`
		Timezone struct {
			Name      string  `json:"timeZoneId"`
			GMTOffset float64 `json:"gmtOffset"`
			DSTOffset float64 `json:"dstOffset"`
		} `json:"timezone"`
		Elevation   int32  `json:"elevation"`
		SRTM3       uint64 `json:"srtm3"`
		Astergdem   uint64 `json:"astergdem"`
		BoundingBox struct {
			West  float64 `json:"west"`
			East  float64 `json:"east"`
			North float64 `json:"north"`
			South float64 `json:"south"`
		} `json:"bbox"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.GeoName = parent
	v.ContinentCode = raw.ContinentCode
	v.ASCIIName = raw.ASCIIName

	v.AlternateNames = make([]value.AlternateName, 0, len(raw.AlternateNames))
	for _, name := range raw.AlternateNames {
		v.AlternateNames = append(v.AlternateNames, value.AlternateName{Language: name.Language, Value: name.Value})
	}

	v.Timezone = value.Timezone{
		Name:      raw.Timezone.Name,
		GMTOffset: raw.Timezone.GMTOffset,
		DSTOffset: raw.Timezone.DSTOffset,
	}
	v.Elevation = raw.Elevation
	v.SRTM3 = raw.SRTM3
	v.Astergdem = raw.Astergdem
	v.BoundingBox = value.BoundingBox{
		East:  raw.BoundingBox.East,
		West:  raw.BoundingBox.West,
		North: raw.BoundingBox.North,
		South: raw.BoundingBox.South,
	}

	return nil
}

type CountryDetailed struct {
	ID               uint64
	Code             value.CountryCode
	Name             string
	ContinentCode    value.ContinentCode
	ContinentName    string
	Capital          string
	Languages        []string
	BoundingBox      value.BoundingBox
	IsoAlpha3        string
	IsoNumeric       uint64
	FipsCode         string
	Population       int64
	AreaInSqKm       float64
	PostalCodeFormat string
	CurrencyCode     string
}

func (v *CountryDetailed) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		GeoNameID        uint64              `json:"geonameId"`
		Continent        value.ContinentCode `json:"continent"`
		CountryCode      value.CountryCode   `json:"countryCode"`
		ContinentName    string              `json:"continentName"`
		CountryName      string              `json:"countryName"`
		IsoAlpha3        string              `json:"isoAlpha3"`
		IsoNumeric       string              `json:"isoNumeric"`
		FipsCode         string              `json:"fipsCode"`
		Capital          string              `json:"capital"`
		Languages        string              `json:"languages"`
		PostalCodeFormat string              `json:"postalCodeFormat"`
		CurrencyCode     string              `json:"currencyCode"`
		Population       string              `json:"population"`
		AreaInSqKm       string              `json:"areaInSqKm"`
		South            float64             `json:"south"`
		North            float64             `json:"north"`
		East             float64             `json:"east"`
		West             float64             `json:"west"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.ID = raw.GeoNameID
	v.Code = raw.CountryCode
	v.Name = raw.CountryName
	v.ContinentCode = raw.Continent
	v.ContinentName = raw.ContinentName
	v.Capital = raw.Capital
	v.Languages = value.ParseMultipleValues[string](raw.Languages)
	v.BoundingBox = value.BoundingBox{
		East:  raw.East,
		West:  raw.West,
		North: raw.North,
		South: raw.South,
	}
	v.IsoAlpha3 = raw.IsoAlpha3
	v.FipsCode = raw.FipsCode
	v.PostalCodeFormat = raw.PostalCodeFormat
	v.CurrencyCode = raw.CurrencyCode

	if v.IsoNumeric, err = value.ParseUint64(raw.IsoNumeric); err != nil {
		return fmt.Errorf("parse IsoNumeric => %w", err)
	}

	if v.Population, err = value.ParseInt64(raw.Population); err != nil {
		return fmt.Errorf("parse Population => %w", err)
	}

	if v.AreaInSqKm, err = value.ParseFloat64(raw.AreaInSqKm); err != nil {
		return fmt.Errorf("parse AreaInSqKm => %w", err)
	}

	return nil
}

type CountryNearby struct {
	Code      value.CountryCode
	Name      string
	Languages []string
	Distance  float64
}

func (v *CountryNearby) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		Languages   string            `json:"languages"`
		Distance    string            `json:"distance"`
		CountryCode value.CountryCode `json:"countryCode"`
		CountryName string            `json:"countryName"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Code = raw.CountryCode
	v.Name = raw.CountryName
	v.Languages = value.ParseMultipleValues[string](raw.Languages)

	if v.Distance, err = value.ParseFloat64(raw.Distance); err != nil {
		return fmt.Errorf("parse Distance => %w", err)
	}

	return nil
}

type PostalCode struct {
	Code           string
	CountryCode    value.CountryCode
	AdminDivisions value.AdminSubdivisions
	PlaceName      string
	Position       value.Position
}

func (v *PostalCode) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		CountryCode value.CountryCode `json:"countryCode"`
		PostalCode  string            `json:"postalCode"`
		PlaceName   string            `json:"placeName"`
		Latitude    float64           `json:"lat"`
		Longitude   float64           `json:"lng"`
		AdminCode1  string            `json:"adminCode1"`
		AdminName1  string            `json:"adminName1"`
		AdminCode2  string            `json:"adminCode2"`
		AdminName2  string            `json:"adminName2"`
		AdminCode3  string            `json:"adminCode3"`
		AdminName3  string            `json:"adminName3"`
		AdminCode4  string            `json:"adminCode4"`
		AdminName4  string            `json:"adminName4"`
		AdminCode5  string            `json:"adminCode5"`
		AdminName5  string            `json:"adminName5"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Code = raw.PostalCode
	v.CountryCode = raw.CountryCode
	v.AdminDivisions = value.AdminSubdivisions{
		First: value.AdminSubdivision{
			Code: raw.AdminCode1,
			Name: raw.AdminName1,
		},
		Second: value.AdminSubdivision{
			Code: raw.AdminCode2,
			Name: raw.AdminName2,
		},
		Third: value.AdminSubdivision{
			Code: raw.AdminCode3,
			Name: raw.AdminName3,
		},
		Fourth: value.AdminSubdivision{
			Code: raw.AdminCode4,
			Name: raw.AdminName4,
		},
		Fifth: value.AdminSubdivision{
			Code: raw.AdminCode5,
			Name: raw.AdminName5,
		},
	}
	v.PlaceName = raw.PlaceName
	v.Position = value.Position{
		Latitude:  raw.Latitude,
		Longitude: raw.Longitude,
	}

	return nil
}

type Wikipedia struct {
	// ID of record in geonames database
	ID           uint64
	CountryCode  value.CountryCode
	Feature      string
	Title        string
	Summary      string
	Position     value.Position
	Language     string
	ThumbnailURL string
	WikipediaURL string
	Rank         uint32
	Elevation    int32
}

func (v *Wikipedia) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		GeoNameID    uint64            `json:"geonameId"`
		CountryCode  value.CountryCode `json:"countryCode"`
		Title        string            `json:"title"`
		Summary      string            `json:"summary"`
		Latitude     float64           `json:"lat"`
		Longitude    float64           `json:"lng"`
		Language     string            `json:"lang"`
		WikipediaURL string            `json:"wikipediaUrl"`
		ThumbnailImg string            `json:"thumbnailImg"`
		Rank         uint32            `json:"rank"`
		Elevation    int32             `json:"elevation"`
		Feature      string            `json:"feature"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.ID = raw.GeoNameID
	v.CountryCode = raw.CountryCode
	v.Feature = raw.Feature
	v.Title = raw.Title
	v.Summary = raw.Summary
	v.Position = value.Position{
		Latitude:  raw.Latitude,
		Longitude: raw.Longitude,
	}
	v.Language = raw.Language
	v.ThumbnailURL = raw.ThumbnailImg
	v.WikipediaURL = raw.WikipediaURL
	v.Rank = raw.Rank
	v.Elevation = raw.Elevation

	return nil
}

type WikipediaNearby struct {
	Wikipedia

	// Distance in km from the point specified via lat and lng that a result was found
	Distance float64
}

func (v *WikipediaNearby) UnmarshalJSON(data []byte) error {
	var err error

	var parent Wikipedia

	if err = json.Unmarshal(data, &parent); err != nil {
		return err
	}

	var raw struct {
		Distance string `json:"distance"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Wikipedia = parent

	if v.Distance, err = value.ParseFloat64(raw.Distance); err != nil {
		return fmt.Errorf("parse Distance => %w", err)
	}

	return nil
}
