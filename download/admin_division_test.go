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

func Test_Client_AdminDivisionFirst(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AdminDivision, []error) {
		return collect(client.AdminDivisionFirst(ctx))
	}

	testCase := testSuite[AdminDivision]{
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
						Body:       testutil.MustOpen(testdata.FS, "adminDivision.txt"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[AdminDivision]{
			res: []AdminDivision{
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
			err: []error{
				errors.New("parse ID => strconv.ParseUint: parsing \"v\": invalid syntax"),
				errors.New("invalid row length, expected 4, got 3"),
			},
		},
	}

	testCase.run(t, caller)
}

func Test_Client_AdminDivisionSecond(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AdminDivision, []error) {
		return collect(client.AdminDivisionSecond(ctx))
	}

	testCase := testSuite[AdminDivision]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"admin2Codes.txt",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "adminDivision.txt"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[AdminDivision]{
			res: []AdminDivision{
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
			err: []error{
				errors.New("parse ID => strconv.ParseUint: parsing \"v\": invalid syntax"),
				errors.New("invalid row length, expected 4, got 3"),
			},
		},
	}

	testCase.run(t, caller)
}

func Test_Client_AdminDivisionFifth(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AdminCode5, []error) {
		return collect(client.AdminDivisionFifth(ctx))
	}

	testCase := testSuite[AdminCode5]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"adminCode5.zip",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "adminCode5.zip"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[AdminCode5]{
			res: []AdminCode5{
				{
					ID:   1,
					Code: "XX",
				},
				{
					ID:   2,
					Code: "YY",
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
