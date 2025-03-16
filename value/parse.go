package value

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseMultipleValues[T ~string](given string) []T {
	rawValues := strings.Split(given, ",")
	values := make([]T, 0, len(rawValues))

	for _, val := range rawValues {
		if val = strings.TrimSpace(val); val == "" {
			continue
		}

		values = append(values, T(val))
	}

	return values
}

func ParsePosition(latitude string, longitude string) (Position, error) {
	var (
		res Position
		err error
	)

	if latitude != "" {
		res.Latitude, err = strconv.ParseFloat(latitude, 64)
		if err != nil {
			return Position{}, fmt.Errorf("latitude => %w", err)
		}
	}

	if longitude != "" {
		res.Longitude, err = strconv.ParseFloat(longitude, 64)
		if err != nil {
			return Position{}, fmt.Errorf("longitude => %w", err)
		}
	}

	return res, nil
}

func ParseInt64(given string) (int64, error) {
	var (
		res int64
		err error
	)

	if given != "" {
		res, err = strconv.ParseInt(given, 10, 64)
	}

	return res, err
}

func ParseUint64(given string) (uint64, error) {
	var (
		res uint64
		err error
	)

	if given != "" {
		res, err = strconv.ParseUint(given, 10, 64)
	}

	return res, err
}

func ParseFloat64(given string) (float64, error) {
	var (
		res float64
		err error
	)

	if given != "" {
		res, err = strconv.ParseFloat(given, 64)
	}

	return res, err
}

func ParseBool(given string) bool {
	return given == "1"
}

func ParseDate(given string) (time.Time, error) {
	var (
		res time.Time
		err error
	)

	if given != "" {
		res, err = time.Parse(time.DateOnly, given)
	}

	return res, err
}
