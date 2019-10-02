package nse

import "encoding/binary"


func Int64sToBytes(data []int64) ([]byte, error) {
  if data == nil {
    return nil, NilArgumentError{"Data"}
  }

  dataLength := len(data)
  resultLength := dataLength << 3
  result := make([]byte, resultLength)
  for index := 0; index < resultLength; index += 8 {
    binary.PutVarint(result[index:], data[index >> 3])
  }
  return result, nil
}


func BytesToInt64s(data []byte) ([]int64, error) {
  if data == nil {
    return nil, NilArgumentError{"Data"}
  }

  dataLength := len(data)
  if dataLength & 7 > 0 {
    return nil, BytesDivisionError{dataLength}
  }

  resultLength := dataLength >> 3
  result := make([]int64, resultLength)
  for index := 0; index < dataLength; index += 8 {
    result[index >> 3], _ = binary.Varint(data[index:index+8])
  }
  return result, nil
}
