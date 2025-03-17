package download

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
)

func Test_Client_AdminSubdivision1(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AdminSubdivision, error) {
		res := make([]AdminSubdivision, 0)

		err := client.AdminSubdivision1(ctx, func(parsed AdminSubdivision) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCase := testSuite[AdminSubdivision]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"admin1CodesASCII.txt",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "adminSubdivision1.txt"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[AdminSubdivision]{
			res: []AdminSubdivision{
				{
					ID:        1,
					Code:      "XX.XY",
					Name:      "Foo1",
					NameASCII: "Foo2",
				},
				{
					ID:        2,
					Code:      "XY.YX",
					Name:      "Bar1",
					NameASCII: "Bar2",
				},
			},
			err: nil,
		},
	}

	testCase.run(t, caller)
}
