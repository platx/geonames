package download

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

	testCases := []testSuite[TimeZone]{
		{
			name: "success",
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
							Body:       testutil.MustOpen(testdata.FS, "timeZonesSuccess.txt"),
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
		},
		{
			name: "context canceled",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "timeZonesSuccess.txt"),
						},
						nil,
					)
				}),
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()

					return ctx
				}(),
			},
			exp: exp[TimeZone]{
				res: []TimeZone{},
				err: errors.New("parse file => context canceled"),
			},
		},
		{
			name: "content copy failed",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body: testutil.MockReadCloser(func(m *testutil.ReadCloserMock) {
								m.On("Read", mock.Anything).Return(0, assert.AnError)
								m.On("Close").Return(nil)
							}),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[TimeZone]{
				res: []TimeZone{},
				err: errors.New("download file => copy file content => assert.AnError general error for testing"),
			},
		},
		{
			name: "invalid status code",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(strings.NewReader("")),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[TimeZone]{
				res: []TimeZone{},
				err: errors.New("download file => unexpected status code: 500"),
			},
		},
		{
			name: "http send request failed",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(nil, assert.AnError)
				}),
				ctx: context.Background(),
			},
			exp: exp[TimeZone]{
				res: []TimeZone{},
				err: errors.New("download file => http client do => assert.AnError general error for testing"),
			},
		},
		{
			name: "missing context",
			args: args{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				ctx:        nil,
			},
			exp: exp[TimeZone]{
				res: []TimeZone{},
				err: errors.New("download file => create http request => net/http: nil Context"),
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testCase.run(t, caller)
		})
	}
}
