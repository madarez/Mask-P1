package main

import (
	"math"
)

// reference for most of the calculations here is https://www.movable-type.co.uk/scripts/latlong.html and http://www.edwilliams.org/avform.htm#LL

const EarthRadius = 6371e3                    // (meters) Wikipedia distances from points on the surface to the center range from 6,353 km to 6,384 km, FAI standard for aviation records is 6,371 km
const SafetyRadial float64 = 2. / EarthRadius // (radian distance) safetly distance set by WHO

// convert radiants to degrees
func GetRadiant(val float64) float64 {
	return val * math.Pi / 180
}

// convert degrees to radiants
func GetDegrees(val float64) float64 {
	return val * 180 / math.Pi
}

// calculate the distance of two GPS data locations using Law of Cosines rather than Haversine formula:
// https://gis.stackexchange.com/questions/4906/why-is-law-of-cosines-more-preferable-than-haversine-when-calculating-distance-b
func Distance(g1, g2 *GPSSentence) (d float64) {
	deltaLambda := g2[1] - g1[1]

	return math.Acos(math.Sin(g1[0])*math.Sin(g2[0])+
		math.Cos(g1[0])*math.Cos(g2[0])*math.Cos(deltaLambda)) * EarthRadius
}

// Interpolate direct path between two GPS data points
func InterpolatePath(g1, g2 *GPSSentence, t1, t2, t int) *GPSSentence {
	w1 := float64(t-t1) / float64(t2-t1)
	w2 := float64(t2-t) / float64(t2-t1)
	WeightedMean := func(a, b float64) float64 {
		return a*w1 + b*w2
	}
	lat := WeightedMean(g1[0], g2[0])
	lon := WeightedMean(g1[1], g2[1])
	return &GPSSentence{lat, lon}
}

// proximate a bounding box around the GPS data and return SW corner (min) and NE corner (max)
func ProximateBox(g *GPSSentence) (*GPSSentence, *GPSSentence) {
	var dlat float64
	dlat = SafetyRadial
	dlon := math.Asin(math.Sin(dlat) / math.Cos(g[0])) // to be more accurate: math.Atan2(math.Sin(rd)*math.Cos(g[0]), math.Cos(rd)-math.Sin(g[0])*math.Sin(g[0]))
	return &GPSSentence{g[0] - dlat, g[1] - dlon},
		&GPSSentence{g[0] + dlat, g[1] + dlon}
}
