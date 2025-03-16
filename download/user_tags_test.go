package download

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
)

func Test_Client_UserTags(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]UserTag, error) {
		res := make([]UserTag, 0)

		err := client.UserTags(ctx, func(parsed UserTag) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCase := testSuite[UserTag]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"userTags.zip",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "userTags.zip"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[UserTag]{
			res: []UserTag{
				{
					ID:    1,
					Value: "Foo",
				},
				{
					ID:    2,
					Value: "Bar",
				},
			},
			err: nil,
		},
	}

	testCase.run(t, caller)
}
