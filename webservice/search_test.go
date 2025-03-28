package webservice

import (
	"context"
	"net/http"
	"net/url"
	"testing"

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
								"/searchJSON",
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
						ID: 1,
						Country: value.Country{
							ID:   11,
							Code: value.CountryCodeUnitedStates,
							Name: "United States",
						},
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
						Feature: value.Feature{
							Class:     "A",
							ClassName: "Test class",
							Code:      "AAAA",
							CodeName:  "Test code",
						},
						Name:        "New York City",
						ToponymName: "New York City",
						Position: value.Position{
							Latitude:  1.111,
							Longitude: -1.111,
						},
						Population: 111111,
					},
					{
						ID: 2,
						Country: value.Country{
							ID:   22,
							Code: value.CountryCodeUnitedKingdom,
							Name: "United Kingdom",
						},
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
						Feature: value.Feature{
							Class:     "A",
							ClassName: "Test class",
							Code:      "AAAA",
							CodeName:  "Test code",
						},
						Name:        "London",
						ToponymName: "London",
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
								"/searchJSON",
								url.Values{
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.run(t, caller)
		})
	}
}
