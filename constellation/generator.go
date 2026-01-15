package constellation

import (
	"fmt"
	"math"
	"math/rand"
)

// GenerateConstellation creates satellites in orbital shells + ground stations
func GenerateConstellation(numSatellites, numGroundStations int, orbitRadius, commRange float64) *Constellation {
	c := New()

	// Place satellites in a sphere (simplified - real constellations use orbital planes)
	for i := 0; i < numSatellites; i++ {
		// Distribute points on a sphere using golden spiral
		phi := math.Acos(1 - 2*float64(i+1)/float64(numSatellites+1))
		theta := math.Pi * (1 + math.Sqrt(5)) * float64(i)

		pos := Vec3{
			X: orbitRadius * math.Sin(phi) * math.Cos(theta),
			Y: orbitRadius * math.Sin(phi) * math.Sin(theta),
			Z: orbitRadius * math.Cos(phi),
		}

		c.AddNode(&Node{
			ID:        fmt.Sprintf("SAT-%d", i),
			Type:      Satellite,
			Position:  pos,
			CommRange: commRange,
		})
	}

	// Place ground stations on Earth's surface (radius ~6371 km, we'll use smaller scale)
	earthRadius := orbitRadius * 0.6 // ground is below satellites
	for i := 0; i < numGroundStations; i++ {
		lat := (rand.Float64() - 0.5) * math.Pi     // -90 to 90
		lon := (rand.Float64() - 0.5) * 2 * math.Pi // -180 to 180

		pos := Vec3{
			X: earthRadius * math.Cos(lat) * math.Cos(lon),
			Y: earthRadius * math.Cos(lat) * math.Sin(lon),
			Z: earthRadius * math.Sin(lat),
		}

		c.AddNode(&Node{
			ID:        fmt.Sprintf("GS-%d", i),
			Type:      GroundStation,
			Position:  pos,
			CommRange: commRange,
		})
	}

	c.UpdateLinks()
	return c
}
