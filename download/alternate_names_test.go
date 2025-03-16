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
)

func Test_Client_AlternateNames(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AlternateName, error) {
		res := make([]AlternateName, 0)

		err := client.AlternateNames(ctx, func(parsed AlternateName) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCases := []testSuite[AlternateName]{
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
								"alternateNamesV2.zip",
							)
						}),
					).Once().Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "alternateNamesV2Success.zip"),
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
		},
		{
			name: "context canceled",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "alternateNamesV2Success.zip"),
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
			exp: exp[AlternateName]{
				res: []AlternateName{},
				err: errors.New("parse file \"alternateNamesV2.txt\" in archive => parse file => context canceled"),
			},
		},
		{
			name: "missing target file in archive",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "alternateNamesV2Missing.zip"),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[AlternateName]{
				res: []AlternateName{},
				err: errors.New("parse file \"alternateNamesV2.txt\" in archive => file not found in archive"),
			},
		},
		{
			name: "invalid archive",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       testutil.MustOpen(testdata.FS, "alternateNamesV2Invalid.zip"),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[AlternateName]{
				res: []AlternateName{},
				err: errors.New("parse file \"alternateNamesV2.txt\" in archive => open zip archive => zip: not a valid zip file"),
			},
		},
		{
			name: "no content",
			args: args{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(
						&http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("")),
						},
						nil,
					)
				}),
				ctx: context.Background(),
			},
			exp: exp[AlternateName]{
				res: []AlternateName{},
				err: errors.New("parse file \"alternateNamesV2.txt\" in archive => open zip archive => zip: not a valid zip file"),
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
			exp: exp[AlternateName]{
				res: []AlternateName{},
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
			exp: exp[AlternateName]{
				res: []AlternateName{},
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
			exp: exp[AlternateName]{
				res: []AlternateName{},
				err: errors.New("download file => http client do => assert.AnError general error for testing"),
			},
		},
		{
			name: "missing context",
			args: args{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				ctx:        nil,
			},
			exp: exp[AlternateName]{
				res: []AlternateName{},
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
