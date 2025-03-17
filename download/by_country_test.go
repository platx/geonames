package download

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
)

func Test_Client_ByCountry(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]GeoName, error) {
		res := make([]GeoName, 0)

		err := client.ByCountry(ctx, value.CountryCodeUnitedStates, func(parsed GeoName) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
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
							"US.zip",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "byCountry.zip"),
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
			err: nil,
		},
	}

	testCase.run(t, caller)
}
