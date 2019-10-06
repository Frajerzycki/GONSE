package nse

import (
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

type IVZeroLengthError struct{}

type BytesDivisionError struct {
	dataLength int
}

func (err DifferentIVLengthError) Error() string {
	return fmt.Sprintf("Intialization vector is different length than data: %v != %v.", err.dataLength, err.IVLength)
}

func (err NotPositiveIntegerKeyError) Error() string {
	return fmt.Sprintf("Key has to be positive integer, but is %v.", err.key)
}

func (err NilArgumentError) Error() string {
	return fmt.Sprintf("%v mustn't be nil nor empty slice", err.argumentName)
}

func (err IVZeroLengthError) Error() string {
	return "Initialization vector length has to be positive."
}

func (err BytesDivisionError) Error() string {
	return fmt.Sprintf("There are %v bytes that cannot be divided equally into eight byte parts.", err.dataLength)
}
