package download

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
)

func Test_Client_UserTags(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]UserTag, []error) {
		return collect(client.UserTags(ctx))
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
			err: []error{
				errors.New("parse ID => strconv.ParseUint: parsing \"v\": invalid syntax"),
				errors.New("invalid row length, expected 2, got 1"),
			},
		},
	}

	testCase.run(t, caller)
}
