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
	for i := 1; i <= 512; i++ {
		data := make([]byte, i)
		key := make([]byte, i)
		rand.Read(data)
		rand.Read(key)

		// Just cast to signed in simple way.
		rotatedData := bits.RightRotate(data, 0, 0)
		derivedKey := bits.RightRotate(key, 0, 0)

		IV, err := nse.GenerateIV(i, rotatedData, derivedKey)
		if err != nil {
			t.Error(err)
		}

		if len(IV) != i {
			t.Errorf("%v is not the same length as %v", len(IV), i)
		}
	}
}
