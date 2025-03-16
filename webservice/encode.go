package webservice

import (
	"net/url"
	"reflect"
	"strconv"
)

// URLEncoder is a struct that can encode a struct into url.Values.
type URLEncoder struct {
	values url.Values
}

func NewURLEncoder(values url.Values) *URLEncoder {
	return &URLEncoder{values: values}
}

// Encode encodes a struct into url.Values using reflection.
func (enc *URLEncoder) Encode(v any) {
	val := reflect.ValueOf(v)
	valType := reflect.TypeOf(v)

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		value := val.Field(i)

		key := field.Tag.Get("url")
		if key == "" || key == "-" {
			continue
		}

		if field.Type.Kind() == reflect.Struct && key == ",dive" {
			enc.Encode(value.Interface())

			continue
		}

		enc.encodeValue(key, value)
	}
}

// encodeValue encodes a value into url.Values using reflection type switch.
func (enc *URLEncoder) encodeValue(key string, value reflect.Value) {
	switch value.Kind() {
	case reflect.String:
		enc.encodeString(key, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		enc.encodeInt(key, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		enc.encodeUint(key, value)
	case reflect.Float32, reflect.Float64:
		enc.encodeFloat(key, value)
	case reflect.Bool:
		enc.encodeBool(key, value)
	case reflect.Slice:
		enc.encodeSlice(key, value)
	case
		reflect.Invalid,
		reflect.Uintptr,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Array,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Struct,
		reflect.UnsafePointer:
		return
	}
}

// encodeString encodes a string value into url.Values. Empty string is not encoded.
func (enc *URLEncoder) encodeString(key string, value reflect.Value) {
	if value.String() != "" {
		enc.values.Set(key, value.String())
	}
}

// encodeInt encodes a signed integer value into url.Values. Zero value is not encoded.
func (enc *URLEncoder) encodeInt(key string, value reflect.Value) {
	if value.Int() != 0 {
		enc.values.Set(key, strconv.FormatInt(value.Int(), 10))
	}
}

// encodeUint encodes an unsigned integer value into url.Values. Zero value is not encoded.
func (enc *URLEncoder) encodeUint(key string, value reflect.Value) {
	if value.Uint() != 0 {
		enc.values.Set(key, strconv.FormatUint(value.Uint(), 10))
	}
}

// encodeFloat encodes a float value into url.Values. Zero value is not encoded.
func (enc *URLEncoder) encodeFloat(key string, value reflect.Value) {
	if value.Float() != 0 {
		enc.values.Set(key, strconv.FormatFloat(value.Float(), 'f', -1, 32))
	}
}

// encodeBool encodes a boolean value into url.Values. False value is not encoded.
func (enc *URLEncoder) encodeBool(key string, value reflect.Value) {
	if value.Bool() {
		enc.values.Set(key, "true")
	}
}

// encodeBool encodes a slice value into url.Values.
func (enc *URLEncoder) encodeSlice(key string, value reflect.Value) {
	for j := 0; j < value.Len(); j++ {
		item := value.Index(j)
		enc.values.Add(key, item.String())
	}
}
