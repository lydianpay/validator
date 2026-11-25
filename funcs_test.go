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

func TestValidateStructHostname(t *testing.T) {
	type hostnames struct {
		Hostname        string  `validator:"hostname"`
		HostnamePointer *string `validator:"hostname"`
	}
	tooLong := strings.Repeat(strings.Repeat("a", 62)+".", 4) + "com"
	testCases := map[string]struct {
		hostname string
		isValid  bool
	}{
		"a subdomain":                {"checkout.example.com", true},
		"mixed case and hyphens":     {"Crypto.Checkout-Sandbox.Maya.PH", true},
		"an empty value":             {"", true},
		"a single label":             {"localhost", false},
		"a URL":                      {"https://checkout.example.com", false},
		"a path":                     {"checkout.example.com/session", false},
		"a port":                     {"checkout.example.com:443", false},
		"surrounding whitespace":     {" checkout.example.com ", false},
		"a trailing dot":             {"checkout.example.com.", false},
		"a leading hyphen":           {"-checkout.example.com", false},
		"an underscore":              {"check_out.example.com", false},
		"an IP address":              {"10.0.0.1", false},
		"a label over 63 characters": {strings.Repeat("a", 64) + ".example.com", false},
		"over 253 characters":        {tooLong, false},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			hostname := testCase.hostname
			err := ValidateStruct(hostnames{Hostname: hostname, HostnamePointer: &hostname})
			if testCase.isValid && err != nil {
				t.Errorf("expected %q to pass, got %v", hostname, err)
			}
			if !testCase.isValid && (err == nil || strings.Count(err.Error(), "is not a valid hostname") != 2) {
				t.Errorf("expected %q to fail on both fields, got %v", hostname, err)
			}
		})
	}

	t.Run("an unset pointer", func(t *testing.T) {
		if err := ValidateStruct(hostnames{}); err != nil {
			t.Errorf("expected an unset pointer to pass, got %v", err)
		}
	})

	t.Run("paired with required", func(t *testing.T) {
		type requiredHostname struct {
			Hostname string `validator:"required;hostname"`
		}
		if err := ValidateStruct(requiredHostname{}); err == nil || !strings.Contains(err.Error(), "is required") {
			t.Errorf("expected an empty required hostname to fail as required, got %v", err)
		}
		if err := ValidateStruct(requiredHostname{Hostname: "localhost"}); err == nil || !strings.Contains(err.Error(), "is not a valid hostname") {
			t.Errorf("expected a single label to fail, got %v", err)
		}
	})

	t.Run("a non-string field", func(t *testing.T) {
		type numericHostname struct {
			Hostname int `validator:"hostname"`
		}
		if err := ValidateStruct(numericHostname{Hostname: 1}); err == nil {
			t.Error("expected a non-string field to fail")
		}
	})
}
