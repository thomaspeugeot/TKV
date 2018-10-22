package barneshut

import (
	"math"
	"math/rand"
	"sync/atomic"

	"github.com/thomaspeugeot/tkv/quadtree"
)

// get modulo distance between alpha and beta.
//
// alpha and beta are between 0.0 and 1.0
// the modulo distance cannot be above 0.5
func getModuloDistance(alpha, beta float64) (dist float64) {

	dist = beta - alpha
	if dist > 0.5 {
		dist -= 1.0
	}
	if dist < -0.5 {
		dist += 1.0
	}

	return dist
}

// get modulo distance from border
//
// alpha and beta are between 0.0 and 1.0
func getModuloDistanceFromBoder(A *quadtree.Body) (distX, distY float64) {
	return getModuloDistance(A.X, 0.0), getModuloDistance(A.Y, 0.0)
}

// get modulo distance between alpha and beta in a given village.
//
// alpha and beta are between left and rigth
// the modulo distance cannot be above (rigth-left) /2
func getModuloDistanceLocal(alpha, beta, left, rigth float64) (dist float64) {

	dist = beta - alpha
	maxDist := rigth - left
	halfMaxDist := maxDist / 2.0
	if dist > halfMaxDist {
		dist -= maxDist
	}
	if dist < -halfMaxDist {
		dist += maxDist
	}

	return dist
}

// function used to spread bodies randomly on
// the unit square
func SpreadOnCircle(bodies *[]quadtree.Body) {
	for idx := range *bodies {

		body := &((*bodies)[idx])

		radius := rand.Float64()
		angle := 2.0 * math.Pi * rand.Float64()

		if idx%2 == 0 {
			body.X = 0.2
			body.Y = 0.7
			radius *= 0.15
		} else {
			body.X = 0.6
			body.Y = 0.4
			radius *= 0.25
		}

		body.M = 0.1000000
		body.X += math.Cos(angle) * radius
		body.Y += math.Sin(angle) * radius
	}
}

// compute modulo distance
func getModuloDistanceBetweenBodies(A, B *quadtree.Body) float64 {

	x := getModuloDistance(B.X, A.X)
	y := getModuloDistance(B.Y, A.Y)

	distSquared := (x*x + y*y)

	return math.Sqrt(distSquared)
}

// compute mirror vector berween bodies A & B
//
// x == -1, B's x position is mirrored relative to x=0
// x == 0, B's x position is unchanged
// x == 1, B's x position is mirrored relative to x=1
//
// idem for y for B's y position
func getVectorBetweenBodiesWithMirror(A, B *quadtree.Body, x, y int) (vX, xY float64) {

	xB := B.X
	yB := B.Y

	if x == -1 {
		xB = 0.0 - xB
	}
	if x == 1 {
		xB = 2.0 - xB
	}

	if y == -1 {
		yB = 0.0 - yB
	}
	if y == 1 {
		yB = 2.0 - yB
	}

	return xB - A.X, yB - A.Y
}

// compute distance between A and B with xM, yM transformation
func getDistanceBetweenBodiesWithMirror(A, B *quadtree.Body, xM, yM int) float64 {

	xV, yV := getVectorBetweenBodiesWithMirror(A, B, xM, yM)
	distSquared := (xV*xV + yV*yV)

	res := math.Sqrt(distSquared)

	// Trace.Printf("A %f %f B %f %f d %f", A.X, A.Y, B.X, B.Y, res)

	return res
}

// compute repulsion force vector between body A and body B
// applied to body A
// proportional to the inverse of the distance squared
// return x, y of repulsion vector and distance between A & B
// return energy as the repulsion energy
func getRepulsionVector(A, B *quadtree.Body, xM, yM int) (x, y, energy float64) {

	atomic.AddUint64(&nbComputationPerStep, 1)

	// Trace.Printf("getRepulsionVector A %f %f B %f %f", A.X, A.Y, B.X, B.Y)

	x, y = getVectorBetweenBodiesWithMirror(A, B, xM, yM)

	distQuared := (x*x + y*y)
	absDistance := math.Sqrt(distQuared + ETA)

	distPow3 := distQuared * absDistance

	// repulsion is proportional to mass
	massCombined := A.M * B.M
	x *= massCombined
	y *= massCombined

	// repulsion is inversly proportional to the square of the distance (1/r2)
	x = -x / distPow3
	y = -y / distPow3

	if absDistance > CutoffDistance {
		x = 0.0
		y = 0.0
	}
	return x, y, massCombined / absDistance
}
