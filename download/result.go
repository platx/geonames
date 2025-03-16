package download

import (
	"errors"
	"fmt"
	"time"

	"github.com/platx/geonames/value"
)

var ErrInvalidRowLength = errors.New("invalid row length")

type GeoName struct {
	// ID of record in geonames database
	ID uint64
	// Name of geographical point (utf8), max length is 200 characters
	Name string
	// Name of geographical point in plain ascii characters, max length is 200 characters
	ASCIIName string
	// ascii names automatically transliterated, convenience attribute from alternatename table
	AlternateNames []string
	// Position represents the latitude and longitude of the toponym
	Position value.Position
	// FeatureClass see http://www.geonames.org/export/codes.html
	FeatureClass string
	// FeatureCode see http://www.geonames.org/export/codes.html
	FeatureCode string
	// CountryCode ISO-3166 2-letter country code
	CountryCode value.CountryCode
	// AlternateCountryCodes alternate country codes, ISO-3166 2-letter country code
	AlternateCountryCodes []value.CountryCode
	// AdminCode represents codes of administrative subdivision
	AdminCode value.AdminCode
	// Population represents the population of the toponym
	Population int64
	// Elevation in meters
	Elevation int64
	// DigitalElevationModel represents the digital elevation model, srtm3 or gtopo30,
	// average elevation of 3''x3'' (ca 90mx90m) or 30''x30'' (ca 900mx900m) area in meters.
	// srtm processed by cgiar/ciat.
	DigitalElevationModel int64
	// Timezone the iana timezone id (see file timeZone.txt), max length is 40 characters
	Timezone string
	// ModificationDate date of last modification
	ModificationDate time.Time
}

func (v *GeoName) UnmarshalRow(row []string) error {
	const columns = 19

	var err error
	if err = checkColumns(row, columns); err != nil {
		return err
	}

	if v.ID, err = value.ParseUint64(row[0]); err != nil {
		return fmt.Errorf("parse ID => %w", err)
	}

	if v.Position, err = value.ParsePosition(row[4], row[5]); err != nil {
		return fmt.Errorf("parse Position => %w", err)
	}

	if v.Population, err = value.ParseInt64(row[14]); err != nil {
		return fmt.Errorf("parse Population => %w", err)
	}

	if v.Elevation, err = value.ParseInt64(row[15]); err != nil {
		return fmt.Errorf("parse Elevation => %w", err)
	}

	if v.DigitalElevationModel, err = value.ParseInt64(row[16]); err != nil {
		return fmt.Errorf("parse DigitalElevationModel => %w", err)
	}

	if v.ModificationDate, err = value.ParseDate(row[18]); err != nil {
		return fmt.Errorf("parse ModificationDate => %w", err)
	}

	v.Name = row[1]
	v.ASCIIName = row[2]
	v.AlternateNames = value.ParseMultipleValues[string](row[3])
	v.FeatureClass = row[6]
	v.FeatureCode = row[7]
	v.CountryCode = value.CountryCode(row[8])
	v.AlternateCountryCodes = value.ParseMultipleValues[value.CountryCode](row[9])
	v.AdminCode.First = row[10]
	v.AdminCode.Second = row[11]
	v.AdminCode.Third = row[12]
	v.AdminCode.Fourth = row[13]
	v.Timezone = row[17]

	return nil
}

type AlternateName struct {
	// ID the id of this alternate name
	ID uint64
	// GeoNameID referring to GeoName.ID
	GeoNameID uint64
	// Language iso 639 language code 2- or 3-characters, optionally followed by a hyphen and a country code
	// for country specific variants (ex:zh-CN) or by a variant name (ex: zh-Hant);
	// 4-characters 'post' for postal codes and 'iata', 'icao' and 'faac' for airport codes,
	// fr_1793 for French Revolution names,  abbr for abbreviation, link to a website (mostly to wikipedia),
	// wkdt for the wikidataid. Max length is 7 characters.
	Language string
	// Value alternate name or name variant, max length is 400 characters.
	Value string
	// Preferred whether this is an official/preferred name
	Preferred bool
	// Short whether this is a short name like 'California' for 'State of California'
	Short bool
	// Colloquial whether this is a colloquial or slang term
	Colloquial bool
	// Historic whether this is historic and was used in the past. Example 'Bombay' for 'Mumbai'.
	Historic bool
	// From period when the name was used
	From string
	// To period when the name was used
	To string
}

func (v *AlternateName) UnmarshalRow(row []string) error {
	const columns = 10

	var err error
	if err = checkColumns(row, columns); err != nil {
		return err
	}

	if v.ID, err = value.ParseUint64(row[0]); err != nil {
		return fmt.Errorf("parse ID => %w", err)
	}

	if v.GeoNameID, err = value.ParseUint64(row[1]); err != nil {
		return fmt.Errorf("parse GeoNameID => %w", err)
	}

	v.Language = row[2]
	v.Value = row[3]
	v.Preferred = value.ParseBool(row[4])
	v.Short = value.ParseBool(row[5])
	v.Colloquial = value.ParseBool(row[6])
	v.Historic = value.ParseBool(row[7])
	v.From = row[8]
	v.To = row[9]

	return nil
}

type Country struct {
	ID                 uint64
	Code               value.CountryCode
	Name               string
	ContinentCode      value.ContinentCode
	Domain             string
	Capital            string
	Languages          []string
	IsoAlpha3          string
	IsoNumeric         uint64
	FipsCode           string
	Population         int64
	AreaInSqKm         float64
	PostalCodeFormat   string
	PostalCodeRegex    string
	CurrencyCode       string
	CurrencyName       string
	Phone              string
	Neighbours         []value.CountryCode
	EquivalentFipsCode string
}

func (v *Country) UnmarshalRow(row []string) error {
	const columns = 19

	var err error
	if err = checkColumns(row, columns); err != nil {
		return err
	}

	if v.IsoNumeric, err = value.ParseUint64(row[2]); err != nil {
		return fmt.Errorf("parse IsoNumeric => %w", err)
	}

	if v.AreaInSqKm, err = value.ParseFloat64(row[6]); err != nil {
		return fmt.Errorf("parse AreaInSqKm => %w", err)
	}

	if v.Population, err = value.ParseInt64(row[7]); err != nil {
		return fmt.Errorf("parse Population => %w", err)
	}

	if v.ID, err = value.ParseUint64(row[16]); err != nil {
		return fmt.Errorf("parse ID => %w", err)
	}

	v.Code = value.CountryCode(row[0])
	v.IsoAlpha3 = row[1]
	v.FipsCode = row[3]
	v.Name = row[4]
	v.Capital = row[5]
	v.ContinentCode = value.ContinentCode(row[8])
	v.Domain = row[9]
	v.CurrencyCode = row[10]
	v.CurrencyName = row[11]
	v.Phone = row[12]
	v.PostalCodeFormat = row[13]
	v.PostalCodeRegex = row[14]
	v.Languages = value.ParseMultipleValues[string](row[15])
	v.Neighbours = value.ParseMultipleValues[value.CountryCode](row[17])
	v.EquivalentFipsCode = row[18]

	return nil
}

type TimeZone struct {
	// CountryCode ISO-3166 2-letter country code
	CountryCode value.CountryCode
	// Name timezoneId
	Name string
	// GMTOffset the offset in hours from GMT on 1st of January
	GMTOffset float64
	// GMTOffset the offset to GMT on 1st of July (of the current year)
	DSTOffset float64
	// RawOffset without DST
	RawOffset float64
}

func (v *TimeZone) UnmarshalRow(row []string) error {
	const columns = 5

	var err error
	if err = checkColumns(row, columns); err != nil {
		return err
	}

	if v.GMTOffset, err = value.ParseFloat64(row[2]); err != nil {
		return fmt.Errorf("parse GMTOffset => %w", err)
	}

	if v.DSTOffset, err = value.ParseFloat64(row[3]); err != nil {
		return fmt.Errorf("parse DSTOffset => %w", err)
	}

	if v.RawOffset, err = value.ParseFloat64(row[4]); err != nil {
		return fmt.Errorf("parse RawOffset => %w", err)
	}

	v.CountryCode = value.CountryCode(row[0])
	v.Name = row[1]

	return nil
}

type Feature struct {
	Code        string
	Name        string
	Description string
}

func (v *Feature) UnmarshalRow(row []string) error {
	const columns = 3

	if err := checkColumns(row, columns); err != nil {
		return err
	}

	v.Code = row[0]
	v.Name = row[1]
	v.Description = row[2]

	return nil
}

type UserTag struct {
	ID    uint64
	Value string
}

func (v *UserTag) UnmarshalRow(row []string) error {
	const columns = 2

	var err error

	if err = checkColumns(row, columns); err != nil {
		return err
	}

	if v.ID, err = value.ParseUint64(row[0]); err != nil {
		return fmt.Errorf("parse ID => %w", err)
	}

	v.Value = row[1]

	return nil
}

type Language struct {
	// ISO639-1 2-letter code
	ISO6391 string
	// ISO639-2 3-letter code
	ISO6392 string
	// ISO639-3 3-letter code
	ISO6393 string
	// Language name
	Name string
}

func (v *Language) UnmarshalRow(row []string) error {
	const columns = 4

	var err error

	if err = checkColumns(row, columns); err != nil {
		return err
	}

	v.ISO6391 = row[2]
	v.ISO6392 = row[1]
	v.ISO6393 = row[0]
	v.Name = row[3]

	return nil
}

func checkColumns(row []string, expected int) error {
	if len(row) != expected {
		return fmt.Errorf("%w, expected %d, got %d", ErrInvalidRowLength, expected, len(row))
	}

	return nil
}
