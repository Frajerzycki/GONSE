package nse

import (
	"crypto/rand"
	"crypto/sha512"
	"io"
	"math/big"

	"github.com/ikcilrep/gonse/internal/bits"
	"github.com/ikcilrep/gonse/internal/errors"
	"golang.org/x/crypto/hkdf"
)

type NSEKey struct {
	Data          []uint16
	BitsToRotate  byte
	BytesToRotate int
}

var primes = [256]uint16{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
	31, 37, 41, 43, 47, 53, 59, 61, 67, 71,
	73, 79, 83, 89, 97, 101, 103, 107, 109, 113,
	127, 131, 137, 139, 149, 151, 157, 163, 167, 173,
	179, 181, 191, 193, 197, 199, 211, 223, 227, 229,
	233, 239, 241, 251, 257, 263, 269, 271, 277, 281,
	283, 293, 307, 311, 313, 317, 331, 337, 347, 349,
	353, 359, 367, 373, 379, 383, 389, 397, 401, 409,
	419, 421, 431, 433, 439, 443, 449, 457, 461, 463,
	467, 479, 487, 491, 499, 503, 509, 521, 523, 541,
	547, 557, 563, 569, 571, 577, 587, 593, 599, 601,
	607, 613, 617, 619, 631, 641, 643, 647, 653, 659,
	661, 673, 677, 683, 691, 701, 709, 719, 727, 733,
	739, 743, 751, 757, 761, 769, 773, 787, 797, 809,
	811, 821, 823, 827, 829, 839, 853, 857, 859, 863,
	877, 881, 883, 887, 907, 911, 919, 929, 937, 941,
	947, 953, 967, 971, 977, 983, 991, 997, 1009, 1013,
	1019, 1021, 1031, 1033, 1039, 1049, 1051, 1061, 1063, 1069,
	1087, 1091, 1093, 1097, 1103, 1109, 1117, 1123, 1129, 1151,
	1153, 1163, 1171, 1181, 1187, 1193, 1201, 1213, 1217, 1223,
	1229, 1231, 1237, 1249, 1259, 1277, 1279, 1283, 1289, 1291,
	1297, 1301, 1303, 1307, 1319, 1321, 1327, 1361, 1367, 1373,
	1381, 1399, 1409, 1423, 1427, 1429, 1433, 1439, 1447, 1451,
	1453, 1459, 1471, 1481, 1483, 1487, 1489, 1493, 1499, 1511,
	1523, 1531, 1543, 1549, 1553, 1559, 1567, 1571, 1579, 1583,
	1597, 1601, 1607, 1609, 1613, 1619,
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

func isDifferenceOrthogonal(derivedKey []uint16, IV, rotatedData []int8) bool {
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
		Data:          make([]uint16, dataLength)}
	unsignedDerivedKey := make([]byte, dataLength)
	keyCopy := new(big.Int)
	keyCopy.SetBytes(key.Bytes())
	_, err = io.ReadFull(hkdf.New(sha512.New, keyCopy.Bytes(), salt, nil), unsignedDerivedKey)
	if err != nil {
		return derivedKey, err
	}
	for i, v := range unsignedDerivedKey {
		derivedKey.Data[i] = primes[v]
	}
	keyCopy.Add(keyCopy, bigOne)

	return
}
