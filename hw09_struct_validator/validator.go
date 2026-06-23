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
	ErrRequiredField   = errors.New("field required")
	ErrInvalidLength   = errors.New("invalid length")
	ErrValueOutOfRange = errors.New("value out of range")
	ErrInvalidFormat   = errors.New("invalid format")
	ErrValueNotInList  = errors.New("value not in allowed list")
	ErrTagNotProvided  = errors.New("the tag is not provided by the validator")
	ErrLessThanMin     = errors.New("field is less than min")
	ErrGraterThanMax   = errors.New("field is greater than max")
	ErrExpectedStruct  = errors.New("expected a struct")
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
		switch field.Type.Kind() { //nolint
		case reflect.String:
			err := validateString(field.Name, val.Field(i).String(), tags)
			if err.Err == nil {
				continue
			}
			if errors.Is(err.Err, ErrTagNotProvided) {
				return errors.New("tag not provided, validation break")
			}
			result = append(result, err)
		case reflect.Int:
			err := validateInt(field.Name, int(val.Field(i).Int()), tags)
			if errors.Is(err.Err, ErrTagNotProvided) {
				return errors.New("tag not provided, validation break")
			}
			if err.Err == nil {
				continue
			}
			result = append(result, err)
		case reflect.Slice:
			switch kind := field.Type.Elem().Kind(); kind { //nolint
			case reflect.String:
				err := validateSliceString(field.Name, val.Field(i), tags)
				if err.Err == nil {
					continue
				}
				result = append(result, err)
			case reflect.Int:
				err := validateSliceInt(field.Name, val.Field(i), tags)
				if err.Err == nil {
					continue
				}
				result = append(result, err)
			}
		}
	}
	return result
}

func validateSliceString(fieldName string, field reflect.Value, tags map[string]string) ValidationError {
	for k, v := range tags {
		switch k {
		case TagLength:
			lenField, err := strconv.Atoi(v)
			if err != nil {
				return ValidationError{fieldName, ErrTagNotProvided}
			}
			tmp := field.Interface().([]string)
			for _, str := range tmp {
				if len(str) != lenField {
					return ValidationError{fieldName, ErrInvalidLength}
				}
			}
		default:
			return ValidationError{fieldName, ErrTagNotProvided}
		}
	}
	return ValidationError{fieldName, nil}
}

func validateSliceInt(fieldName string, field reflect.Value, tags map[string]string) ValidationError {
	for k, v := range tags {
		switch k {
		case TagLength:
			lenField, _ := strconv.Atoi(v)
			if len(field.Interface().([]int)) != lenField {
				return ValidationError{fieldName, ErrInvalidLength}
			}
		default:
			return ValidationError{fieldName, ErrTagNotProvided}
		}
	}
	return ValidationError{fieldName, nil}
}

func validateInt(fieldName string, field int, tags map[string]string) ValidationError {
	for k, v := range tags {
		switch k {
		case TagMax:
			val, err := strconv.Atoi(v)
			if err != nil {
				return ValidationError{fieldName, ErrTagNotProvided}
			}
			if field > val {
				return ValidationError{fieldName, ErrGraterThanMax}
			}
		case TagMin:
			val, err := strconv.Atoi(v)
			if err != nil {
				return ValidationError{fieldName, ErrTagNotProvided}
			}
			if field < val {
				return ValidationError{fieldName, ErrLessThanMin}
			}
		case TagIn:
			if !strings.Contains(v, string(rune(field))) {
				return ValidationError{fieldName, ErrValueOutOfRange}
			}
		default:
			return ValidationError{fieldName, ErrTagNotProvided}
		}
	}
	return ValidationError{fieldName, nil}
}

func validateString(fieldName string, field string, tags map[string]string) ValidationError {
	for k, v := range tags {
		switch k {
		case TagRequire:
			if field == "" {
				return ValidationError{fieldName, ErrRequiredField}
			}
		case TagLength:
			lenField, err := strconv.Atoi(v)
			if err != nil {
				return ValidationError{fieldName, ErrTagNotProvided}
			}
			if len([]rune(field)) != lenField {
				return ValidationError{fieldName, ErrInvalidLength}
			}
		case TagRegex:
			re, err := regexp.Compile(v)
			if err != nil {
				return ValidationError{fieldName, ErrTagNotProvided}
			}
			if !re.MatchString(field) {
				return ValidationError{fieldName, ErrInvalidFormat}
			}
		case TagIn:
			parts := strings.Split(v, ",")
			if !slices.Contains(parts, field) {
				return ValidationError{fieldName, ErrValueNotInList}
			}

		default:
			return ValidationError{fieldName, ErrTagNotProvided}
		}
	}
	return ValidationError{fieldName, nil}
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
