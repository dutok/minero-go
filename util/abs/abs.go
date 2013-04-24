package abs

func Pos(i float64) int32 {
	return int32(i * 32)
}

func Look(l float32) int8 {
	return int8(int(l / 180.0 * 128))
}

func RealPos(i int32) float64 {
	return float64(i << 5)
}

func RealLook(l int8) float32 {
	return float32(l) / 128 * 180.0
}
