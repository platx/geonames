package value

type Position struct {
	// Latitude in decimal degrees (wgs84)
	Latitude float64 `url:"lat"`
	// Longitude in decimal degrees (wgs84)
	Longitude float64 `url:"lng"`
}
