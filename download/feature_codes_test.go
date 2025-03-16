package download

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
)

func Test_Client_FeatureCodes(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]Feature, error) {
		res := make([]Feature, 0)

		err := client.FeatureCodes(ctx, "en", func(parsed Feature) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCase := testSuite[Feature]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"featureCodes_en.txt",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "featureCodes.txt"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[Feature]{
			res: []Feature{
				{
					Code:        "AAAA",
					Name:        "First name",
					Description: "First description",
				},
				{
					Code:        "BBBB",
					Name:        "Second name",
					Description: "Second description",
				},
			},
			err: nil,
		},
	}

	testCase.run(t, caller)
}
