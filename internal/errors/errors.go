package errors

import (
	"errors"
	"fmt"
	"math/big"
)

type DifferentIVLengthError struct {
	IVLength   int
	DataLength int
}

type NotPositiveIntegerKeyError struct {
	Key *big.Int
}

type NotPositiveDataLengthError struct {
	DataName string
}

var WrongDataFormatError error = errors.New("Wrong data format.")

func (err *DifferentIVLengthError) Error() string {
	return fmt.Sprintf("Intialization vector is different length than data: %v != %v.", err.DataLength, err.IVLength)
}

func (err *NotPositiveIntegerKeyError) Error() string {
	return fmt.Sprintf("Key has to be positive integer, but is %v.", err.Key)
}

func (err *NotPositiveDataLengthError) Error() string {
	return fmt.Sprintf("%v has to has positive length.", err.DataName)
}
