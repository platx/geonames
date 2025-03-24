package download

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
)

func Test_Client_Cities500(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]GeoName, []error) {
		return collect(client.Cities500(ctx))
	}

	testCase := testSuite[GeoName]{
		name: "success",
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"cities500.zip",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "cities500.zip"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[GeoName]{
			res: []GeoName{
				{
					ID:             1,
					Name:           "New York City",
					NameASCII:      "New York City",
					AlternateNames: []string{"NYC", "NewYork", "Нью-Йорк"},
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					FeatureClass: "A",
					FeatureCode:  "AAAA",
					CountryCode:  value.CountryCodeUnitedStates,
					AlternateCountryCodes: []value.CountryCode{
						value.CountryCodeUnitedKingdom,
						value.CountryCodeUkraine,
					},
					AdminCode: value.AdminCode{
						First:  "FOO",
						Second: "BAR",
						Third:  "BAZ",
						Fourth: "FOOBAR",
					},
					Population:            111111,
					Elevation:             111,
					DigitalElevationModel: 11,
					Timezone:              "America/New_York",
					ModificationDate:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:             2,
					Name:           "London",
					NameASCII:      "London",
					AlternateNames: []string{"Landan", "Лондон"},
					Position: value.Position{
						Latitude:  2.222,
						Longitude: -2.222,
					},
					FeatureClass:          "B",
					FeatureCode:           "BBBB",
					CountryCode:           value.CountryCodeUnitedKingdom,
					AlternateCountryCodes: []value.CountryCode{},
					AdminCode: value.AdminCode{
						First:  "FOO",
						Second: "BAR",
						Third:  "BAZ",
						Fourth: "FOOBAR",
					},
					Population:            222222,
					Elevation:             222,
					DigitalElevationModel: 22,
					Timezone:              "Europe/London",
					ModificationDate:      time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			err: []error{
				errors.New("parse ID => strconv.ParseUint: parsing \"v\": invalid syntax"),
				errors.New("parse Position => latitude => strconv.ParseFloat: parsing \"v\": invalid syntax"),
				errors.New("parse Position => longitude => strconv.ParseFloat: parsing \"v\": invalid syntax"),
				errors.New("parse Population => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse Elevation => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse DigitalElevationModel => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse ModificationDate => parsing time \"v\" as \"2006-01-02\": cannot parse \"v\" as \"2006\""),
				errors.New("invalid row length, expected 19, got 3"),
			},
		},
	}

	testCase.run(t, caller)
}

func Test_Client_Cities1000(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]GeoName, []error) {
		return collect(client.Cities1000(ctx))
	}

	testCase := testSuite[GeoName]{
		name: "success",
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"cities1000.zip",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "cities1000.zip"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[GeoName]{
			res: []GeoName{
				{
					ID:             1,
					Name:           "New York City",
					NameASCII:      "New York City",
					AlternateNames: []string{"NYC", "NewYork", "Нью-Йорк"},
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					FeatureClass: "A",
					FeatureCode:  "AAAA",
					CountryCode:  value.CountryCodeUnitedStates,
					AlternateCountryCodes: []value.CountryCode{
						value.CountryCodeUnitedKingdom,
						value.CountryCodeUkraine,
					},
					AdminCode: value.AdminCode{
						First:  "FOO",
						Second: "BAR",
						Third:  "BAZ",
						Fourth: "FOOBAR",
					},
					Population:            111111,
					Elevation:             111,
					DigitalElevationModel: 11,
					Timezone:              "America/New_York",
					ModificationDate:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:             2,
					Name:           "London",
					NameASCII:      "London",
					AlternateNames: []string{"Landan", "Лондон"},
					Position: value.Position{
						Latitude:  2.222,
						Longitude: -2.222,
					},
					FeatureClass:          "B",
					FeatureCode:           "BBBB",
					CountryCode:           value.CountryCodeUnitedKingdom,
					AlternateCountryCodes: []value.CountryCode{},
					AdminCode: value.AdminCode{
						First:  "FOO",
						Second: "BAR",
						Third:  "BAZ",
						Fourth: "FOOBAR",
					},
					Population:            222222,
					Elevation:             222,
					DigitalElevationModel: 22,
					Timezone:              "Europe/London",
					ModificationDate:      time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			err: []error{
				errors.New("parse ID => strconv.ParseUint: parsing \"v\": invalid syntax"),
				errors.New("parse Position => latitude => strconv.ParseFloat: parsing \"v\": invalid syntax"),
				errors.New("parse Position => longitude => strconv.ParseFloat: parsing \"v\": invalid syntax"),
				errors.New("parse Population => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse Elevation => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse DigitalElevationModel => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse ModificationDate => parsing time \"v\" as \"2006-01-02\": cannot parse \"v\" as \"2006\""),
				errors.New("invalid row length, expected 19, got 3"),
			},
		},
	}

	testCase.run(t, caller)
}

func Test_Client_Cities5000(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]GeoName, []error) {
		return collect(client.Cities5000(ctx))
	}

	testCase := testSuite[GeoName]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"cities5000.zip",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "cities5000.zip"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[GeoName]{
			res: []GeoName{
				{
					ID:             1,
					Name:           "New York City",
					NameASCII:      "New York City",
					AlternateNames: []string{"NYC", "NewYork", "Нью-Йорк"},
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					FeatureClass: "A",
					FeatureCode:  "AAAA",
					CountryCode:  value.CountryCodeUnitedStates,
					AlternateCountryCodes: []value.CountryCode{
						value.CountryCodeUnitedKingdom,
						value.CountryCodeUkraine,
					},
					AdminCode: value.AdminCode{
						First:  "FOO",
						Second: "BAR",
						Third:  "BAZ",
						Fourth: "FOOBAR",
					},
					Population:            111111,
					Elevation:             111,
					DigitalElevationModel: 11,
					Timezone:              "America/New_York",
					ModificationDate:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:             2,
					Name:           "London",
					NameASCII:      "London",
					AlternateNames: []string{"Landan", "Лондон"},
					Position: value.Position{
						Latitude:  2.222,
						Longitude: -2.222,
					},
					FeatureClass:          "B",
					FeatureCode:           "BBBB",
					CountryCode:           value.CountryCodeUnitedKingdom,
					AlternateCountryCodes: []value.CountryCode{},
					AdminCode: value.AdminCode{
						First:  "FOO",
						Second: "BAR",
						Third:  "BAZ",
						Fourth: "FOOBAR",
					},
					Population:            222222,
					Elevation:             222,
					DigitalElevationModel: 22,
					Timezone:              "Europe/London",
					ModificationDate:      time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			err: []error{
				errors.New("parse ID => strconv.ParseUint: parsing \"v\": invalid syntax"),
				errors.New("parse Position => latitude => strconv.ParseFloat: parsing \"v\": invalid syntax"),
				errors.New("parse Position => longitude => strconv.ParseFloat: parsing \"v\": invalid syntax"),
				errors.New("parse Population => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse Elevation => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse DigitalElevationModel => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse ModificationDate => parsing time \"v\" as \"2006-01-02\": cannot parse \"v\" as \"2006\""),
				errors.New("invalid row length, expected 19, got 3"),
			},
		},
	}

	testCase.run(t, caller)
}

func Test_Client_Cities15000(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]GeoName, []error) {
		return collect(client.Cities15000(ctx))
	}

	testCase := testSuite[GeoName]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"cities15000.zip",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "cities15000.zip"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[GeoName]{
			res: []GeoName{
				{
					ID:             1,
					Name:           "New York City",
					NameASCII:      "New York City",
					AlternateNames: []string{"NYC", "NewYork", "Нью-Йорк"},
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					FeatureClass: "A",
					FeatureCode:  "AAAA",
					CountryCode:  value.CountryCodeUnitedStates,
					AlternateCountryCodes: []value.CountryCode{
						value.CountryCodeUnitedKingdom,
						value.CountryCodeUkraine,
					},
					AdminCode: value.AdminCode{
						First:  "FOO",
						Second: "BAR",
						Third:  "BAZ",
						Fourth: "FOOBAR",
					},
					Population:            111111,
					Elevation:             111,
					DigitalElevationModel: 11,
					Timezone:              "America/New_York",
					ModificationDate:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:             2,
					Name:           "London",
					NameASCII:      "London",
					AlternateNames: []string{"Landan", "Лондон"},
					Position: value.Position{
						Latitude:  2.222,
						Longitude: -2.222,
					},
					FeatureClass:          "B",
					FeatureCode:           "BBBB",
					CountryCode:           value.CountryCodeUnitedKingdom,
					AlternateCountryCodes: []value.CountryCode{},
					AdminCode: value.AdminCode{
						First:  "FOO",
						Second: "BAR",
						Third:  "BAZ",
						Fourth: "FOOBAR",
					},
					Population:            222222,
					Elevation:             222,
					DigitalElevationModel: 22,
					Timezone:              "Europe/London",
					ModificationDate:      time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			err: []error{
				errors.New("parse ID => strconv.ParseUint: parsing \"v\": invalid syntax"),
				errors.New("parse Position => latitude => strconv.ParseFloat: parsing \"v\": invalid syntax"),
				errors.New("parse Position => longitude => strconv.ParseFloat: parsing \"v\": invalid syntax"),
				errors.New("parse Population => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse Elevation => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse DigitalElevationModel => strconv.ParseInt: parsing \"v\": invalid syntax"),
				errors.New("parse ModificationDate => parsing time \"v\" as \"2006-01-02\": cannot parse \"v\" as \"2006\""),
				errors.New("invalid row length, expected 19, got 3"),
			},
		},
	}

	testCase.run(t, caller)
}
