package main

import (
	"testing"
)

func TestWhenInvokeMaxWithSecondParameterBiggerThanFirstOneThenReturnTheBiggestOne(t *testing.T) {
	max := Max(5, 10)
	if max != 10 {
		t.Fatalf("Max is not retuning the max value 10, returned %v", max)
	}
}

func TestWhenInvokeMaxWithFirstParameterBiggerThanSecondOneThenReturnTheBiggestOne(t *testing.T) {
	max := Max(10, 5)
	if max != 10 {
		t.Fatalf("Max is not retuning the max value 10, returned %v", max)
	}
}

func TestWhenInvokeMaxWithNegativeValuesThenReturnTheNearestToZero(t *testing.T) {
	max := Max(-10, -5)
	if max != -5 {
		t.Fatalf("Max is not retuning the max value -5, returned %v", max)
	}
}
