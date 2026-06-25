package validator

import (
	"strings"
	"testing"
)

type TS2 struct {
	Username             string  `json:"username" validator:"required"`
	MultipleRequirements int     `json:"multi" validator:"required;gt=0"`
	Specific             string  `json:"specific" validator:"exact=mustbe"`
	PetChoice            string  `json:"petChoice" validator:"oneof=cat,dog,possum,raccoon"`
	Email                string  `json:"email" validator:"email"`
	NoTagValue           float64 `json:"noTagValue"`
	Price                float64 `json:"price" validator:"gte=.01"`
}
type TS3 struct {
	Field1 string `json:"field1" validator:"oneof"`
	Field2 string `json:"field2" validator:"madeUp=Elvis"`
}

func TestValidateStruct(t *testing.T) {
	goodTS2 := &TS2{
		Username:             "fake username",
		MultipleRequirements: 1,
		Specific:             "mustbe",
		PetChoice:            "possum",
		Email:                "test@test.test",
		NoTagValue:           1.23,
		Price:                .02,
	}
	if err := ValidateStruct(goodTS2); err != nil {
		t.Fatal("validatestruct happy path")
	}
	err := ValidateStruct(TS2{
		MultipleRequirements: -1,
		Price:                -.01,
	},
	)
	if err == nil {
		t.Fatal("validatestruct several issues")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Fatal("missing required field")
	}
	if !strings.Contains(err.Error(), "must have the exact value") {
		t.Fatal("missing exact field")
	}
	if !strings.Contains(err.Error(), "failed the comparison") {
		t.Fatal("missing numeric comparison")
	}
	if !strings.Contains(err.Error(), "invalid email address") {
		t.Fatal("missing email")
	}
	err = ValidateStruct(TS3{})
	if err == nil || !strings.Contains(err.Error(), "has an incorrectly defined validator tag") || !strings.Contains(err.Error(), "don't have a handler for") {
		t.Fatal("invald tag")
	}
	if err = ValidateStruct("hello"); err == nil {
		t.Fatal("attempted to validate non-struct")
	}
}

type TSEdge struct {
	BigNum         int64   `validator:"gt=1"`          // unsupported numeric type -> vCompareNumber default
	BadIntCmp      int     `validator:"gt=abc"`        // prepInts parse failure
	BadFloatCmp    float64 `validator:"gt=abc"`        // prepFloats parse failure
	EmailNonString int     `validator:"email"`         // email on non-string field
	PetPtr         *string `validator:"oneof=cat,dog"` // oneof pointer handling
}

func TestValidateStructEdgeCases(t *testing.T) {
	err := ValidateStruct(TSEdge{}) // nil PetPtr exercises the oneof nil-pointer branch
	if err == nil {
		t.Fatal("expected validation errors for edge struct")
	}
	for _, want := range []string{
		"failed the comparison",
		"unable to validate email addresses on non-string field",
		"does not contain a value on the list of allowed values",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("expected error to contain %q, got %q", want, err.Error())
		}
	}

	pet := "cat" // non-nil pointer exercises the oneof deref path
	if err := ValidateStruct(TSEdge{PetPtr: &pet}); err == nil {
		t.Error("expected errors from the other fields even with a valid PetPtr")
	}
}
