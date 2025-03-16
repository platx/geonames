package download

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
)

func Test_Client_AllCountries(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]GeoName, error) {
		res := make([]GeoName, 0)

		err := client.AllCountries(ctx, func(parsed GeoName) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCases := []testSuite[GeoName]{
		{
			name: "success",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On(
						"Do",
						mock.MatchedBy(func(given *http.Request) bool {
							return assertRequest(
								t,
								given,
								"allCountries.zip",
							)
						}),
					).Once().Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "allCountries.zip"),
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
						ASCIIName:      "New York City",
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
						ASCIIName:      "London",
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
				err: nil,
			},
		},
		{
			name: "context canceled",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "allCountries.zip"),
						},
						nil,
					)
				}),
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()

					return ctx
				}(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("parse file \"allCountries.txt\" in archive => parse file => context canceled"),
			},
		},
		{
			name: "missing target file in archive",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "missing.zip"),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("parse file \"allCountries.txt\" in archive => file not found in archive"),
			},
		},
		{
			name: "invalid archive",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "invalid.zip"),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("parse file \"allCountries.txt\" in archive => open zip archive => zip: not a valid zip file"),
			},
		},
		{
			name: "no content",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("")),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("parse file \"allCountries.txt\" in archive => open zip archive => zip: not a valid zip file"),
			},
		},
		{
			name: "content copy failed",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body: testutil.MockReadCloser(func(m *testutil.ReadCloserMock) {
								m.On("Read", mock.Anything).Return(0, assert.AnError)
								m.On("Close").Return(nil)
							}),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => copy file content => assert.AnError general error for testing"),
			},
		},
		{
			name: "invalid status code",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(strings.NewReader("")),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => unexpected status code: 500"),
			},
		},
		{
			name: "http send request failed",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(nil, assert.AnError)
				}),
				ctx: context.Background(),
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => http client do => assert.AnError general error for testing"),
			},
		},
		{
			name: "missing context",
			args: args{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				ctx:        nil,
			},
			exp: exp[GeoName]{
				res: []GeoName{},
				err: errors.New("download file => create http request => net/http: nil Context"),
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.run(t, caller)
		})
	}
}
