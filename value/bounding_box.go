package value

type BoundingBox struct {
	East  float64 `url:"east"`
	West  float64 `url:"west"`
	North float64 `url:"north"`
	South float64 `url:"south"`
}
