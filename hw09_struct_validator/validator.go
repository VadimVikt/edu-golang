package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	TagLength  = "len"
	TagMin     = "min"
	TagMax     = "max"
	TagIn      = "in"
	TagRequire = "required"
	TagRegex   = "regexp"
)

var (
	ErrRequiredField         = errors.New("field required")
	ErrInvalidLength         = errors.New("invalid length")
	ErrValueOutOfRange       = errors.New("value out of range")
	ErrInvalidFormat         = errors.New("invalid format")
	ErrValueNotInList        = errors.New("value not in allowed list")
	ErrLessThanMin           = errors.New("field is less than min")
	ErrGraterThanMax         = errors.New("field is greater than max")
	ErrExpectedStruct        = errors.New("expected a struct")
	ErrNoValidTagFormat      = errors.New("tag value must be in digital format")
	ErrNoValidRegularExpr    = errors.New("required field value must be in digital format")
	ErrNoValidTagNotProvider = errors.New("tag is not provided in this structure")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "no errors"
	}
	errorMessages := make([]string, 0, len(v))
	for _, err := range v {
		errorMessages = append(errorMessages, err.Err.Error())
	}
	return strings.Join(errorMessages, "\n ")
}

func (v ValidationErrors) Is(target error) bool {
	for _, err := range v {
		if errors.Is(err.Err, target) {
			return true
		}
	}
	return false
}

func Validate(v interface{}) error {
	var result ValidationErrors
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return ErrExpectedStruct
	}
	val := reflect.ValueOf(v)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagString := field.Tag.Get("validate")
		if tagString == "" {
			continue
		}
		tags := parseTags(tagString)
		res, err := ValidatorFunc(field, val.Field(i), tags)
		if err != nil {
			return err // Прерывание валидации
		}
		if res.Err == nil {
			continue
		}
		result = append(result, res)
	}
	return result
}
func ValidatorFunc(
	fieldType reflect.StructField,
	fieldValue reflect.Value,
	tags map[string]string) (ValidationError, error) {
	switch fieldType.Type.Kind() { //nolint
	case reflect.String:
		validErr, err := validateString(fieldType.Name, fieldValue.String(), tags)
		if err != nil {
			return ValidationError{}, errors.New("tag not provided, validation break")
		}
		if validErr.Err == nil {
			return ValidationError{}, nil
		}
		return validErr, nil
	case reflect.Int:
		validErr, err := validateInt(fieldType.Name, int(fieldValue.Int()), tags)
		if err != nil {
			return ValidationError{}, errors.New("tag not provided, validation break")
		}
		if validErr.Err == nil {
			return ValidationError{}, nil
		}
		return validErr, nil
	case reflect.Slice:
		switch kind := fieldType.Type.Elem().Kind(); kind { //nolint
		case reflect.String:
			validErr, err := validateSliceString(fieldType.Name, fieldValue, tags)
			if err != nil {
				return ValidationError{}, errors.New("tag not provided, validation break")
			}
			if validErr.Err == nil {
				return ValidationError{}, nil
			}
			return validErr, nil
		case reflect.Int:
			validErr, err := validateSliceInt(fieldType.Name, fieldValue, tags)
			if err != nil {
				return ValidationError{}, errors.New("tag not provided, validation break")
			}
			if validErr.Err == nil {
				return ValidationError{}, nil
			}
			return validErr, nil
		}
	}
	return ValidationError{}, nil
}

func validateSliceString(fieldName string, field reflect.Value, tags map[string]string) (ValidationError, error) {
	for k, v := range tags {
		switch k {
		case TagLength:
			lenField, err := strconv.Atoi(v)
			if err != nil {
				return ValidationError{}, ErrNoValidTagFormat
			}
			tmp := field.Interface().([]string)
			for _, str := range tmp {
				if len(str) != lenField {
					return ValidationError{fieldName, ErrInvalidLength}, nil
				}
			}
		default:
			return ValidationError{}, ErrNoValidTagNotProvider
		}
	}
	return ValidationError{fieldName, nil}, nil
}

func validateSliceInt(fieldName string, field reflect.Value, tags map[string]string) (ValidationError, error) {
	for k, v := range tags {
		switch k {
		case TagLength:
			lenField, err := strconv.Atoi(v)
			if err != nil {
				return ValidationError{}, ErrNoValidTagFormat
			}
			if len(field.Interface().([]int)) != lenField {
				return ValidationError{fieldName, ErrInvalidLength}, nil
			}
		default:
			return ValidationError{}, ErrNoValidTagNotProvider
		}
	}
	return ValidationError{fieldName, nil}, nil
}

func validateInt(fieldName string, field int, tags map[string]string) (ValidationError, error) {
	for k, v := range tags {
		switch k {
		case TagMax:
			val, err := strconv.Atoi(v)
			if err != nil {
				return ValidationError{}, ErrNoValidTagFormat
			}
			if field > val {
				return ValidationError{fieldName, ErrGraterThanMax}, nil
			}
		case TagMin:
			val, err := strconv.Atoi(v)
			if err != nil {
				return ValidationError{}, ErrNoValidTagFormat
			}
			if field < val {
				return ValidationError{fieldName, ErrLessThanMin}, nil
			}
		case TagIn:
			if !strings.Contains(v, string(rune(field))) {
				return ValidationError{fieldName, ErrValueOutOfRange}, nil
			}
		default:
			return ValidationError{}, ErrNoValidTagNotProvider
		}
	}
	return ValidationError{fieldName, nil}, nil
}

func validateString(fieldName string, field string, tags map[string]string) (ValidationError, error) {
	for k, v := range tags {
		switch k {
		case TagRequire:
			if field == "" {
				return ValidationError{fieldName, ErrRequiredField}, nil
			}
		case TagLength:
			lenField, err := strconv.Atoi(v)
			if err != nil {
				return ValidationError{}, ErrNoValidTagFormat
			}
			if len([]rune(field)) != lenField {
				return ValidationError{fieldName, ErrInvalidLength}, nil
			}
		case TagRegex:
			re, err := regexp.Compile(v)
			if err != nil {
				return ValidationError{}, ErrNoValidRegularExpr
			}
			if !re.MatchString(field) {
				return ValidationError{fieldName, ErrInvalidFormat}, nil
			}
		case TagIn:
			parts := strings.Split(v, ",")
			if !slices.Contains(parts, field) {
				return ValidationError{fieldName, ErrValueNotInList}, nil
			}

		default:
			return ValidationError{}, ErrNoValidTagNotProvider
		}
	}
	return ValidationError{fieldName, nil}, nil
}

func parseTags(tagString string) map[string]string {
	tagPairs := make(map[string]string)
	part := tagString
	var listTag []string
	splitTag := strings.IndexAny(part, "|")
	if splitTag != -1 {
		partOne := part[:splitTag]
		partTwo := part[splitTag+1:]
		listTag = append(listTag, partOne)
		listTag = append(listTag, partTwo)
	} else {
		listTag = append(listTag, part)
	}
	for _, res := range listTag {
		keyEnd := strings.IndexAny(res, "=: ")
		if keyEnd == -1 {
			tagPairs[res] = ""
			continue
		}
		key := res[:keyEnd]
		val := res[keyEnd+1:]
		if len(val) > 0 && val[0] == '"' {
			val = strings.Trim(val, `"`)
		}
		tagPairs[key] = val
	}
	return tagPairs
}
