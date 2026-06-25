package validator

import (
	"fmt"
	"net/mail"
	"reflect"
	"slices"
	"strings"
)

func ValidateStruct(input any) error {
	if !isStruct(input) {
		return fmt.Errorf("not a struct")
	}
	var expectedFields []reflect.StructField
	if reflect.TypeOf(input).Kind() == reflect.Pointer {
		expectedFields = reflect.VisibleFields(reflect.TypeOf(input).Elem())
	} else {
		expectedFields = reflect.VisibleFields(reflect.TypeOf(input))
	}
	dereffedInput := reflect.ValueOf(input)
	if dereffedInput.Kind() == reflect.Pointer {
		dereffedInput = dereffedInput.Elem()
	}
	problems := []string{}
	for _, baselineField := range expectedFields {
		prefix, param := getTagPartsFromField(baselineField)
		problems = append(problems, validateField(dereffedInput.FieldByName(baselineField.Name), dereffedInput, baselineField, prefix, param)...)
	}
	if len(problems) > 0 {
		return fmt.Errorf("validation failed : %s", strings.Join(problems, ", "))
	}
	return nil
}

func validateField(fieldVal, dereffedInput reflect.Value, baselineField reflect.StructField, prefix, param string) []string {
	problems := []string{}
	switch prefix {
	case VRequired:
		if fieldVal.IsZero() {
			problems = append(problems, fmt.Sprintf("field '%s' is required", baselineField.Name))
		}
	case VExactStr:
		if good, expected := vExactStr(fieldVal, param); !good {
			problems = append(problems, fmt.Sprintf("field '%s' must have the exact value '%s'", baselineField.Name, expected))
		}
	case VGreaterThan, VLessThan, VGreaterThanOrEqual, VLessThanOrEqual, VEqualTo:
		if ok, expected := vCompareNumber(baselineField, fieldVal, param, prefix); !ok {
			problems = append(problems, fmt.Sprintf("field '%s' failed the comparison '%s' vs the value '%s'", baselineField.Name, prefix, expected))
		}
	case VOneOf:
		if ok, allowed := vOneOf(fieldVal, param); !ok {
			problems = append(problems, fmt.Sprintf("field '%s' does not contain a value on the list of allowed values (%s)", baselineField.Name, allowed))
		}
	case VEmail:
		if fieldVal.Kind() != reflect.String {
			problems = append(problems, fmt.Sprintf("unable to validate email addresses on non-string field %s", baselineField.Name))
		} else if _, err := mail.ParseAddress(fieldVal.String()); err != nil {
			problems = append(problems, fmt.Sprintf("invalid email address for field %s", baselineField.Name))
		}
	case vMultipleTags:
		multi := getMultipleTagParts(baselineField)
		multiProbs := []string{}
		for multiPrefix, multiParam := range multi {
			multiProbs = append(multiProbs, validateField(fieldVal, dereffedInput, baselineField, multiPrefix, multiParam)...)
		}
		problems = append(problems, multiProbs...)
	case vNoTag:
		return nil
	case vInvalid:
		problems = append(problems, fmt.Sprintf("field '%s' has an incorrectly defined validator tag", baselineField.Name))
	default:
		problems = append(problems, fmt.Sprintf("don't have a handler for %s yet", prefix))
	}
	return problems
}

func vExactStr(fieldVal reflect.Value, param string) (bool, string) {
	return fieldVal.String() == param, param
}

func vCompareNumber(field reflect.StructField, fieldVal reflect.Value, param, comparisonOperator string) (bool, string) {
	switch field.Type {
	case reflect.TypeFor[float64]():
		actual, expected, ok := prepFloats(fieldVal, param)
		if !ok {
			return false, fmt.Sprint(expected)
		}
		return compareFloats(actual, expected, comparisonOperator), fmt.Sprint(expected)
	case reflect.TypeFor[int]():
		actual, expected, ok := prepInts(fieldVal, param)
		if !ok {
			return false, fmt.Sprint(expected)
		}
		return compareInts(actual, expected, comparisonOperator), fmt.Sprint(expected)
	default:
		return false, "don't currently do greaterthan for this type"
	}
}

func vOneOf(fieldVal reflect.Value, param string) (bool, string) {
	allowed := strings.Split(param, MultiChoiceSeparator)
	if fieldVal.Kind() == reflect.Pointer {
		if fieldVal.IsNil() {
			return false, param
		}
		fieldVal = fieldVal.Elem()
	}
	return slices.Contains(allowed, fieldVal.String()), param
}
