package download

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
)

func Test_Client_AlternateNamesModifications(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AlternateName, error) {
		res := make([]AlternateName, 0)

		err := client.AlternateNamesModifications(ctx, func(parsed AlternateName) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCase := testSuite[AlternateName]{
		name: "success",
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							fmt.Sprintf("alternateNamesModifications-%s.txt", yesterday().Format(time.DateOnly)),
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "alternateNames.txt"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[AlternateName]{
			res: []AlternateName{
				{
					ID:         1,
					GeoNameID:  11,
					Language:   "en-US",
					Value:      "New York City",
					Preferred:  true,
					Short:      true,
					Colloquial: true,
					Historic:   true,
					From:       "1901",
					To:         "2000",
				},
				{
					ID:         2,
					GeoNameID:  22,
					Language:   "en-GB",
					Value:      "London",
					Preferred:  false,
					Short:      false,
					Colloquial: false,
					Historic:   false,
					From:       "",
					To:         "",
				},
			},
			err: nil,
		},
	}

	testCase.run(t, caller)
}
