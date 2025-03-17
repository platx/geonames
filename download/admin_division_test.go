package download

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
)

func Test_Client_AdminDivisionFirst(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AdminDivision, error) {
		res := make([]AdminDivision, 0)

		err := client.AdminDivisionFirst(ctx, func(parsed AdminDivision) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
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
			err: nil,
		},
	}

	testCase.run(t, caller)
}

func Test_Client_AdminDivisionSecond(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AdminDivision, error) {
		res := make([]AdminDivision, 0)

		err := client.AdminDivisionSecond(ctx, func(parsed AdminDivision) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
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
			err: nil,
		},
	}

	testCase.run(t, caller)
}

func Test_Client_AdminDivisionFifth(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AdminCode5, error) {
		res := make([]AdminCode5, 0)

		err := client.AdminDivisionFifth(ctx, func(parsed AdminCode5) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
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
			err: nil,
		},
	}

	testCase.run(t, caller)
}
