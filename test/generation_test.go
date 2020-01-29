package test

import (
	"github.com/ikcilrep/gonse/internal/bits"
	"github.com/ikcilrep/gonse/pkg/nse"
	"math/rand"
	"testing"
	"time"
)

func Test_nse_GenerateIV(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for dataLength := 1; dataLength <= 512; dataLength++ {
		data := make([]byte, dataLength)
		key := make([]byte, dataLength)
		rand.Read(data)
		rand.Read(key)

		// Just cast to signed in simple way.
		rotatedData := bits.RightRotate(data, 0, 0)
		derivedKey := &nse.NSEKey{Data: bits.RightRotate(key, 0, 0)}

		IV, err := nse.GenerateIV(dataLength, rotatedData, derivedKey)
		if err != nil {
			t.Error(err)
		}

		if len(IV) != dataLength {
			t.Errorf("%v is not the same length as %v", len(IV), dataLength)
		}
	}
}
