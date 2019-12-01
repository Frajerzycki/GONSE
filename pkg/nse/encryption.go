package nse

import (
	"github.com/ikcilrep/gonse/internal/bits"
	"github.com/ikcilrep/gonse/internal/errors"
	"math/big"
)

var bigZero *big.Int = big.NewInt(0)

// Encrypt encrypts data with given salt, IV and key using NSE algorithm.
// It returns an error if len(data) < 1, len(data) != len(IV) or if key is not a positive integer.
func Encrypt(data, salt []byte, IV []int8, key *big.Int) ([]int64, error) {
	var err error

	dataLength := len(data)
	IVLength := len(IV)

	switch {
	case dataLength < 1:
		return nil, &errors.NotPositiveDataLengthError{"Data"}
	case dataLength != IVLength:
		return nil, &errors.DifferentIVLengthError{IVLength, dataLength}
	case key.Cmp(big.NewInt(0)) <= 0:
		return nil, &errors.NotPositiveIntegerKeyError{key}
	}

	bitsToRotate, bytesToRotate, derivedKey, err := deriveKey(key, salt, dataLength)
	if err != nil {
		return nil, err
	}
	rotated := bits.RightRotate(data, bitsToRotate, bytesToRotate)
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

// Decrypt decrypts encryptedData with given salt, IV and key using NSE algorithm.
// It returns an error if len(data) < 1, len(data) != len(IV) or if key is not a positive integer.
func Decrypt(encryptedData []int64, salt []byte, IV []int8, key *big.Int) ([]byte, error) {
	var err error

	dataLength := len(encryptedData)
	IVLength := len(IV)

	switch {
	case dataLength < 1:
		return nil, &errors.NotPositiveDataLengthError{"Ciphertext"}
	case dataLength != IVLength:
		return nil, &errors.DifferentIVLengthError{IVLength, dataLength}
	case key.Cmp(bigZero) <= 0:
		return nil, &errors.NotPositiveIntegerKeyError{key}
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
		rotated[index] = bits.AsUnsigned(int8(((encryptedData64[index]+((derivedKey64[index]*sum3)<<1))*sum1 - ((derivedKey64[index] * sum2) << 1)) / sum1Square))
	}
	return bits.LeftRotate(rotated, bitsToRotate, bytesToRotate), nil
}
