package bits

func AsSigned(b byte) int8 {
	if b < 128 {
		return int8(b)
	}
	return int8(int16(b) - 256)
}

func AsUnsigned(b int8) byte {
	if b < 0 {
		return byte(int16(b) + 256)
	}
	return byte(b)
}
