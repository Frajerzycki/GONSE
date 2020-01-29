package nse

import (
	"github.com/ikcilrep/gonse/internal/bits"
	"github.com/ikcilrep/gonse/internal/errors"
	"math/big"
)

var bigZero *big.Int = big.NewInt(0)

// Encrypt encrypts data using NSE algorithm with given key derived using DeriveKey function.
// It returns encryptedData, IV and err.
// err != nil if len(data) < 1 or if GenerateIV function returned an error.
func Encrypt(data []byte, key *NSEKey) (encryptedData []int64, IV []int8, err error) {
	dataLength := len(data)

	if dataLength < 1 {
		return nil, nil, &errors.NotPositiveDataLengthError{"Data"}
	}

	rotated := bits.RightRotate(data, key.BitsToRotate, key.BytesToRotate)
	IV, err = GenerateIV(dataLength, rotated, key)
	if err != nil {
		return nil, nil, err
	}
	rotated64, IV64, derivedKey64 := make([]int64, dataLength), make([]int64, dataLength), make([]int64, dataLength)

	var sum1, sum2 int64 = 0, 0
	for index := 0; index < dataLength; index++ {
		rotated64[index], IV64[index], derivedKey64[index] = int64(rotated[index]), int64(IV[index]), int64(key.Data[index])
		sum1 += derivedKey64[index] * derivedKey64[index]
		sum2 += derivedKey64[index] * (rotated64[index] - IV64[index])
	}
	encryptedData = make([]int64, dataLength)

	for index := range encryptedData {
		encryptedData[index] = rotated64[index]*sum1 - ((derivedKey64[index] * sum2) << 1)
	}

	return
}

// Decrypt decrypts encryptedData using NSE algorithm with given IV and key derived using DeriveKey function.
// It returns decryptedData and err.
// err != nil if len(data) < 1, len(data) != len(IV).
func Decrypt(encryptedData []int64, IV []int8, key *NSEKey) (decryptedData []byte, err error) {
	dataLength := len(encryptedData)
	IVLength := len(IV)

	switch {
	case dataLength < 1:
		return nil, &errors.NotPositiveDataLengthError{"Ciphertext"}
	case dataLength != IVLength:
		return nil, &errors.DifferentIVLengthError{IVLength, dataLength}
	}

	rotated := make([]byte, dataLength)
	encryptedData64, IV64, derivedKey64 := make([]int64, dataLength), make([]int64, dataLength), make([]int64, dataLength)

	var sum1, sum2, sum3 int64 = 0, 0, 0
	for index := 0; index < dataLength; index++ {
		encryptedData64[index], IV64[index], derivedKey64[index] = int64(encryptedData[index]), int64(IV[index]), int64(key.Data[index])
		sum1 += derivedKey64[index] * derivedKey64[index]
		sum2 += derivedKey64[index] * encryptedData64[index]
		sum3 += derivedKey64[index] * IV64[index]
	}

	sum1Square := sum1 * sum1

	for index := range encryptedData {
		rotated[index] = bits.AsUnsigned(int8(((encryptedData64[index]+((derivedKey64[index]*sum3)<<1))*sum1 - ((derivedKey64[index] * sum2) << 1)) / sum1Square))
	}
	return bits.LeftRotate(rotated, key.BitsToRotate, key.BytesToRotate), nil
}
