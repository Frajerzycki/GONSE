package nse

import (
	"crypto/rand"
	"crypto/sha512"
	"golang.org/x/crypto/hkdf"
	"io"
	"math/big"
)

func GenerateIV(length int) ([]int8, error) {
	if length < 1 {
		return nil, IVZeroLengthError{}
	}

	unsignedIV := make([]byte, length)
	_, err := rand.Read(unsignedIV)
	if err != nil {
		return nil, nil
	}
	IV := make([]int8, length)
	for index, value := range unsignedIV {
		IV[index] = asSigned(value)
	}

	return IV, nil
}

func deriveKey(key *big.Int, salt []byte, dataLength int) (bitsToRotate byte, bytesToRotate int, derivedKey []int8, err error) {
	var bigKeyWithExcludedLength big.Int
	bigKeyWithExcludedLength.Mod(key, big.NewInt(int64(dataLength<<3)))
	keyWithExcludedLength := bigKeyWithExcludedLength.Uint64()
	bitsToRotate = byte(keyWithExcludedLength & 7)
	bytesToRotate = int(keyWithExcludedLength >> 3)
	unsignedDerivedKey := make([]byte, dataLength)
	derivedKey = make([]int8, dataLength)
	io.ReadFull(hkdf.New(sha512.New, key.Bytes(), salt, nil), unsignedDerivedKey)
	for i, v := range unsignedDerivedKey {
		derivedKey[i] = asSigned(v)
	}
	return
}
