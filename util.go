package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func DescribeBase(input any) error {
	if !isStruct(input) {
		return fmt.Errorf("not a struct")
	}
	dereffedInput := reflect.ValueOf(input)
	if dereffedInput.Kind() == reflect.Pointer {
		dereffedInput = dereffedInput.Elem()
	}
	expectedFields := reflect.VisibleFields(dereffedInput.Type())
	outputFmt := "%20s |%20s |%20s\n"
	fmt.Printf(outputFmt, "Field Name", "Field Type", "Validation")
	for _, field := range expectedFields {
		fmt.Printf(outputFmt, field.Name, field.Type.String(), field.Tag.Get(TagVBase))
	}
	fmt.Println("---")
	return nil
}

func isStruct(input any) bool {
	dri := reflect.ValueOf(input)
	if dri.Kind() == reflect.Pointer {
		dri = dri.Elem()
	}
	if dri.Kind() == reflect.Struct {
		return true
	}
	return false
}

func getMultipleTagParts(field reflect.StructField) map[string]string {
	entireTag := field.Tag.Get(TagVBase)
	if !strings.Contains(entireTag, MultiValidatorSeparator) {
		return nil
	}
	multiParts := strings.Split(entireTag, MultiValidatorSeparator)
	output := map[string]string{}
	for _, multiPart := range multiParts {
		prefix, param := getTagPartsFromEntireTag(multiPart)
		output[prefix] = param
	}
	return output
}

func getTagPartsFromField(field reflect.StructField) (string, string) {
	entireTag := field.Tag.Get(TagVBase)
	if entireTag == "" {
		return vNoTag, ""
	}
	if strings.Contains(entireTag, MultiValidatorSeparator) {
		return vMultipleTags, ""
	}
	return getTagPartsFromEntireTag(entireTag)
}

func getTagPartsFromEntireTag(entireTag string) (string, string) {
	if entireTag == VRequired || entireTag == VEmail {
		return entireTag, ""
	}
	parts := strings.Split(entireTag, TagSeparator)
	if len(parts) == 2 {
		return parts[0] + TagSeparator, parts[1]
	}
	return vInvalid, ""
}

func prepFloats(fieldVal reflect.Value, subtag string) (float64, float64, bool) {
	if !fieldVal.CanFloat() {
		return 0, 0, false
	}
	expectedVal, err := strconv.ParseFloat(subtag, 64)
	if err != nil {
		return 0, 0, false
	}
	return fieldVal.Float(), expectedVal, true
}

func prepInts(fieldVal reflect.Value, subtag string) (int, int, bool) {
	if !fieldVal.CanInt() {
		return 0, 0, false
	}
	expectedVal, err := strconv.Atoi(subtag)
	if err != nil {
		return 0, 0, false
	}
	return int(fieldVal.Int()), expectedVal, true
}

func compareFloats(actual, expected float64, comparisonTag string) bool {
	switch comparisonTag {
	case VGreaterThan:
		return actual > expected
	case VGreaterThanOrEqual:
		return actual >= expected
	case VLessThan:
		return actual < expected
	case VLessThanOrEqual:
		return actual <= expected
	case VEqualTo:
		return actual == expected
	}
	return false
}

func compareInts(actual, expected int, comparisonTag string) bool {
	switch comparisonTag {
	case VGreaterThan:
		return actual > expected
	case VGreaterThanOrEqual:
		return actual >= expected
	case VLessThan:
		return actual < expected
	case VLessThanOrEqual:
		return actual <= expected
	case VEqualTo:
		return actual == expected
	}
	return false
}
