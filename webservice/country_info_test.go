package webservice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
	"github.com/platx/geonames/webservice/testdata"
)

func Test_Client_CountryInfo(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req CountryInfoRequest) ([]CountryDetailed, error) {
		return client.CountryInfo
	}

	testCases := []testSuite[CountryInfoRequest, []CountryDetailed]{
		{
			name: "success with request values",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On(
						"Do",
						mock.MatchedBy(func(given *http.Request) bool {
							return assertRequest(
								t,
								given,
								"/countryInfoJSON",
								url.Values{
									"country":  []string{"DE", "UA"},
									"lang":     []string{"en"},
									"type":     []string{"json"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "country_detailed.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{
					Country:  []value.CountryCode{value.CountryCodeGermany, value.CountryCodeUkraine},
					Language: "en",
				},
			},
			exp: exp[[]CountryDetailed]{
				res: []CountryDetailed{
					{
						ID:            1,
						Code:          value.CountryCodeGermany,
						Name:          "Germany",
						ContinentCode: value.ContinentCodeEurope,
						ContinentName: "Europe",
						Capital:       "Berlin",
						Languages:     []string{"de"},
						BoundingBox: value.BoundingBox{
							East:  1.1,
							West:  1.2,
							North: -1.1,
							South: -1.2,
						},
						IsoAlpha3:        "DEU",
						IsoNumeric:       276,
						FipsCode:         "GM",
						Population:       111111,
						AreaInSqKm:       11111.1,
						PostalCodeFormat: "XXXXX",
						CurrencyCode:     "EUR",
					},
					{
						ID:            2,
						Code:          value.CountryCodeUkraine,
						Name:          "Ukraine",
						ContinentCode: value.ContinentCodeEurope,
						ContinentName: "Europe",
						Capital:       "Kyiv",
						Languages:     []string{"uk", "ru-UA"},
						BoundingBox: value.BoundingBox{
							East:  2.1,
							West:  2.2,
							North: -2.1,
							South: -2.2,
						},
						IsoAlpha3:        "UKR",
						IsoNumeric:       804,
						FipsCode:         "UP",
						Population:       222222,
						AreaInSqKm:       22222.2,
						PostalCodeFormat: "XXXXX",
						CurrencyCode:     "UAH",
					},
				},
				err: nil,
			},
		},
		{
			name: "empty without request values",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On(
						"Do",
						mock.MatchedBy(func(given *http.Request) bool {
							return assertRequest(
								t,
								given,
								"/countryInfoJSON",
								url.Values{
									"type":     []string{"json"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "geonames_empty.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: []CountryDetailed{},
				err: nil,
			},
		},
		{
			name: "invalid population",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"population": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: nil,
				err: errors.New("decode response => parse Population => strconv.ParseInt: parsing \"invalid\": invalid syntax"),
			},
		},
		{
			name: "invalid isoNumeric",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"isoNumeric": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: nil,
				err: errors.New("decode response => parse IsoNumeric => strconv.ParseUint: parsing \"invalid\": invalid syntax"),
			},
		},
		{
			name: "invalid areaInSqKm",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"areaInSqKm": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: nil,
				err: errors.New("decode response => parse AreaInSqKm => strconv.ParseFloat: parsing \"invalid\": invalid syntax"),
			},
		},
		{
			name: "invalid success response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geo`)),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: nil,
				err: errors.New("decode response => unexpected EOF"),
			},
		},
		{
			name: "error response",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       testutil.MustOpen(testdata.FS, "authorization_error.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: nil,
				err: errors.New("decode response => got error response => code: 10, message: \"user does not exist.\""),
			},
		},
		{
			name: "invalid error response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(strings.NewReader(`{"stat`)),
					})
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: nil,
				err: errors.New("decode response => unexpected EOF"),
			},
		},
		{
			name: "send request failed",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(nil, assert.AnError)
				}),
				userName: "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: context.Background(),
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: nil,
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[CountryInfoRequest]{
				ctx: nil,
				req: CountryInfoRequest{},
			},
			exp: exp[[]CountryDetailed]{
				res: nil,
				err: errors.New("create http request => net/http: nil Context"),
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
