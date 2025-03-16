package download

import (
	"context"
	"net/http"
	"testing"

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

	testCase := testSuite[Country]{
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
	}

	testCase.run(t, caller)
}
