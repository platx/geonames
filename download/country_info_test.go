package download

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
)

func Test_Client_CountryInfo(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]Country, error) {
		res := make([]Country, 0)

		err := client.CountryInfo(ctx, func(parsed Country) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCases := []testSuite[Country]{
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
								"countryInfo.txt",
							)
						}),
					).Once().Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "countryInfo.txt"),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[Country]{
				res: []Country{
					{
						ID:                 1,
						Code:               value.CountryCodeUnitedStates,
						Name:               "United States",
						ContinentCode:      value.ContinentCodeNorthAmerica,
						Domain:             ".us",
						Capital:            "Washington",
						Languages:          []string{"en", "en-US", "es"},
						IsoAlpha3:          "USA",
						IsoNumeric:         111,
						FipsCode:           "UN",
						Population:         111111,
						AreaInSqKm:         11111,
						PostalCodeFormat:   "XXXXX",
						PostalCodeRegex:    "[a-zA-Z0-9]{5}",
						CurrencyCode:       "USD",
						CurrencyName:       "Dollar",
						Phone:              "+1",
						Neighbours:         []value.CountryCode{value.CountryCodeMexico, value.CountryCodeCanada},
						EquivalentFipsCode: "UNA",
					},
					{
						ID:                 2,
						Code:               value.CountryCodeUnitedKingdom,
						Name:               "United Kingdom",
						ContinentCode:      value.ContinentCodeEurope,
						Domain:             ".uk",
						Capital:            "London",
						Languages:          []string{"en-GB"},
						IsoAlpha3:          "GRB",
						IsoNumeric:         222,
						FipsCode:           "UK",
						Population:         222222,
						AreaInSqKm:         22222,
						PostalCodeFormat:   "YYYYY",
						PostalCodeRegex:    "[0-9]{5}",
						CurrencyCode:       "GBP",
						CurrencyName:       "Pound",
						Phone:              "+2",
						Neighbours:         []value.CountryCode{},
						EquivalentFipsCode: "GBN",
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
							Body:       testutil.MustOpen(testdata.FS, "countryInfo.txt"),
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
			exp: exp[Country]{
				res: []Country{},
				err: errors.New("parse file => context canceled"),
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
			exp: exp[Country]{
				res: []Country{},
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
			exp: exp[Country]{
				res: []Country{},
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
			exp: exp[Country]{
				res: []Country{},
				err: errors.New("download file => http client do => assert.AnError general error for testing"),
			},
		},
		{
			name: "missing context",
			args: args{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				ctx:        nil,
			},
			exp: exp[Country]{
				res: []Country{},
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
