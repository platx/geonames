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

func Test_Client_Ocean(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req OceanRequest) (Ocean, error) {
		return client.Ocean
	}

	testCases := []testSuite[OceanRequest, Ocean]{
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
								"/oceanJSON",
								url.Values{
									"lat":      []string{"1.111"},
									"lng":      []string{"-1.111"},
									"radius":   []string{"10"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "ocean.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[OceanRequest]{
				ctx: context.Background(),
				req: OceanRequest{
					Position: value.Position{
						Latitude:  1.111,
						Longitude: -1.111,
					},
					Radius: 10,
				},
			},
			exp: exp[Ocean]{
				res: Ocean{
					ID:       1,
					Distance: 0.111,
					Name:     "Test ocean",
				},
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
