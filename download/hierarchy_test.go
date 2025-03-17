package download

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
)

func Test_Client_Hierarchy(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]HierarchyItem, error) {
		res := make([]HierarchyItem, 0)

		err := client.Hierarchy(ctx, func(parsed HierarchyItem) error {
			res = append(res, parsed)

			return nil
		})

		return res, err
	}

	testCase := testSuite[HierarchyItem]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							"hierarchy.zip",
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "hierarchy.zip"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[HierarchyItem]{
			res: []HierarchyItem{
				{
					ParentID: 1,
					ChildID:  2,
					Type:     "XX",
				},
				{
					ParentID: 2,
					ChildID:  3,
					Type:     "YY",
				},
				{
					ParentID: 3,
					ChildID:  4,
					Type:     "",
				},
			},
			err: nil,
		},
	}

	testCase.run(t, caller)
}
