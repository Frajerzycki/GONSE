package nse

import (
	"crypto/rand"
	"crypto/sha512"
	"github.com/ikcilrep/gonse/internal/bits"
	"github.com/ikcilrep/gonse/internal/errors"
	"golang.org/x/crypto/hkdf"
	"io"
	"math/big"
)

var bigOne *big.Int = big.NewInt(1)

// GenerateIV generates IV of a given length for NSE algorithm.
// It returns an error if length < 1 or if crypto.rand.Read returns an error.
func GenerateIV(length int, rotatedData, derivedKey []int8) ([]int8, error) {
	if length < 1 {
		return nil, &errors.NotPositiveDataLengthError{"Initialization vector"}
	}

	var unsignedIV []byte
	var IV []int8
	for ok := true; ok; ok = isDifferenceOrthogonal(derivedKey, IV, rotatedData) {
		unsignedIV = make([]byte, length)
		_, err := rand.Read(unsignedIV)
		if err != nil {
			return nil, err
		}
		IV = make([]int8, length)
		for index, value := range unsignedIV {
			IV[index] = bits.AsSigned(value)
		}
	}

	return IV, nil
}

func isDifferenceOrthogonal(derivedKey, IV, rotatedData []int8) bool {
	var sum int64 = 0
	for index, keyElement := range derivedKey {
		sum += int64(keyElement) * (int64(rotatedData[index]) - int64(IV[index]))
	}
	return sum == 0
}

func isZeroVector(vector []int8) bool {
	for _, value := range vector {
		if value != 0 {
			return false
		}
	}
	return true
}

func deriveKey(key *big.Int, salt []byte, dataLength int) (bitsToRotate byte, bytesToRotate int, derivedKey []int8, err error) {
	var bigKeyWithExcludedLength big.Int
	bigKeyWithExcludedLength.Mod(key, big.NewInt(int64(dataLength<<3)))
	keyWithExcludedLength := bigKeyWithExcludedLength.Uint64()
	bitsToRotate = byte(keyWithExcludedLength & 7)
	bytesToRotate = int(keyWithExcludedLength >> 3)
	unsignedDerivedKey := make([]byte, dataLength)
	derivedKey = make([]int8, dataLength)
	keyCopy := new(big.Int)
	keyCopy.SetBytes(key.Bytes())
	for ok := true; ok; ok = isZeroVector(derivedKey) {
		io.ReadFull(hkdf.New(sha512.New, keyCopy.Bytes(), salt, nil), unsignedDerivedKey)
		for i, v := range unsignedDerivedKey {
			derivedKey[i] = bits.AsSigned(v)
		}
		keyCopy.Add(keyCopy, bigOne)
	}

	return
}
