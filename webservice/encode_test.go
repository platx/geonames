package webservice

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_URLEncoder_Encode(t *testing.T) {
	t.Parallel()

	values := url.Values{}
	encoder := NewURLEncoder(values)

	input := testStruct{
		StringField:  "test",
		IntField:     42,
		UintField:    99,
		FloatField:   3.14,
		BoolField:    true,
		SliceField:   []string{"a", "b", "c"},
		IgnoredField: "ignored",
		DiveField: diveStruct{
			InnerField:   "inner",
			IgnoredField: "ignored",
		},
	}

	encoder.Encode(input)

	assert.Equal(t, "test", values.Get("string_field"))
	assert.Equal(t, "42", values.Get("int_field"))
	assert.Equal(t, "99", values.Get("uint_field"))
	assert.Equal(t, "3.14", values.Get("float_field"))
	assert.Equal(t, "true", values.Get("bool_field"))
	assert.ElementsMatch(t, []string{"a", "b", "c"}, values["slice_field"])
	assert.Empty(t, values.Get("ignored_field"))
	assert.Equal(t, "inner", values.Get("inner_field"))
	assert.Empty(t, values.Get("ignored_struct"))
}

func Test_URLEncoder_Encode_EmptyValues(t *testing.T) {
	t.Parallel()

	values := url.Values{}
	encoder := NewURLEncoder(values)

	input := testStruct{}
	encoder.Encode(input)

	require.Empty(t, values)
}

type testStruct struct {
	StringField   string     `url:"string_field"`
	IntField      int        `url:"int_field"`
	UintField     uint       `url:"uint_field"`
	FloatField    float64    `url:"float_field"`
	BoolField     bool       `url:"bool_field"`
	IgnoredField  string     `url:"-"`
	SliceField    []string   `url:"slice_field"`
	DiveField     diveStruct `url:",dive"`
	IgnoredStruct struct{}   `url:"ignored_struct"`
}

type diveStruct struct {
	InnerField   string `url:"inner_field"`
	IgnoredField string `url:"-"`
}
