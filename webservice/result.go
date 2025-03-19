package webservice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/platx/geonames/value"
)

type GeoName struct {
	// ID of record in geonames database
	ID               uint64
	Country          value.Country
	AdminSubdivision value.AdminDivisions
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
	v.Country.Code = raw.CountryCode
	v.Country.Name = raw.CountryName
	v.AdminSubdivision = value.AdminDivisions{
		First:  value.AdminDivision{ID: 0, Code: raw.AdminCode1, Name: raw.AdminName1},
		Second: value.AdminDivision{ID: 0, Code: raw.AdminCode2, Name: raw.AdminName2},
		Third:  value.AdminDivision{ID: 0, Code: raw.AdminCode3, Name: raw.AdminName3},
		Fourth: value.AdminDivision{ID: 0, Code: raw.AdminCode4, Name: raw.AdminName4},
		Fifth:  value.AdminDivision{ID: 0, Code: raw.AdminCode5, Name: raw.AdminName5},
	}
	v.FeatureClass = raw.FeatureClass
	v.FeatureClassName = raw.FeatureClassName
	v.FeatureCode = raw.FeatureCode
	v.FeatureCodeName = raw.FeatureCodeName
	v.Name = raw.Name
	v.ToponymName = raw.ToponymName
	v.Population = raw.Population

	if v.Country.ID, err = value.ParseUint64(raw.CountryID); err != nil {
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
	value.Country
	Continent        value.Continent
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
	value.Country
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
	AdminDivisions value.AdminDivisions
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

type Timezone struct {
	Name      string
	Country   value.Country
	Position  value.Position
	Time      time.Time
	Sunrise   time.Time
	Sunset    time.Time
	GMTOffset int
	DSTOffset int
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
	GeoNameID uint64
	value.Country
	Codes         []value.AdminLevelCode
	AdminDivision value.AdminDivisions
	Distance      float64
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
	ID        string
	Position  value.Position
	Depth     float64
	Source    string
	Magnitude float64
	Time      time.Time
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
	ID       uint64
	Distance float64
	Name     string
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
