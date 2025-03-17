package download

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
)

func Test_Client_Languages(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]Language, error) {
		res := make([]Language, 0)

		err := client.Languages(ctx, func(parsed Language) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCase := testSuite[Language]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"iso-languagecodes.txt",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "languages.txt"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[Language]{
			res: []Language{
				{
					ISO6391: "xx",
					ISO6392: "xxx",
					ISO6393: "xxy",
					Name:    "Foo",
				},
				{
					ISO6391: "yy",
					ISO6392: "yyy",
					ISO6393: "yyx",
					Name:    "Bar",
				},
			},
			err: nil,
		},
	}

	testCase.run(t, caller)
}
