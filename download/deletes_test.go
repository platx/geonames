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

func Test_Client_Deletes(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]GeoNameDeleted, error) {
		res := make([]GeoNameDeleted, 0)

		err := client.Deletes(ctx, func(parsed GeoNameDeleted) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCase := testSuite[GeoNameDeleted]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							fmt.Sprintf("deletes-%s.txt", yesterday().Format(time.DateOnly)),
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "deletes.txt"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[GeoNameDeleted]{
			res: []GeoNameDeleted{
				{
					ID:      1,
					Name:    "Name 1",
					Comment: "Comment 1",
				},
				{
					ID:      2,
					Name:    "Name 2",
					Comment: "Comment 2",
				},
			},
			err: nil,
		},
	}

	testCase.run(t, caller)
}
