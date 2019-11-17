package nse

import (
	"crypto/rand"
	"crypto/sha512"
	"github.com/ikcilrep/gonse/internal/errors"
	"golang.org/x/crypto/hkdf"
	"io"
	"math/big"
)

var bigOne *big.Int = big.NewInt(1)

func GenerateIV(length int) ([]int8, error) {
	if length < 1 {
		return nil, errors.NotPositiveDataLengthError{"Initialization vector"}
	}
	unsignedIV := make([]byte, length)
	_, err := rand.Read(unsignedIV)
	if err != nil {
		return nil, err
	}
	IV := make([]int8, length)
	for index, value := range unsignedIV {
		IV[index] = asSigned(value)
	}

	return IV, nil
}

func isNonZeroVector(vector []int8) bool {
	for _, v := range vector {
		if v != 0 {
			return true
		}
	}
	return false
}

func deriveKey(key *big.Int, salt []byte, dataLength int) (bitsToRotate byte, bytesToRotate int, derivedKey []int8, err error) {
	var bigKeyWithExcludedLength big.Int
	bigKeyWithExcludedLength.Mod(key, big.NewInt(int64(dataLength<<3)))
	keyWithExcludedLength := bigKeyWithExcludedLength.Uint64()
	bitsToRotate = byte(keyWithExcludedLength & 7)
	bytesToRotate = int(keyWithExcludedLength >> 3)
	unsignedDerivedKey := make([]byte, dataLength)
	derivedKey = make([]int8, dataLength)
	keyCopy := *key
	keyCopyRef := &keyCopy
	for ok := true; ok; ok = !isNonZeroVector(derivedKey) {
		io.ReadFull(hkdf.New(sha512.New, keyCopyRef.Bytes(), salt, nil), unsignedDerivedKey)
		for i, v := range unsignedDerivedKey {
			derivedKey[i] = asSigned(v)
		}
		keyCopyRef.Add(keyCopyRef, bigOne)
	}

	return
}
