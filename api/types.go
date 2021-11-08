package api

// GeoPoint stores latitude and longitude information
type GeoPoint struct {
	Latitude  float64
	Longitude float64
}

// GeoPoints type alias for slice of GeoPoint
type GeoPoints []GeoPoint

// Add adds a GeoPoint to the GeoPoints slice
func (gp *GeoPoints) Add(point GeoPoint) {
	*gp = append(*gp, point)
}
