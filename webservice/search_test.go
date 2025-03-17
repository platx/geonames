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

func Test_Client_GeoNameSearch(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req SearchRequest) ([]GeoName, error) {
		return client.Search
	}

	testCases := []testSuite[SearchRequest, []GeoName]{
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
								"/search",
								url.Values{
									"q":               []string{"London"},
									"name":            []string{"London"},
									"name_equals":     []string{"London"},
									"name_startsWith": []string{"Lon"},
									"maxRows":         []string{"2"},
									"startRow":        []string{"1"},
									"country":         []string{"GB", "US"},
									"countryBias":     []string{"GB"},
									"continentCode":   []string{"EU"},
									"adminCode1":      []string{"FOO1"},
									"adminCode2":      []string{"FOO2"},
									"adminCode3":      []string{"FOO3"},
									"adminCode4":      []string{"FOO4"},
									"adminCode5":      []string{"FOO5"},
									"featureClass":    []string{"A"},
									"featureCode":     []string{"AAAA", "AAAB"},
									"cities":          []string{"cities5000"},
									"lang":            []string{"en"},
									"searchLanguage":  []string{"en-GB"},
									"isNameRequired":  []string{"true"},
									"tag":             []string{"tag1"},
									"operator":        []string{"AND"},
									"fuzzy":           []string{"0.5"},
									"west":            []string{"1"},
									"east":            []string{"2"},
									"north":           []string{"-1"},
									"south":           []string{"-2"},
									"orderby":         []string{"population"},
									"type":            []string{"json"},
									"username":        []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "geonames.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[SearchRequest]{
				ctx: context.Background(),
				req: SearchRequest{
					Query:          "London",
					Name:           "London",
					NameEquals:     "London",
					NameStartsWith: "Lon",
					MaxRows:        2,
					StartRow:       1,
					Country:        []value.CountryCode{value.CountryCodeUnitedKingdom, value.CountryCodeUnitedStates},
					CountryBias:    value.CountryCodeUnitedKingdom,
					ContinentCode:  value.ContinentCodeEurope,
					AdminCode: value.AdminCode{
						First:  "FOO1",
						Second: "FOO2",
						Third:  "FOO3",
						Fourth: "FOO4",
						Fifth:  "FOO5",
					},
					FeatureClass:   []string{"A"},
					FeatureCode:    []string{"AAAA", "AAAB"},
					Cities:         value.Cities5000,
					Language:       "en",
					SearchLanguage: "en-GB",
					NameRequired:   true,
					Tag:            "tag1",
					Operator:       value.OperatorAnd,
					Fuzzy:          0.5,
					BoundingBox: value.BoundingBox{
						West:  1.0,
						East:  2.0,
						North: -1.0,
						South: -2.0,
					},
					OrderBy: value.OrderByPopulation,
				},
			},
			exp: exp[[]GeoName]{
				res: []GeoName{
					{
						ID:          1,
						CountryID:   11,
						CountryCode: value.CountryCodeUnitedStates,
						CountryName: "United States",
						AdminSubdivision: value.AdminDivisions{
							First: value.AdminDivision{
								Code: "FOO",
								Name: "Foo",
							},
							Second: value.AdminDivision{
								Code: "BAR",
								Name: "Bar",
							},
							Third: value.AdminDivision{
								Code: "BAZ",
								Name: "Baz",
							},
							Fourth: value.AdminDivision{
								Code: "FOOBAR",
								Name: "FooBar",
							},
							Fifth: value.AdminDivision{},
						},
						FeatureClass:     "A",
						FeatureClassName: "Test class",
						FeatureCode:      "AAAA",
						FeatureCodeName:  "Test code",
						Name:             "New York City",
						ToponymName:      "New York City",
						Position: value.Position{
							Latitude:  1.111,
							Longitude: -1.111,
						},
						Population: 111111,
					},
					{
						ID:          2,
						CountryID:   22,
						CountryCode: value.CountryCodeUnitedKingdom,
						CountryName: "United Kingdom",
						AdminSubdivision: value.AdminDivisions{
							First: value.AdminDivision{
								Code: "FOO",
								Name: "Foo",
							},
							Second: value.AdminDivision{
								Code: "BAR",
								Name: "Bar",
							},
							Third: value.AdminDivision{
								Code: "BAZ",
								Name: "Baz",
							},
							Fourth: value.AdminDivision{
								Code: "FOOBAR",
								Name: "FooBar",
							},
							Fifth: value.AdminDivision{},
						},
						FeatureClass:     "A",
						FeatureClassName: "Test class",
						FeatureCode:      "AAAA",
						FeatureCodeName:  "Test code",
						Name:             "London",
						ToponymName:      "London",
						Position: value.Position{
							Latitude:  2.222,
							Longitude: -2.222,
						},
						Population: 222222,
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
								"/search",
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
			args: args[SearchRequest]{
				ctx: context.Background(),
				req: SearchRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName{},
				err: nil,
			},
		},
		{
			name: "invalid country id",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"countryId": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[SearchRequest]{
				ctx: context.Background(),
				req: SearchRequest{},
			},
			exp: exp[[]GeoName]{
				res: nil,
				err: errors.New("decode response => parse CountryID => strconv.ParseUint: parsing \"invalid\": invalid syntax"),
			},
		},
		{
			name: "invalid longitude",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"lng": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[SearchRequest]{
				ctx: context.Background(),
				req: SearchRequest{},
			},
			exp: exp[[]GeoName]{
				res: nil,
				err: errors.New("decode response => parse Position => longitude => strconv.ParseFloat: parsing \"invalid\": invalid syntax"),
			},
		},
		{
			name: "invalid latitude",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geonames": [{"lat": "invalid"}]}`)),
					})
				}),
				userName: "test-user",
			},
			args: args[SearchRequest]{
				ctx: context.Background(),
				req: SearchRequest{},
			},
			exp: exp[[]GeoName]{
				res: nil,
				err: errors.New("decode response => parse Position => latitude => strconv.ParseFloat: parsing \"invalid\": invalid syntax"),
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
			args: args[SearchRequest]{
				ctx: context.Background(),
				req: SearchRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
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
			args: args[SearchRequest]{
				ctx: context.Background(),
				req: SearchRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
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
			args: args[SearchRequest]{
				ctx: context.Background(),
				req: SearchRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
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
			args: args[SearchRequest]{
				ctx: context.Background(),
				req: SearchRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[SearchRequest]{
				ctx: nil,
				req: SearchRequest{},
			},
			exp: exp[[]GeoName]{
				res: []GeoName(nil),
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
