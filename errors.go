package nse

import (
	"errors"
	"fmt"
	"math/big"
)

type DifferentIVLengthError struct {
	IVLength   int
	dataLength int
}

type NotPositiveIntegerKeyError struct {
	key *big.Int
}

type NilArgumentError struct {
	argumentName string
}

var IVZeroLengthError error = errors.New("Initialization vector length has to be positive.")
var WrongDataFormatError error = errors.New("Wrong data format.")

func (err DifferentIVLengthError) Error() string {
	return fmt.Sprintf("Intialization vector is different length than data: %v != %v.", err.dataLength, err.IVLength)
}

func (err NotPositiveIntegerKeyError) Error() string {
	return fmt.Sprintf("Key has to be positive integer, but is %v.", err.key)
}

func (err NilArgumentError) Error() string {
	return fmt.Sprintf("%v mustn't be nil nor empty slice", err.argumentName)
}
