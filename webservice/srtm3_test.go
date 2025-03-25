package webservice

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
	"github.com/platx/geonames/webservice/testdata"
)

func Test_Client_SRTM3(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, position value.Position) (int32, error) {
		return client.SRTM3
	}

	testCases := []testSuite[value.Position, int32]{
		{
			name: "success with request values",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On(
						"Do",
						mock.MatchedBy(func(given *http.Request) bool {
							return assertRequest(
								t,
								given,
								"/srtm3JSON",
								url.Values{
									"lat":      []string{"1.111"},
									"lng":      []string{"-1.111"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "srtm3.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[value.Position]{
				ctx: context.Background(),
				req: value.Position{
					Latitude:  1.111,
					Longitude: -1.111,
				},
			},
			exp: exp[int32]{
				res: 111,
				err: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.run(t, caller)
		})
	}
}
