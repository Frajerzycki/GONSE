package bits

func rightRotateBits(b1, b2 byte, bitsToRotate byte) byte {
	return (b1 << (8 - bitsToRotate)) | (b2 >> bitsToRotate)
}

func leftRotateBits(b1, b2 byte, bitsToRotate byte) byte {
	return (b1 << bitsToRotate) | (b2 >> (8 - bitsToRotate))
}

func RightRotate(dataToRotate []byte, bitsToRotate byte, bytesToRotate int) []int8 {
	length := len(dataToRotate)
	bitRotated := make([]int8, length)
	byteRotated := make([]byte, length)

	for index := range dataToRotate[:bytesToRotate] {
		byteRotated[index] = dataToRotate[length+index-bytesToRotate]
	}

	for index := range dataToRotate[bytesToRotate:] {
		byteRotated[index+bytesToRotate] = dataToRotate[index]
	}

	bitRotated[0] = AsSigned(rightRotateBits(byteRotated[length-1], byteRotated[0], bitsToRotate))

	for index := 1; index < length; index++ {
		bitRotated[index] = AsSigned(rightRotateBits(byteRotated[index-1], byteRotated[index], bitsToRotate))
	}

	return bitRotated
}

func LeftRotate(dataToRotate []byte, bitsToRotate byte, bytesToRotate int) []byte {
	length := len(dataToRotate)
	bitRotated := make([]byte, length)
	byteRotated := make([]byte, length)

	{
		lastIndex := length - 1
		for index := range dataToRotate[:lastIndex] {
			bitRotated[index] = leftRotateBits(dataToRotate[index], dataToRotate[index+1], bitsToRotate)
		}

		bitRotated[lastIndex] = leftRotateBits(dataToRotate[lastIndex], dataToRotate[0], bitsToRotate)
	}

	{
		limit := length - bytesToRotate
		for index := range dataToRotate[:limit] {
			byteRotated[index] = bitRotated[index+bytesToRotate]
		}

		for index := limit; index < length; index++ {
			byteRotated[index] = bitRotated[index+bytesToRotate-length]
		}
	}

	return byteRotated
}
