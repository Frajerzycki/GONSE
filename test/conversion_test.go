package test

import (
	"bytes"
	"github.com/ikcilrep/gonse/pkg/nse"
	"math/rand"
	"testing"
	"time"
)

func TestUnsigningBytes(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 1; i < 512; i++ {
		data := make([]byte, i)
		rand.Read(data)
		convertedData := nse.BytesToInt8s(data)
		unconvertedData := nse.Int8sToBytes(convertedData)
		if len(convertedData) != len(data) {
			t.Errorf("%v has different length than %v", data, convertedData)
		}
		if !bytes.Equal(data, unconvertedData) {
			t.Errorf("%v as unsigned is %v, but %v as signed is %v which is not the same", data, convertedData, convertedData, unconvertedData)
		}
	}
}
