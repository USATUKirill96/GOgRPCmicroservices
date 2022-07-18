package location

import "math"

type Distance struct {
	FromLon float64
	FromLat float64
	ToLon   float64
	ToLAt   float64
}

//Meters returns distance in kilometers between two points http://en.wikipedia.org/wiki/Haversine_formula
func (d Distance) Meters() float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = d.FromLat * math.Pi / 180
	lo1 = d.FromLon * math.Pi / 180
	la2 = d.ToLAt * math.Pi / 180
	lo2 = d.ToLon * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * r * math.Asin(math.Sqrt(h))
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
