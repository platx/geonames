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

func Test_Client_TimeZones(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]TimeZone, error) {
		res := make([]TimeZone, 0)

		err := client.TimeZones(ctx, func(parsed TimeZone) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCase := testSuite[TimeZone]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"timeZones.txt",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "timeZones.txt"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[TimeZone]{
			res: []TimeZone{
				{
					CountryCode: value.CountryCodeUnitedStates,
					Name:        "America/Unites_States",
					GMTOffset:   -1.1,
					DSTOffset:   0.1,
					RawOffset:   1.1,
				},
				{
					CountryCode: value.CountryCodeUnitedKingdom,
					Name:        "Europe/United_Kingdom",
					GMTOffset:   -2.2,
					DSTOffset:   0.2,
					RawOffset:   2.2,
				},
			},
			err: nil,
		},
	}

	testCase.run(t, caller)
}
