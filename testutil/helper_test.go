package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/platx/geonames/download/testdata"
)

func Test_MustOpen(t *testing.T) {
	t.Parallel()

	t.Run("exist", func(t *testing.T) {
		t.Parallel()

		got := MustOpen(testdata.FS, "countryInfo.txt")

		assert.NotNil(t, got)
	})

	t.Run("not exist", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			got := MustOpen(testdata.FS, "invalid.txt")

			assert.Nil(t, got)
		})
	})
}
