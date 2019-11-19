package test

import (
	"github.com/ikcilrep/gonse/pkg/nse"
	"testing"
)

func TestGeneratingIV(t *testing.T) {
	for i := 1; i <= 512; i++ {
		IV, err := nse.GenerateIV(i)
		if err != nil {
			t.Error(err)
		}

		if len(IV) != i {
			t.Errorf("%v is not the same length as %v", len(IV), i)
		}
	}
}
