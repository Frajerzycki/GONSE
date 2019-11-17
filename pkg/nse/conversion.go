package nse

import (
	"encoding/binary"
	"github.com/ikcilrep/gonse/internal/errors"
)

func Int64sToBytes(data []int64) []byte {
	dataLength := len(data)
	resultLength := dataLength * 9
	result := make([]byte, resultLength)
	resultIndex := 0
	for dataIndex := 0; dataIndex < dataLength; dataIndex++ {
		buffer := make([]byte, 8)
		binary.PutVarint(buffer, data[dataIndex])
		lastNonZeroIndex := 7
		for ; buffer[lastNonZeroIndex] == 0; lastNonZeroIndex-- {
		}
		result[resultIndex] = byte(lastNonZeroIndex + 1)
		resultIndex++
		for index := 0; index <= lastNonZeroIndex; index++ {
			result[resultIndex] = buffer[index]
			resultIndex++
		}
	}
	return result[:resultIndex]
}

func BytesToInt64s(data []byte) ([]int64, error) {
	dataLength := len(data)

	resultLength := dataLength
	result := make([]int64, resultLength)
	resultIndex := 0
	for dataIndex := 0; dataIndex < dataLength; resultIndex++ {
		newDataIndex := dataIndex + int(data[dataIndex]) + 1
		if newDataIndex > dataLength {
			return nil, errors.WrongDataFormatError
		}

		result[resultIndex], _ = binary.Varint(data[dataIndex+1 : newDataIndex])
		dataIndex = newDataIndex
	}
	return result[:resultIndex], nil
}

func Int8sToBytes(data []int8) []byte {
	dataLength := len(data)
	result := make([]byte, dataLength)
	for index, value := range data {
		result[index] = asUnsigned(value)
	}
	return result
}

func BytesToInt8s(data []byte) []int8 {
	dataLength := len(data)
	result := make([]int8, dataLength)
	for index, value := range data {
		result[index] = asSigned(value)
	}
	return result
}
