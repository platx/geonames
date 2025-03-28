package download

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/platx/geonames/download/testdata"
	"github.com/platx/geonames/testutil"
)

func Test_Client_AlternateNamesDeletes(t *testing.T) {
	t.Parallel()

	caller := func(client *Client, ctx context.Context) ([]AlternateNameDeleted, []error) {
		return collect(client.AlternateNamesDeletes(ctx))
	}

	testCase := testSuite[AlternateNameDeleted]{
		args: args{
			httpClient: testutil.MockHTTPClient(func(m *testutil.HTTPClientMock) {
				m.On(
					"Do",
					mock.MatchedBy(func(given *http.Request) bool {
						return assertRequest(
							t,
							given,
							fmt.Sprintf("alternateNamesDeletes-%s.txt", yesterday().Format(time.DateOnly)),
						)
					}),
				).Once().Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       testutil.MustOpen(testdata.FS, "alternateNamesDeletes.txt"),
					},
					nil,
				)
			}),
			ctx: context.Background(),
		},
		exp: exp[AlternateNameDeleted]{
			res: []AlternateNameDeleted{
				{
					AlternateNameID: 1,
					GeoNameID:       11,
					Name:            "Name 1",
					Comment:         "Comment 1",
				},
				{
					AlternateNameID: 2,
					GeoNameID:       22,
					Name:            "Name 2",
					Comment:         "Comment 2",
				},
			},
			err: []error{
				errors.New("parse AlternateNameID => strconv.ParseUint: parsing \"v\": invalid syntax"),
				errors.New("parse GeoNameID => strconv.ParseUint: parsing \"vv\": invalid syntax"),
				errors.New("invalid row length, expected 4, got 3"),
			},
		},
	}

	testCase.run(t, caller)
}
