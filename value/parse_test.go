package value

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseMultipleValues(t *testing.T) {
	t.Parallel()

	given := "a,b , c,, "

	expected := []string{"a", "b", "c"}

	actual := ParseMultipleValues[string](given)

	assert.Equal(t, expected, actual)
}

func Test_ParsePosition(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		latitude, longitude := "1.1", "2.2"

		expected := Position{
			Latitude:  1.1,
			Longitude: 2.2,
		}

		actual, err := ParsePosition(latitude, longitude)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("invalid latitude", func(t *testing.T) {
		t.Parallel()

		latitude, longitude := "v", "2.2"

		actual, err := ParsePosition(latitude, longitude)

		require.EqualError(t, err, "latitude => strconv.ParseFloat: parsing \"v\": invalid syntax")
		assert.Empty(t, actual)
	})

	t.Run("invalid longitude", func(t *testing.T) {
		t.Parallel()

		latitude, longitude := "11", "v"

		actual, err := ParsePosition(latitude, longitude)

		require.EqualError(t, err, "longitude => strconv.ParseFloat: parsing \"v\": invalid syntax")
		assert.Empty(t, actual)
	})
}

func Test_ParseInt64(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		given := "1"

		expected := int64(1)

		actual, err := ParseInt64(given)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		given := "v"

		actual, err := ParseInt64(given)

		require.EqualError(t, err, "strconv.ParseInt: parsing \"v\": invalid syntax")
		assert.Empty(t, actual)
	})
}

func Test_ParseUint64(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		given := "1"

		expected := uint64(1)

		actual, err := ParseUint64(given)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		given := "v"

		actual, err := ParseUint64(given)

		require.EqualError(t, err, "strconv.ParseUint: parsing \"v\": invalid syntax")
		assert.Empty(t, actual)
	})
}

func Test_ParseFloat64(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		given := "1.1"

		expected := 1.1

		actual, err := ParseFloat64(given)

		require.NoError(t, err)
		assert.InEpsilon(t, expected, actual, 0)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		given := "v"

		actual, err := ParseFloat64(given)

		require.EqualError(t, err, "strconv.ParseFloat: parsing \"v\": invalid syntax")
		assert.Empty(t, actual)
	})
}

func Test_ParseBool(t *testing.T) {
	t.Parallel()

	t.Run("true", func(t *testing.T) {
		t.Parallel()

		given := "1"

		assert.True(t, ParseBool(given))
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()

		given := "0"

		assert.False(t, ParseBool(given))
	})
}

func Test_ParseDate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		given := "2021-01-01"

		expected := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

		actual, err := ParseDate(given)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		given := "v"

		actual, err := ParseDate(given)

		require.EqualError(t, err, "parsing time \"v\" as \"2006-01-02\": cannot parse \"v\" as \"2006\"")
		assert.Empty(t, actual)
	})
}
