package webservice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/platx/geonames/value"
)

type GeoName struct {
	// ID of record in geonames database
	ID uint64
	// Country information
	Country value.Country
	// AdminSubdivision is the administrative subdivision of a toponym, such as a state or province.
	AdminSubdivision value.AdminDivisions
	// Feature class and code
	Feature value.Feature
	// Position coordinates of the geographical point
	Position value.Position
	// Name is a localized name of geographical point, the preferred name in the language passed
	// in the optional 'lang' parameter or the name that triggered the response in a 'startWith' search.
	Name string
	// ToponymName is the main name of the toponym as displayed on the google maps interface page
	// or in the geoname file in the download. The 'name' attribute is derived from the alternate names.
	ToponymName string
	// Population of the toponym
	Population uint64
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
		Latitude         float64           `json:"lat,string"`
		Longitude        float64           `json:"lng,string"`
		Population       uint64            `json:"population"`
		CountryID        uint64            `json:"countryId,string"`
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
	v.Country.ID = raw.CountryID
	v.Country.Code = raw.CountryCode
	v.Country.Name = raw.CountryName
	v.AdminSubdivision = value.AdminDivisions{
		First:  value.AdminDivision{ID: 0, Code: raw.AdminCode1, Name: raw.AdminName1},
		Second: value.AdminDivision{ID: 0, Code: raw.AdminCode2, Name: raw.AdminName2},
		Third:  value.AdminDivision{ID: 0, Code: raw.AdminCode3, Name: raw.AdminName3},
		Fourth: value.AdminDivision{ID: 0, Code: raw.AdminCode4, Name: raw.AdminName4},
		Fifth:  value.AdminDivision{ID: 0, Code: raw.AdminCode5, Name: raw.AdminName5},
	}
	v.Feature = value.Feature{
		Class:     raw.FeatureClass,
		ClassName: raw.FeatureClassName,
		Code:      raw.FeatureCode,
		CodeName:  raw.FeatureCodeName,
	}
	v.Position = value.Position{
		Latitude:  raw.Latitude,
		Longitude: raw.Longitude,
	}
	v.Name = raw.Name
	v.ToponymName = raw.ToponymName
	v.Population = raw.Population

	return nil
}

type GeoNameNearby struct {
	GeoName

	// Distance in km from the point specified via lat and lng that a result was found
	Distance float64
}

func (v *GeoNameNearby) UnmarshalJSON(data []byte) error {
	var err error

	var parent GeoName

	if err = json.Unmarshal(data, &parent); err != nil {
		return err
	}

	var raw struct {
		Distance float64 `json:"distance,string"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.GeoName = parent
	v.Distance = raw.Distance

	return nil
}

type GeoNameDetailed struct {
	GeoName

	// ContinentCode continent code
	ContinentCode value.ContinentCode
	// ASCIIName name of geographical point in plain ascii characters, varchar(200)
	ASCIIName string
	// AlternateNames alternate names of the geographical point
	AlternateNames []value.AlternateName
	// Timezone the iana timezone id
	Timezone value.Timezone
	// Elevation in meters
	Elevation int32
	// SRTM3 srtm3 digital elevation model
	SRTM3 int32
	// Astergdem aster gdem digital elevation model
	Astergdem int32
	// BoundingBox bounding box surrounding the geographical point
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
		Elevation   int32 `json:"elevation"`
		SRTM3       int32 `json:"srtm3"`
		Astergdem   int32 `json:"astergdem"`
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
	value.Country
	// Continent continent code
	Continent value.Continent
	// Capital administrative capital city of the country
	Capital string
	// Languages spoken in the country
	Languages []string
	// BoundingBox bounding box surrounding the country
	BoundingBox value.BoundingBox
	// IsoAlpha3 3-letter ISO country code
	IsoAlpha3 string
	// IsoNumeric 3-digit ISO country code
	IsoNumeric uint64
	// FipsCode fips code
	FipsCode string
	// Population of the country
	Population int64
	// AreaInSqKm area in square km
	AreaInSqKm float64
	// PostalCodeFormat postal code format
	PostalCodeFormat string
	// CurrencyCode official currency code
	CurrencyCode string
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
		IsoNumeric       uint64              `json:"isoNumeric,string"`
		FipsCode         string              `json:"fipsCode"`
		Capital          string              `json:"capital"`
		Languages        string              `json:"languages"`
		PostalCodeFormat string              `json:"postalCodeFormat"`
		CurrencyCode     string              `json:"currencyCode"`
		Population       int64               `json:"population,string"`
		AreaInSqKm       float64             `json:"areaInSqKm,string"`
		South            float64             `json:"south"`
		North            float64             `json:"north"`
		East             float64             `json:"east"`
		West             float64             `json:"west"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Country.ID = raw.GeoNameID
	v.Country.Code = raw.CountryCode
	v.Country.Name = raw.CountryName
	v.Continent.Code = raw.Continent
	v.Continent.Name = raw.ContinentName
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
	v.IsoNumeric = raw.IsoNumeric
	v.Population = raw.Population
	v.AreaInSqKm = raw.AreaInSqKm

	return nil
}

type CountryNearby struct {
	value.Country
	// Languages spoken in the country
	Languages []string
	// Distance in km from the point specified via lat and lng that a result was found
	Distance float64
}

func (v *CountryNearby) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		Languages   string            `json:"languages"`
		Distance    float64           `json:"distance,string"`
		CountryCode value.CountryCode `json:"countryCode"`
		CountryName string            `json:"countryName"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Code = raw.CountryCode
	v.Name = raw.CountryName
	v.Languages = value.ParseMultipleValues[string](raw.Languages)
	v.Distance = raw.Distance

	return nil
}

type PostalCode struct {
	// Code postal code
	Code string
	// CountryCode ISO3166 2-letter country code
	CountryCode value.CountryCode
	// AdminDivisions administrative divisions
	AdminDivisions value.AdminDivisions
	// PlaceName name of the place
	PlaceName string
	// Position coordinates of the postal code
	Position value.Position
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
	v.AdminDivisions = value.AdminDivisions{
		First:  value.AdminDivision{ID: 0, Code: raw.AdminCode1, Name: raw.AdminName1},
		Second: value.AdminDivision{ID: 0, Code: raw.AdminCode2, Name: raw.AdminName2},
		Third:  value.AdminDivision{ID: 0, Code: raw.AdminCode3, Name: raw.AdminName3},
		Fourth: value.AdminDivision{ID: 0, Code: raw.AdminCode4, Name: raw.AdminName4},
		Fifth:  value.AdminDivision{ID: 0, Code: raw.AdminCode5, Name: raw.AdminName5},
	}
	v.PlaceName = raw.PlaceName
	v.Position = value.Position{
		Latitude:  raw.Latitude,
		Longitude: raw.Longitude,
	}

	return nil
}

type PostalCodeNearby struct {
	PostalCode

	// Distance in km from the point specified via lat and lng that a result was found
	Distance float64
}

func (v *PostalCodeNearby) UnmarshalJSON(data []byte) error {
	var err error

	var parent PostalCode

	if err = json.Unmarshal(data, &parent); err != nil {
		return err
	}

	var raw struct {
		Distance float64 `json:"distance,string"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.PostalCode = parent
	v.Distance = raw.Distance

	return nil
}

type Wikipedia struct {
	// ID of record in geonames database
	ID uint64
	// CountryCode ISO3166 2-letter country code
	CountryCode value.CountryCode
	// Feature feature class
	Feature string
	// Title of the wikipedia article
	Title string
	// Summary of the wikipedia article
	Summary string
	// Position coordinates of the wikipedia article
	Position value.Position
	// Language of the wikipedia article
	Language string
	// ThumbnailURL URL to the thumbnail image
	ThumbnailURL string
	// WikipediaURL URL to the wikipedia article
	WikipediaURL string
	// Rank of the wikipedia article
	Rank uint32
	// Elevation in meters
	Elevation int32
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
		Distance float64 `json:"distance,string"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Wikipedia = parent
	v.Distance = raw.Distance

	return nil
}

type Timezone struct {
	// Name of the timezone
	Name string
	// Country information
	Country value.Country
	// Position coordinates of the timezone
	Position value.Position
	// Time current time in the timezone
	Time time.Time
	// Sunrise time of sunrise
	Sunrise time.Time
	// Sunset time of sunset
	Sunset time.Time
	// GMTOffset offset to GMT in hours
	GMTOffset int
	// DSTOffset daylight saving time offset in hours
	DSTOffset int
	// RawOffset offset to GMT in hours (without DST)
	RawOffset int
}

func (v *Timezone) UnmarshalJSON(data []byte) error {
	const timeFormat = "2006-01-02 15:04"

	var err error

	var raw struct {
		TimezoneID  string            `json:"timezoneId"`
		CountryCode value.CountryCode `json:"countryCode"`
		CountryName string            `json:"countryName"`
		Latitude    float64           `json:"lat"`
		Longitude   float64           `json:"lng"`
		Time        string            `json:"time"`
		Sunset      string            `json:"sunset"`
		Sunrise     string            `json:"sunrise"`
		GMTOffset   int               `json:"gmtOffset"`
		DSTOffset   int               `json:"dstOffset"`
		RawOffset   int               `json:"rawOffset"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Name = raw.TimezoneID
	v.Country.Code = raw.CountryCode
	v.Country.Name = raw.CountryName
	v.Position = value.Position{
		Latitude:  raw.Latitude,
		Longitude: raw.Longitude,
	}
	v.GMTOffset = raw.GMTOffset
	v.DSTOffset = raw.DSTOffset
	v.RawOffset = raw.RawOffset

	var location *time.Location

	if location, err = time.LoadLocation(v.Name); err != nil {
		return fmt.Errorf("load timezone location => %w", err)
	}

	if v.Time, err = time.ParseInLocation(timeFormat, raw.Time, location); err != nil {
		return fmt.Errorf("parse Time => %w", err)
	}

	if v.Sunset, err = time.ParseInLocation(timeFormat, raw.Sunset, location); err != nil {
		return fmt.Errorf("parse Sunset => %w", err)
	}

	if v.Sunrise, err = time.ParseInLocation(timeFormat, raw.Sunrise, location); err != nil {
		return fmt.Errorf("parse Sunrise => %w", err)
	}

	return nil
}

// WeatherObservation represents a weather observation (https://en.wikipedia.org/wiki/METAR).
type WeatherObservation struct {
	// Position coordinates of the weather station
	Position value.Position
	// Observation METAR raw weather observation data
	Observation string
	// ICAO code of the weather station
	ICAO string
	// StationName name of the weather station
	StationName string
	// CloudsCode cloud coverage description
	CloudsCode string
	// CloudsName clouds coverage description
	CloudsName string
	// WeatherCondition weather condition description (if available)
	WeatherCondition string
	// Temperature represents air temperature in Celsius
	Temperature int
	// DewPoint temperature in Celsius
	DewPoint int
	// Humidity percentage
	Humidity int
	// WindDirection wind direction in degrees (0-360)
	WindDirection int
	// WindSpeed wind speed in knots
	WindSpeed int
	// UpdatedAt timestamp in UTC
	UpdatedAt time.Time
}

func (v *WeatherObservation) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		Latitude         float64 `json:"lat"`
		Longitude        float64 `json:"lng"`
		Observation      string  `json:"observation"`
		ICAO             string  `json:"ICAO"`
		StationName      string  `json:"stationName"`
		CloudsCode       string  `json:"cloudsCode"`
		CloudsName       string  `json:"clouds"`
		WeatherCondition string  `json:"weatherCondition"`
		Temperature      int     `json:"temperature,string"`
		DewPoint         int     `json:"dewPoint,string"`
		Humidity         int     `json:"humidity"`
		WindDirection    int     `json:"windDirection"`
		WindSpeed        int     `json:"windSpeed,string"`
		Datetime         string  `json:"datetime"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Position = value.Position{
		Latitude:  raw.Latitude,
		Longitude: raw.Longitude,
	}
	v.Observation = raw.Observation
	v.ICAO = raw.ICAO
	v.StationName = raw.StationName
	v.CloudsCode = raw.CloudsCode
	v.CloudsName = raw.CloudsName
	v.WeatherCondition = raw.WeatherCondition
	v.Temperature = raw.Temperature
	v.DewPoint = raw.DewPoint
	v.Humidity = raw.Humidity
	v.WindDirection = raw.WindDirection
	v.WindSpeed = raw.WindSpeed

	if v.UpdatedAt, err = time.Parse(time.DateTime, raw.Datetime); err != nil {
		return fmt.Errorf("parse UpdatedAt => %w", err)
	}

	return nil
}

// WeatherObservationNearby represents a weather observation (https://en.wikipedia.org/wiki/METAR).
type WeatherObservationNearby struct {
	WeatherObservation
	// CountryCode ISO3166 2-letter country code
	CountryCode value.CountryCode
	// Elevation in meters
	Elevation int32
	// HectoPascAltimeter altimeter pressure
	HectoPascAltimeter int32
}

func (v *WeatherObservationNearby) UnmarshalJSON(data []byte) error {
	var err error

	var parent WeatherObservation

	if err = json.Unmarshal(data, &parent); err != nil {
		return err
	}

	var raw struct {
		CountryCode        value.CountryCode `json:"countryCode"`
		Elevation          int32             `json:"elevation"`
		HectoPascAltimeter int32             `json:"hectoPascAltimeter"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.WeatherObservation = parent
	v.CountryCode = raw.CountryCode
	v.Elevation = raw.Elevation
	v.HectoPascAltimeter = raw.HectoPascAltimeter

	return nil
}

type CountrySubdivision struct {
	value.Country
	// GeoNameID ID of record in geonames database
	GeoNameID uint64
	// Codes administrative level codes
	Codes []value.AdminLevelCode
	// AdminDivision administrative divisions
	AdminDivision value.AdminDivisions
	// Distance in km from the point specified via lat and lng that a result was found
	Distance float64
}

func (v *CountrySubdivision) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		GeoNameID       uint64  `json:"geonameId"`
		CountryCode     string  `json:"countryCode"`
		CountryName     string  `json:"countryName"`
		Admin1GeonameID uint64  `json:"admin1geonameId"`
		AdminCode1      string  `json:"adminCode1"`
		AdminName1      string  `json:"adminName1"`
		Admin2GeonameID uint64  `json:"admin2geonameId"`
		AdminCode2      string  `json:"adminCode2"`
		AdminName2      string  `json:"adminName2"`
		Admin3GeonameID uint64  `json:"admin3geonameId"`
		AdminCode3      string  `json:"adminCode3"`
		AdminName3      string  `json:"adminName3"`
		Admin4GeonameID uint64  `json:"admin4geonameId"`
		AdminCode4      string  `json:"adminCode4"`
		AdminName4      string  `json:"adminName4"`
		Admin5GeonameID uint64  `json:"admin5geonameId"`
		AdminCode5      string  `json:"adminCode5"`
		AdminName5      string  `json:"adminName5"`
		Distance        float64 `json:"distance"`
		Codes           []struct {
			Code  string `json:"code"`
			Level uint8  `json:"level,string"`
			Type  string `json:"type"`
		} `json:"codes"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.GeoNameID = raw.GeoNameID
	v.Country = value.Country{ID: 0, Code: value.CountryCode(raw.CountryCode), Name: raw.CountryName}
	v.AdminDivision.First = value.AdminDivision{ID: raw.Admin1GeonameID, Code: raw.AdminCode1, Name: raw.AdminName1}
	v.AdminDivision.Second = value.AdminDivision{ID: raw.Admin2GeonameID, Code: raw.AdminCode2, Name: raw.AdminName2}
	v.AdminDivision.Third = value.AdminDivision{ID: raw.Admin3GeonameID, Code: raw.AdminCode3, Name: raw.AdminName3}
	v.AdminDivision.Fourth = value.AdminDivision{ID: raw.Admin4GeonameID, Code: raw.AdminCode4, Name: raw.AdminName4}
	v.AdminDivision.Fifth = value.AdminDivision{ID: raw.Admin5GeonameID, Code: raw.AdminCode5, Name: raw.AdminName5}
	v.Distance = raw.Distance
	v.Codes = make([]value.AdminLevelCode, 0, len(raw.Codes))

	for _, code := range raw.Codes {
		v.Codes = append(v.Codes, value.AdminLevelCode{Code: code.Code, Level: code.Level, Type: code.Type})
	}

	return nil
}

type Earthquake struct {
	// ID Equivalent identification (EqID)
	ID string
	// Position coordinates of the earthquake
	Position value.Position
	// Depth of the earthquake in km
	Depth float64
	// Source of the earthquake
	Source string
	// Magnitude of the earthquake
	Magnitude float64
	// Time of the earthquake
	Time time.Time
}

func (v *Earthquake) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		Datetime  string  `json:"datetime"`
		Depth     float64 `json:"depth"`
		Lng       float64 `json:"lng"`
		Src       string  `json:"src"`
		Eqid      string  `json:"eqid"`
		Magnitude float64 `json:"magnitude"`
		Lat       float64 `json:"lat"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.ID = raw.Eqid
	v.Position = value.Position{
		Latitude:  raw.Lat,
		Longitude: raw.Lng,
	}
	v.Depth = raw.Depth
	v.Source = raw.Src
	v.Magnitude = raw.Magnitude

	if v.Time, err = time.Parse(time.DateTime, raw.Datetime); err != nil {
		return fmt.Errorf("parse Time => %w", err)
	}

	return nil
}

type Ocean struct {
	// ID of record in geonames database
	ID uint64
	// Distance in km from the point specified via lat and lng that a result was found
	Distance float64
	// Name of the ocean/sea
	Name string
}

func (v *Ocean) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		ID       uint64  `json:"geonameId"`
		Distance float64 `json:"distance,string"`
		Name     string  `json:"name"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.ID = raw.ID
	v.Distance = raw.Distance
	v.Name = raw.Name

	return nil
}

type Address struct {
	// Position coordinates of the address
	Position value.Position
	// CountryCode ISO3166 2-letter country code
	CountryCode value.CountryCode
	// AdminDivision administrative divisions
	AdminDivision value.AdminDivisions
	// SourceID source of the address
	SourceID string
	// PostalCode postal code
	PostalCode string
	// Locality city or town
	Locality string
	// Street name
	Street string
	// HouseNumber house number
	HouseNumber string
}

func (v *Address) UnmarshalJSON(data []byte) error {
	var err error

	var raw struct {
		Lat         float64           `json:"lat,string"`
		Lng         float64           `json:"lng,string"`
		CountryCode value.CountryCode `json:"countryCode"`
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
		SourceID    string            `json:"sourceId"`
		PostalCode  string            `json:"postalcode"`
		Locality    string            `json:"locality"`
		Street      string            `json:"street"`
		HouseNumber string            `json:"houseNumber"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Position = value.Position{
		Latitude:  raw.Lat,
		Longitude: raw.Lng,
	}
	v.CountryCode = raw.CountryCode
	v.AdminDivision = value.AdminDivisions{
		First:  value.AdminDivision{ID: 0, Code: raw.AdminCode1, Name: raw.AdminName1},
		Second: value.AdminDivision{ID: 0, Code: raw.AdminCode2, Name: raw.AdminName2},
		Third:  value.AdminDivision{ID: 0, Code: raw.AdminCode3, Name: raw.AdminName3},
		Fourth: value.AdminDivision{ID: 0, Code: raw.AdminCode4, Name: raw.AdminName4},
		Fifth:  value.AdminDivision{ID: 0, Code: raw.AdminCode5, Name: raw.AdminName5},
	}
	v.SourceID = raw.SourceID
	v.PostalCode = raw.PostalCode
	v.Locality = raw.Locality
	v.Street = raw.Street
	v.HouseNumber = raw.HouseNumber

	return nil
}

type AddressNearby struct {
	Address

	// Distance in km from the point specified via lat and lng that a result was found
	Distance float64
}

func (v *AddressNearby) UnmarshalJSON(data []byte) error {
	var err error

	var parent Address

	if err = json.Unmarshal(data, &parent); err != nil {
		return err
	}

	var raw struct {
		Distance float64 `json:"distance,string"`
	}

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Address = parent
	v.Distance = raw.Distance

	return nil
}
