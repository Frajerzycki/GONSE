package test

import (
	"bytes"
	"encoding/binary"
	"github.com/ikcilrep/gonse/pkg/nse"
	"math/rand"
	"testing"
	"time"
)

func Test_nse_BytesToInt8s(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 512; i++ {
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

func randomInt64s(length int) (data []int64) {
	data = make([]int64, length)
	for index := 0; index < length; index++ {
		buffer := make([]byte, 8)
		rand.Read(buffer)
		data[index], _ = binary.Varint(buffer)

	}
	return
}

func equalsInt64s(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for index, value := range a {
		if value != b[index] {
			return false
		}
	}
	return true
}

func Test_nse_Int64sToBytes(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 512; i++ {
		data := randomInt64s(i)
		convertedData := nse.Int64sToBytes(data)
		unconvertedData, err := nse.BytesToInt64s(convertedData)
		if err != nil {
			t.Error(err)
		}
		if !equalsInt64s(data, unconvertedData) {
			t.Errorf("%v as unsigned is %v, but %v as signed is %v which is not the same", data, convertedData, convertedData, unconvertedData)
		}
	}

}
