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

type NSEKey struct {
	Data          []int8
	BitsToRotate  byte
	BytesToRotate int
}

var bigOne *big.Int = big.NewInt(1)

// GenerateIV generates IV of given length for NSE algorithm.
// It returns an error if length < 1 or if crypto.rand.Read returns an error.
func GenerateIV(length int, rotatedData []int8, key *NSEKey) ([]int8, error) {
	if length < 1 {
		return nil, &errors.NotPositiveDataLengthError{"Initialization vector"}
	}

	var unsignedIV []byte
	var IV []int8
	for ok := true; ok; ok = isDifferenceOrthogonal(key.Data, IV, rotatedData) {
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

// DeriveKey derives key from given big integer key, salt. DerivedKey has the same length as data, so it is dataLength.
// It returns derived key as struct NSEKey and err, err != nil if and only if given key is not positive or hkdf returns an error.
func DeriveKey(key *big.Int, salt []byte, dataLength int) (derivedKey *NSEKey, err error) {
	if key.Cmp(big.NewInt(0)) <= 0 {
		return derivedKey, &errors.NotPositiveIntegerKeyError{key}
	}
	var bigKeyWithExcludedLength big.Int
	bigKeyWithExcludedLength.Mod(key, big.NewInt(int64(dataLength<<3)))
	keyWithExcludedLength := bigKeyWithExcludedLength.Uint64()
	derivedKey = &NSEKey{
		BitsToRotate:  byte(keyWithExcludedLength & 7),
		BytesToRotate: int(keyWithExcludedLength >> 3),
		Data:          make([]int8, dataLength)}
	unsignedDerivedKey := make([]byte, dataLength)
	keyCopy := new(big.Int)
	keyCopy.SetBytes(key.Bytes())
	for ok := true; ok; ok = isZeroVector(derivedKey.Data) {
		_, err := io.ReadFull(hkdf.New(sha512.New, keyCopy.Bytes(), salt, nil), unsignedDerivedKey)
		if err != nil {
			return derivedKey, err
		}
		for i, v := range unsignedDerivedKey {
			derivedKey.Data[i] = bits.AsSigned(v)
		}
		keyCopy.Add(keyCopy, bigOne)
	}

	return
}
