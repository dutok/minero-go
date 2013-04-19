package simplex

import (
	"math"
)

// Skewing and unskewing factors for 2, 3, and 4 dimensions
var (
	f2 float64 = 0.5 * (math.Sqrt(3.0) - 1.0)
	g2 float64 = (3.0 - math.Sqrt(3.0)) / 6.0
	f3 float64 = 1.0 / 3.0
	g3 float64 = 1.0 / 6.0
	F4 float64 = (math.Sqrt(5.0) - 1.0) / 4.0
	G4 float64 = (5.0 - math.Sqrt(5.0)) / 20.0
)

// Faster than int(math.Floor(x)), check out simplex_test.go
func fastfloor(x float64) int {
	var xi = int(x)
	if x < float64(xi) {
		return xi - 1
	}
	return xi
}

func dot2(g grad, x, y float64) float64 {
	return g.x*x + g.y*y
}

func dot3(g grad, x, y, z float64) float64 {
	return g.x*x + g.y*y + g.z*z
}

func dot4(g grad, x, y, z, w float64) float64 {
	return g.x*x + g.y*y + g.z*z + g.w*w
}

// At2d computes 2D simplex noise.
func At2d(x, y float64) float64 {
	// Noise contributions from the three corners
	var n0, n1, n2 float64

	// Skew the input space to determine which simplex cell we're in
	var s float64 = (x + y) * f2 // Hairy factor for 2D
	var i int = fastfloor(x + s)
	var j int = fastfloor(y + s)

	var t = float64(i+j) * g2

	// Unskew the cell origin back to (x,y) space
	var ox0 = float64(i) - t
	var oy0 = float64(j) - t

	// The x,y distances from the cell origin
	var x0 float64 = x - ox0
	var y0 float64 = y - oy0

	// For the 2D case, the simplex shape is an equilateral triangle.
	// Determine which simplex we are in.
	var i1, j1 int // Offsets for second (middle) corner of simplex in (i,j) coords
	if x0 > y0 {
		// lower triangle, XY order: (0,0)->(1,0)->(1,1)
		i1, j1 = 1, 0
	} else {
		// upper triangle, YX order: (0,0)->(0,1)->(1,1)
		i1, j1 = 0, 1
	}

	// A step of (1,0) in (i,j) means a step of (1-c,-c) in (x,y), and
	// a step of (0,1) in (i,j) means a step of (-c,1-c) in (x,y), where
	// c = (3-sqrt(3))/6

	// Offsets for middle corner in (x,y) unskewed coords
	var x1 float64 = x0 - float64(i1) + g2
	var y1 float64 = y0 - float64(j1) + g2

	// Offsets for last corner in (x,y) unskewed coords
	var x2 float64 = x0 - 1.0 + 2.0*g2
	var y2 float64 = y0 - 1.0 + 2.0*g2

	// Work out the hashed gradient indices of the three simplex corners
	var ii = i & 0xff
	var jj = j & 0xff
	var gi0 int = permMod12[ii+perm[jj]]
	var gi1 int = permMod12[ii+i1+perm[jj+j1]]
	var gi2 int = permMod12[ii+1+perm[jj+1]]

	// Calculate the contribution from the three corners
	var t0 float64 = 0.5 - x0*x0 - y0*y0
	if t0 < 0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * dot2(grad3[gi0], x0, y0) // (x,y) of grad3 used for 2D gradient
	}

	var t1 float64 = 0.5 - x1*x1 - y1*y1
	if t1 < 0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * dot2(grad3[gi1], x1, y1)
	}

	var t2 float64 = 0.5 - x2*x2 - y2*y2
	if t2 < 0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * dot2(grad3[gi2], x2, y2)
	}

	// Sum and scale to [-1,1]
	return 70.0 * (n0 + n1 + n2)
}

// At3d computes 3D simplex noise.
func At3d(x, y, z float64) float64 {
	// Noise contributions from the four corners
	var n0, n1, n2, n3 float64

	// Skew the input space to determine which simplex cell we're in
	var s float64 = (x + y + z) * f3 // Very nice and simple skew factor for 3D
	var i = int(fastfloor(x + s))
	var j = int(fastfloor(y + s))
	var k = int(fastfloor(z + s))

	var t = float64(i+j+k) * g3

	// Unskew the cell origin back to (x,y,z) space
	var ox0 = float64(i) - t
	var oy0 = float64(j) - t
	var oz0 = float64(k) - t

	// The x,y,z distances from the cell origin
	var x0 float64 = x - ox0
	var y0 float64 = y - oy0
	var z0 float64 = z - oz0

	// For the 3D case, the simplex shape is a slightly irregular tetrahedron.
	// Determine which simplex we are in.
	var i1, j1, k1 int // Offsets for second corner of simplex in (i,j,k) coords
	var i2, j2, k2 int // Offsets for third corner of simplex in (i,j,k) coords

	if x0 >= y0 {
		if y0 >= z0 {
			i1, j1, k1 = 1, 0, 0
			i2, j2, k2 = 1, 1, 0
		} else if x0 >= z0 {
			i1, j1, k1 = 1, 0, 0
			i2, j2, k2 = 1, 0, 1
		} else {
			i1, j1, k1 = 0, 0, 1
			i2, j2, k2 = 1, 0, 1
		}
	} else {
		if y0 < z0 {
			i1, j1, k1 = 0, 0, 1
			i2, j2, k2 = 0, 1, 1
		} else if x0 < z0 {
			i1, j1, k1 = 0, 1, 0
			i2, j2, k2 = 0, 1, 1
		} else {
			i1, j1, k1 = 0, 1, 0
			i2, j2, k2 = 1, 1, 0
		}
	}

	// A step of (1,0,0) in (i,j,k) means a step of (1-c,-c,-c) in (x,y,z),
	// a step of (0,1,0) in (i,j,k) means a step of (-c,1-c,-c) in (x,y,z), and
	// a step of (0,0,1) in (i,j,k) means a step of (-c,-c,1-c) in (x,y,z), where
	// c = 1/6.

	// Offsets for second corner in (x,y,z) coords
	var x1 = x0 - float64(i1) + g3
	var y1 = y0 - float64(j1) + g3
	var z1 = z0 - float64(k1) + g3

	// Offsets for third corner in (x,y,z) coords
	var x2 = x0 - float64(i2) + 2.0*g3
	var y2 = y0 - float64(j2) + 2.0*g3
	var z2 = z0 - float64(k2) + 2.0*g3

	// Offsets for last corner in (x,y,z) coords
	var x3 = x0 - 1.0 + 3.0*g3
	var y3 = y0 - 1.0 + 3.0*g3
	var z3 = z0 - 1.0 + 3.0*g3

	// Work out the hashed gradient indices of the four simplex corners
	var ii int = i & 0xff
	var jj int = j & 0xff
	var kk int = k & 0xff
	var gi0 int = permMod12[ii+perm[jj+perm[kk]]]
	var gi1 int = permMod12[ii+i1+perm[jj+j1+perm[kk+k1]]]
	var gi2 int = permMod12[ii+i2+perm[jj+j2+perm[kk+k2]]]
	var gi3 int = permMod12[ii+1+perm[jj+1+perm[kk+1]]]

	// Calculate the contribution from the four corners
	var t0 float64 = 0.5 - x0*x0 - y0*y0 - z0*z0
	if t0 < 0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * dot3(grad3[gi0], x0, y0, z0)
	}

	var t1 float64 = 0.5 - x1*x1 - y1*y1 - z1*z1
	if t1 < 0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * dot3(grad3[gi1], x1, y1, z1)
	}

	var t2 float64 = 0.5 - x2*x2 - y2*y2 - z2*z2
	if t2 < 0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * dot3(grad3[gi2], x2, y2, z2)
	}

	var t3 float64 = 0.5 - x3*x3 - y3*y3 - z3*z3
	if t3 < 0 {
		n3 = 0.0
	} else {
		t3 *= t3
		n3 = t3 * t3 * dot3(grad3[gi3], x3, y3, z3)
	}

	// Sum and scale to [-1,1]
	return 32.0 * (n0 + n1 + n2 + n3)
}

// At4d computes 4D simplex noise.
func At4d(x, y, z, w float64) float64 {
	// Noise contributions from the five corners
	var n0, n1, n2, n3, n4 float64

	// Skew the (x,y,z,w) space to determine which cell of 24 simplices we're in
	var s float64 = (x + y + z + w) * F4 // Factor for 4D skewing
	var i int = fastfloor(x + s)
	var j int = fastfloor(y + s)
	var k int = fastfloor(z + s)
	var l int = fastfloor(w + s)

	var t = float64(i+j+k+l) * G4 // Factor for 4D unskewing

	// Unskew the cell origin back to (x,y,z,w) space
	var ox0 = float64(i) - t
	var oy0 = float64(j) - t
	var oz0 = float64(k) - t
	var ow0 = float64(l) - t

	// The x,y,z,w distances from the cell origin
	var x0 float64 = x - ox0
	var y0 float64 = y - oy0
	var z0 float64 = z - oz0
	var w0 float64 = w - ow0

	// For the 4D case, the simplex is a 4D shape I won't even try to describe.
	// To find out which of the 24 possible simplices we're in, we need to
	// determine the magnitude ordering of x0, y0, z0 and w0.
	var rankx, ranky, rankz, rankw int

	// Six pair-wise comparisons are performed between each possible pair
	// of the four coordinates, and the results are used to rank the numbers.
	if x0 > y0 {
		rankx++
	} else {
		ranky++
	}
	if x0 > z0 {
		rankx++
	} else {
		rankz++
	}
	if x0 > w0 {
		rankx++
	} else {
		rankw++
	}
	if y0 > z0 {
		ranky++
	} else {
		rankz++
	}
	if y0 > w0 {
		ranky++
	} else {
		rankw++
	}
	if z0 > w0 {
		rankz++
	} else {
		rankw++
	}

	// The integer offsets for the second, third and fourth simplex corner
	var i1, j1, k1, l1 int
	var i2, j2, k2, l2 int
	var i3, j3, k3, l3 int

	// simplex[c] is a 4-vector with the numbers 0, 1, 2 and 3 in some order.
	// Many values of c will never occur, since e.g. x>y>z>w makes x<z, y<w and x<w
	// impossible. Only the 24 indices which have non-zero entries make any sense.
	// We use a thresholding to set the coordinates in turn from the largest magnitude.

	// Rank 3 denotes the largest coordinate.
	if rankx >= 3 {
		i1 = 1
	} else {
		i1 = 0
	}
	if ranky >= 3 {
		j1 = 1
	} else {
		j1 = 0
	}
	if rankz >= 3 {
		k1 = 1
	} else {
		k1 = 0
	}
	if rankw >= 3 {
		l1 = 1
	} else {
		l1 = 0
	}

	// Rank 2 denotes the second largest coordinate.
	if rankx >= 2 {
		i2 = 1
	} else {
		i2 = 0
	}
	if ranky >= 2 {
		j2 = 1
	} else {
		j2 = 0
	}
	if rankz >= 2 {
		k2 = 1
	} else {
		k2 = 0
	}
	if rankw >= 2 {
		l2 = 1
	} else {
		l2 = 0
	}

	// Rank 1 denotes the second smallest coordinate.
	if rankx >= 1 {
		i3 = 1
	} else {
		i3 = 0
	}
	if ranky >= 1 {
		j3 = 1
	} else {
		j3 = 0
	}
	if rankz >= 1 {
		k3 = 1
	} else {
		k3 = 0
	}
	if rankw >= 1 {
		l3 = 1
	} else {
		l3 = 0
	}

	// The fifth corner has all coordinate offsets = 1, so no need to compute that.

	// Offsets for second corner in (x,y,z,w) coords
	var x1 = x0 - float64(i1) + G4
	var y1 = y0 - float64(j1) + G4
	var z1 = z0 - float64(k1) + G4
	var w1 = w0 - float64(l1) + G4

	// Offsets for third corner in (x,y,z,w) coords
	var x2 = x0 - float64(i2) + 2.0*G4
	var y2 = y0 - float64(j2) + 2.0*G4
	var z2 = z0 - float64(k2) + 2.0*G4
	var w2 = w0 - float64(l2) + 2.0*G4

	// Offsets for fourth corner in (x,y,z,w) coords
	var x3 = x0 - float64(i3) + 3.0*G4
	var y3 = y0 - float64(j3) + 3.0*G4
	var z3 = z0 - float64(k3) + 3.0*G4
	var w3 = w0 - float64(l3) + 3.0*G4

	// Offsets for last corner in (x,y,z,w) coords
	var x4 = x0 - 1.0 + 4.0*G4
	var y4 = y0 - 1.0 + 4.0*G4
	var z4 = z0 - 1.0 + 4.0*G4
	var w4 = w0 - 1.0 + 4.0*G4

	// Work out the hashed gradient indices of the five simplex corners
	var ii int = i & 0xff
	var jj int = j & 0xff
	var kk int = k & 0xff
	var ll int = l & 0xff
	var gi0 int = perm[ii+perm[jj+perm[kk+perm[ll]]]] % 32
	var gi1 int = perm[ii+i1+perm[jj+j1+perm[kk+k1+perm[ll+l1]]]] % 32
	var gi2 int = perm[ii+i2+perm[jj+j2+perm[kk+k2+perm[ll+l2]]]] % 32
	var gi3 int = perm[ii+i3+perm[jj+j3+perm[kk+k3+perm[ll+l3]]]] % 32
	var gi4 int = perm[ii+1+perm[jj+1+perm[kk+1+perm[ll+1]]]] % 32

	// Calculate the contribution from the five corners
	var t0 float64 = 0.5 - x0*x0 - y0*y0 - z0*z0 - w0*w0
	if t0 < 0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * dot4(grad4[gi0], x0, y0, z0, w0)
	}

	var t1 float64 = 0.5 - x1*x1 - y1*y1 - z1*z1 - w1*w1
	if t1 < 0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * dot4(grad4[gi1], x1, y1, z1, w1)
	}

	var t2 float64 = 0.5 - x2*x2 - y2*y2 - z2*z2 - w2*w2
	if t2 < 0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * dot4(grad4[gi2], x2, y2, z2, w2)
	}

	var t3 float64 = 0.5 - x3*x3 - y3*y3 - z3*z3 - w3*w3
	if t3 < 0 {
		n3 = 0.0
	} else {
		t3 *= t3
		n3 = t3 * t3 * dot4(grad4[gi3], x3, y3, z3, w3)
	}

	var t4 float64 = 0.5 - x4*x4 - y4*y4 - z4*z4 - w4*w4
	if t4 < 0 {
		n4 = 0.0
	} else {
		t4 *= t4
		n4 = t4 * t4 * dot4(grad4[gi4], x4, y4, z4, w4)
	}

	// Sum and scale to [-1,1]
	return 27.0 * (n0 + n1 + n2 + n3 + n4)
}
