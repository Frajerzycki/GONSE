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
// It returns error if length < 1 or if crypto.rand.Read returns an error. 
func GenerateIV(length int) ([]int8, error) {
	if length < 1 {
		return nil, &errors.NotPositiveDataLengthError{"Initialization vector"}
	}
	unsignedIV := make([]byte, length)
	_, err := rand.Read(unsignedIV)
	if err != nil {
		return nil, err
	}
	IV := make([]int8, length)
	for index, value := range unsignedIV {
		IV[index] = bits.AsSigned(value)
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
	keyCopy := new(big.Int)
	keyCopy.SetBytes(key.Bytes())
	for ok := true; ok; ok = !isNonZeroVector(derivedKey) {
		io.ReadFull(hkdf.New(sha512.New, keyCopy.Bytes(), salt, nil), unsignedDerivedKey)
		for i, v := range unsignedDerivedKey {
			derivedKey[i] = bits.AsSigned(v)
		}
		keyCopy.Add(keyCopy, bigOne)
	}

	return
}
