package nse

import "math/big"

var bigZero *big.Int = big.NewInt(0)

func rightRotateBits(b1, b2 byte, bitsToRotate byte) byte {
	return (b1 << (8 - bitsToRotate)) | (b2 >> bitsToRotate)
}

func leftRotateBits(b1, b2 byte, bitsToRotate byte) byte {
	return (b1 << bitsToRotate) | (b2 >> (8 - bitsToRotate))
}

func rightRotate(dataToRotate []byte, bitsToRotate byte, bytesToRotate int) []int8 {
	length := len(dataToRotate)
	bitRotated := make([]int8, length)
	byteRotated := make([]byte, length)

	for index := range dataToRotate[:bytesToRotate] {
		byteRotated[index] = dataToRotate[length+index-bytesToRotate]
	}

	for index := range dataToRotate[bytesToRotate:] {
		byteRotated[index+bytesToRotate] = dataToRotate[index]
	}

	bitRotated[0] = asSigned(rightRotateBits(byteRotated[length-1], byteRotated[0], bitsToRotate))

	for index := 1; index < length; index++ {
		bitRotated[index] = asSigned(rightRotateBits(byteRotated[index-1], byteRotated[index], bitsToRotate))
	}

	return bitRotated
}

func leftRotate(dataToRotate []byte, bitsToRotate byte, bytesToRotate int) []byte {
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

func asSigned(b byte) int8 {
	if b < 128 {
		return int8(b)
	}
	return int8(int16(b) - 256)
}

func asUnsigned(b int8) byte {
	if b < 0 {
		return byte(int16(b) + 256)
	}
	return byte(b)
}

func Encrypt(data, salt []byte, IV []int8, key *big.Int) ([]int64, error) {
	var err error
	switch {
	case data == nil:
		err = NilArgumentError{"Data"}
	case IV == nil:
		err = NilArgumentError{"Initalization vector"}
	case key == nil:
		err = NilArgumentError{"Key"}
	}

	dataLength := len(data)
	IVLength := len(IV)

	switch {
	case dataLength < 1:
		err = NilArgumentError{"Data"}
	case dataLength != IVLength:
		err = DifferentIVLengthError{IVLength, dataLength}
	case key.Cmp(big.NewInt(0)) <= 0:
		err = NotPositiveIntegerKeyError{key}
	}
	if err != nil {
		return nil, err
	}

	bitsToRotate, bytesToRotate, derivedKey, err := deriveKey(key, salt, dataLength)
	if err != nil {
		return nil, err
	}
	rotated := rightRotate(data, bitsToRotate, bytesToRotate)
	rotated64, IV64, derivedKey64 := make([]int64, dataLength), make([]int64, dataLength), make([]int64, dataLength)

	var sum1, sum2 int64 = 0, 0
	for index := 0; index < dataLength; index++ {
		rotated64[index], IV64[index], derivedKey64[index] = int64(rotated[index]), int64(IV[index]), int64(derivedKey[index])
		sum1 += derivedKey64[index] * derivedKey64[index]
		sum2 += derivedKey64[index] * (rotated64[index] - IV64[index])
	}
	encryptedData := make([]int64, dataLength)

	for index := range encryptedData {
		encryptedData[index] = rotated64[index]*sum1 - ((derivedKey64[index] * sum2) << 1)
	}

	return encryptedData, nil
}

func Decrypt(encryptedData []int64, salt []byte, IV []int8, key *big.Int) ([]byte, error) {
	var err error
	switch {
	case encryptedData == nil:
		err = NilArgumentError{"Encrypted data"}
	case IV == nil:
		err = NilArgumentError{"Initalization vector"}
	case key == nil:
		err = NilArgumentError{"Key"}
	}

	dataLength := len(encryptedData)
	IVLength := len(IV)

	switch {
	case dataLength < 1:
		err = NilArgumentError{"Encrypted data"}
	case dataLength != IVLength:
		err = DifferentIVLengthError{IVLength, dataLength}
	case key.Cmp(bigZero) <= 0:
		err = NotPositiveIntegerKeyError{key}
	}

	if err != nil {
		return nil, err
	}

	bitsToRotate, bytesToRotate, derivedKey, err := deriveKey(key, salt, dataLength)
	if err != nil {
		return nil, err
	}

	rotated := make([]byte, dataLength)
	encryptedData64, IV64, derivedKey64 := make([]int64, dataLength), make([]int64, dataLength), make([]int64, dataLength)

	var sum1, sum2, sum3 int64 = 0, 0, 0
	for index := 0; index < dataLength; index++ {
		encryptedData64[index], IV64[index], derivedKey64[index] = int64(encryptedData[index]), int64(IV[index]), int64(derivedKey[index])
		sum1 += derivedKey64[index] * derivedKey64[index]
		sum2 += derivedKey64[index] * encryptedData64[index]
		sum3 += derivedKey64[index] * IV64[index]
	}

	sum1Square := sum1 * sum1

	for index := range encryptedData {
		rotated[index] = asUnsigned(int8(((encryptedData64[index]+((derivedKey64[index]*sum3)<<1))*sum1 - ((derivedKey64[index] * sum2) << 1)) / sum1Square))
	}
	return leftRotate(rotated, bitsToRotate, bytesToRotate), nil
}
