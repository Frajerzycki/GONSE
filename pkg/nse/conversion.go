package nse

import (
	"encoding/binary"
	"io"

	"github.com/ikcilrep/gonse/internal/bits"
	"github.com/ikcilrep/gonse/internal/errors"
)

// Int64ToBytes converts integer into byte array. It ignores padding, result is as short as possible. First byte is length of the rest.
func Int64ToBytes(integer int64) []byte {
	bytes := make([]byte, 9)
	binary.PutVarint(bytes[1:], integer)
	lastNonZeroIndex := 8
	for ; lastNonZeroIndex > 1 && bytes[lastNonZeroIndex] == 0; lastNonZeroIndex-- {
	}
	bytes[0] = byte(lastNonZeroIndex)
	return bytes[:lastNonZeroIndex+1]
}

// BytesToInt64 converts byte array into int64. It returns converted integer, bytes read and err.
// err != nil if and only if binary.Varint returns an error.
func BytesToInt64(data []byte) (int64, int, error) {
	length := int(data[0]) + 1
	result, bytesRead := binary.Varint(data[1:length])
	if bytesRead <= 0 {
		return int64(0), 0, errors.WrongDataFormatError
	}
	return result, length, nil
}

// BytesToInt64FromReader converts few first bytes into int64. It returns converted integer, bytes read and err.
// err != nil if and only if binary.Varint returns an error or given reader doesn't have that many bytes to read.
func BytesToInt64FromReader(reader io.Reader) (int64, int, error) {
	lengthByte := make([]byte, 1)
	_, err := io.ReadFull(reader, lengthByte)
	if err != nil {
		return int64(0), 0, err
	}
	length := int(lengthByte[0])
	resultBytes := make([]byte, length)
	_, err = io.ReadFull(reader, resultBytes)
	if err != nil {
		return int64(0), 0, err
	}

	result, bytesRead := binary.Varint(resultBytes)
	if bytesRead <= 0 {
		return int64(0), 0, errors.WrongDataFormatError
	}
	return result, length + 1, nil
}

// Int64sToBytes converts []int64 into []byte.
// For each int64 in the slice there is one byte indicating how many bytes to read next and those bytes.
func Int64sToBytes(data []int64) []byte {
	dataLength := len(data)
	resultLength := dataLength * 9
	result := make([]byte, resultLength)
	resultIndex := 0
	for dataIndex := 0; dataIndex < dataLength; dataIndex++ {
		integerBytes := Int64ToBytes(data[dataIndex])
		copy(result[resultIndex:], integerBytes)
		resultIndex += len(integerBytes)
	}
	return result[:resultIndex]
}

// BytesToInt64s converts result of Int64sToBytes back into []int64.
// It returns errors.WrongDataFormatError as an error when data doesn't appear to be a result of Int64sToBytes.
func BytesToInt64s(data []byte) ([]int64, error) {
	dataLength := len(data)

	resultLength := dataLength
	result := make([]int64, resultLength)
	resultIndex := 0
	for dataIndex := 0; dataIndex < dataLength; resultIndex++ {
		var bytesRead int
		var err error
		result[resultIndex], bytesRead, err = BytesToInt64(data[dataIndex:])
		if err != nil {
			return nil, err
		}
		dataIndex += bytesRead
	}
	return result[:resultIndex], nil
}

// Int8sToBytes converts []int8 into []byte.
// Every int8 in the slice is treated like it would be unsigned.
func Int8sToBytes(data []int8) []byte {
	dataLength := len(data)
	result := make([]byte, dataLength)
	for index, value := range data {
		result[index] = bits.AsUnsigned(value)
	}
	return result
}

// BytesToInt8s converts []byte into []int8.
// Every byte in the slice is treated like it would be signed.
func BytesToInt8s(data []byte) []int8 {
	dataLength := len(data)
	result := make([]int8, dataLength)
	for index, value := range data {
		result[index] = bits.AsSigned(value)
	}
	return result
}
