package simplex

// MultiAt2d computes 2D Multi-Octave Simplex noise. For each octave, a higher
// frequency/lower amplitude function will be added to the original. The higher
// the persistence [0, 1], the more of each succeeding octave will be added.
func MultiAt2d(oct int, p, f, x, y float64) float64 {
	var (
		total        float64
		amplitude    float64 = 1.0
		maxAmplitude float64
	)

	for i := 0; i < oct; i++ {
		total += At2d(x*f, y*f) * amplitude
		f *= 2.0
		maxAmplitude += amplitude
		amplitude *= p
	}

	return total / maxAmplitude
}

// MultiAt3d computes 3D Multi-Octave Simplex noise. For each octave, a higher
// frequency/lower amplitude function will be added to the original. The higher
// the persistence [0, 1], the more of each succeeding octave will be added.
func MultiAt3d(oct int, p, f, x, y, z float64) float64 {
	var (
		total        float64
		amplitude    float64 = 1.0
		maxAmplitude float64
	)

	for i := 0; i < oct; i++ {
		total += At3d(x*f, y*f, z*f) * amplitude
		f *= 2.0
		maxAmplitude += amplitude
		amplitude *= p
	}

	return total / maxAmplitude
}

// MultiAt4d computes 4D Multi-Octave Simplex noise. For each octave, a higher
// frequency/lower amplitude function will be added to the original. The higher
// the persistence [0, 1], the more of each succeeding octave will be added.
func MultiAt4d(oct int, p, f, x, y, z, w float64) float64 {
	var (
		total        float64
		amplitude    float64 = 1.0
		maxAmplitude float64
	)

	for i := 0; i < oct; i++ {
		total += At4d(x*f, y*f, z*f, w*f) * amplitude
		f *= 2.0
		maxAmplitude += amplitude
		amplitude *= p
	}

	return total / maxAmplitude
}
