package validator

import (
	"testing"
)

type TS1 struct {
	Username             string `json:"username" validator:"required"`
	MultipleRequirements int    `json:"multi" validator:"required;gt=0"`
}

func TestDescribeBase(t *testing.T) {
	if err := DescribeBase(&TS1{Username: "a"}); err != nil {
		t.Fatal("describe happy path")
	}
	if err := DescribeBase("a"); err == nil {
		t.Fatal("describe non-struct")
	}
}

func TestIsStruct(t *testing.T) {
	if !isStruct(TS1{}) {
		t.Fatal("isStruct happy path")
	}
	if isStruct("") {
		t.Fatal("isStruct non-struct")
	}
}

func TestCompareFloats(t *testing.T) {
	if !compareFloats(1, 2, VLessThan) {
		t.Fatal("happy lt f64")
	}
	if !compareFloats(1, 2, VLessThanOrEqual) {
		t.Fatal("happy lte f64")
	}
	if !compareFloats(1, 1, VEqualTo) {
		t.Fatal("happy eq f64")
	}
	if !compareFloats(2, 1, VGreaterThan) {
		t.Fatal("happy gt f64")
	}
	if !compareFloats(2, 1, VGreaterThanOrEqual) {
		t.Fatal("happy gte f64")
	}
	if compareFloats(1, 1, "are these equal to each other?") {
		t.Fatal("nonsense comparison f64")
	}
}
func TestCompareInts(t *testing.T) {
	if !compareInts(1, 2, VLessThan) {
		t.Fatal("happy lt int")
	}
	if !compareInts(1, 2, VLessThanOrEqual) {
		t.Fatal("happy lte int")
	}
	if !compareInts(1, 1, VEqualTo) {
		t.Fatal("happy eq int")
	}
	if !compareInts(2, 1, VGreaterThan) {
		t.Fatal("happy gt int")
	}
	if !compareInts(2, 1, VGreaterThanOrEqual) {
		t.Fatal("happy gte int")
	}
	if compareInts(1, 1, "are these equal to each other?") {
		t.Fatal("nonsense comparison int")
	}
}
