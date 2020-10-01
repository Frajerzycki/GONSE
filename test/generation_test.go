package test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ikcilrep/gonse/internal/bits"
	"github.com/ikcilrep/gonse/pkg/nse"
)

func Test_nse_GenerateIV(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for dataLength := 1; dataLength <= 512; dataLength++ {
		data := make([]byte, dataLength)
		keyBytes := make([]byte, 2*dataLength)
		rand.Read(data)
		rand.Read(keyBytes)

		key := make([]uint16, dataLength)
		for i := 0; i < 2*dataLength; i += 2 {
			key[i/2] = uint16(keyBytes[i])<<8 + uint16(keyBytes[i+1])
		}

		// Just cast to signed in simple way.
		rotatedData := bits.RightRotate(data, 0, 0)
		derivedKey := &nse.NSEKey{Data: key}

		IV, err := nse.GenerateIV(dataLength, rotatedData, derivedKey)
		if err != nil {
			t.Error(err)
		}

		if len(IV) != dataLength {
			t.Errorf("%v is not the same length as %v", len(IV), dataLength)
		}
	}
}
