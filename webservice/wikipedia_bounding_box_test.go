package webservice

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/testutil"
	"github.com/platx/geonames/value"
	"github.com/platx/geonames/webservice/testdata"
)

func Test_Client_WikipediaBoundingBox(t *testing.T) {
	t.Parallel()

	caller := func(client *Client) func(ctx context.Context, req WikipediaBoundingBoxRequest) ([]Wikipedia, error) {
		return client.WikipediaBoundingBox
	}

	testCases := []testSuite[WikipediaBoundingBoxRequest, []Wikipedia]{
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
								"/wikipediaBoundingBoxJSON",
								url.Values{
									"west":     []string{"1"},
									"east":     []string{"2"},
									"north":    []string{"-1"},
									"south":    []string{"-2"},
									"lang":     []string{"en"},
									"maxRows":  []string{"2"},
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "wikipedia.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[WikipediaBoundingBoxRequest]{
				ctx: context.Background(),
				req: WikipediaBoundingBoxRequest{
					BoundingBox: value.BoundingBox{
						West:  1.0,
						East:  2.0,
						North: -1.0,
						South: -2.0,
					},
					Language: "en",
					MaxRows:  2,
				},
			},
			exp: exp[[]Wikipedia]{
				res: []Wikipedia{
					{
						ID:          1,
						CountryCode: value.CountryCodeUnitedStates,
						Position: value.Position{
							Latitude:  1.111,
							Longitude: -1.111,
						},
						Feature:      "foo",
						Elevation:    111,
						Rank:         100,
						Language:     "en",
						WikipediaURL: "example.com/foo",
						ThumbnailURL: "https://example.com/foo.jpg",
						Title:        "Foo",
						Summary:      "My to considered delightful invitation announcing of no decisively boisterous. Did add dashwoods deficient man concluded additions resources.",
					},
					{
						ID:          2,
						CountryCode: value.CountryCodeUnitedKingdom,
						Position: value.Position{
							Latitude:  2.222,
							Longitude: -2.222,
						},
						Feature:      "bar",
						Elevation:    222,
						Rank:         200,
						Language:     "es",
						WikipediaURL: "example.com/bar",
						ThumbnailURL: "https://example.com/bar.jpg",
						Title:        "Bar",
						Summary:      "Full he none no side. Uncommonly surrounded considered for him are its. It we is read good soon.",
					},
				},
				err: nil,
			},
		},
		{
			name: "empty without request values",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On(
						"Do",
						mock.MatchedBy(func(given *http.Request) bool {
							return assertRequest(
								t,
								given,
								"/wikipediaBoundingBoxJSON",
								url.Values{
									"username": []string{"test-user"},
								},
							)
						}),
					).Once().Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "geonames_empty.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[WikipediaBoundingBoxRequest]{
				ctx: context.Background(),
				req: WikipediaBoundingBoxRequest{},
			},
			exp: exp[[]Wikipedia]{
				res: []Wikipedia{},
				err: nil,
			},
		},
		{
			name: "invalid success response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"geo`)),
					})
				}),
				userName: "test-user",
			},
			args: args[WikipediaBoundingBoxRequest]{
				ctx: context.Background(),
				req: WikipediaBoundingBoxRequest{},
			},
			exp: exp[[]Wikipedia]{
				res: nil,
				err: errors.New("decode response => unexpected EOF"),
			},
		},
		{
			name: "error response",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       testutil.MustOpen(testdata.FS, "authorization_error.json"),
					})
				}),
				userName: "test-user",
			},
			args: args[WikipediaBoundingBoxRequest]{
				ctx: context.Background(),
				req: WikipediaBoundingBoxRequest{},
			},
			exp: exp[[]Wikipedia]{
				res: nil,
				err: errors.New("decode response => got error response => code: 10, message: \"user does not exist.\""),
			},
		},
		{
			name: "invalid error response body",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(strings.NewReader(`{"stat`)),
					})
				}),
				userName: "test-user",
			},
			args: args[WikipediaBoundingBoxRequest]{
				ctx: context.Background(),
				req: WikipediaBoundingBoxRequest{},
			},
			exp: exp[[]Wikipedia]{
				res: nil,
				err: errors.New("decode response => unexpected EOF"),
			},
		},
		{
			name: "send request failed",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
					m.On("Do", mock.Anything).Return(nil, assert.AnError)
				}),
				userName: "test-user",
			},
			args: args[WikipediaBoundingBoxRequest]{
				ctx: context.Background(),
				req: WikipediaBoundingBoxRequest{},
			},
			exp: exp[[]Wikipedia]{
				res: nil,
				err: fmt.Errorf("send http request => %w", assert.AnError),
			},
		},
		{
			name: "context not provided",
			deps: deps{
				httpClient: testutil.MockHTTPClient(func(_ *testutil.HTTPClientMock) {}),
				userName:   "test-user",
			},
			args: args[WikipediaBoundingBoxRequest]{
				ctx: nil,
				req: WikipediaBoundingBoxRequest{},
			},
			exp: exp[[]Wikipedia]{
				res: nil,
				err: errors.New("create http request => net/http: nil Context"),
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
