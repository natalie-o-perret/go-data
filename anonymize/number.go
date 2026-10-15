package anonymize

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func (a *identifierAnonymizer) number(value any) (any, error) {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return value, nil
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.signed(rv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.unsigned(rv)
	case reflect.Float32, reflect.Float64:
		return a.floating(rv)
	default:
		return value, nil
	}
}

func (a *identifierAnonymizer) signed(value reflect.Value) (any, error) {
	original := strconv.FormatInt(value.Int(), 10)
	replacement, err := a.mapped(fmt.Sprintf("%s:%s", value.Type(), original), func() (string, error) {
		for {
			candidate, err := a.randomDecimal(original)
			if err != nil {
				return "", err
			}
			if _, err := strconv.ParseInt(candidate, 10, value.Type().Bits()); err == nil {
				return candidate, nil
			}
		}
	})
	if err != nil {
		return nil, err
	}
	parsed, _ := strconv.ParseInt(replacement, 10, value.Type().Bits())
	result := reflect.New(value.Type()).Elem()
	result.SetInt(parsed)
	return result.Interface(), nil
}

func (a *identifierAnonymizer) unsigned(value reflect.Value) (any, error) {
	original := strconv.FormatUint(value.Uint(), 10)
	replacement, err := a.mapped(fmt.Sprintf("%s:%s", value.Type(), original), func() (string, error) {
		for {
			candidate, err := a.randomDecimal(original)
			if err != nil {
				return "", err
			}
			if _, err := strconv.ParseUint(candidate, 10, value.Type().Bits()); err == nil {
				return candidate, nil
			}
		}
	})
	if err != nil {
		return nil, err
	}
	parsed, _ := strconv.ParseUint(replacement, 10, value.Type().Bits())
	result := reflect.New(value.Type()).Elem()
	result.SetUint(parsed)
	return result.Interface(), nil
}

func (a *identifierAnonymizer) floating(value reflect.Value) (any, error) {
	if math.IsInf(value.Float(), 0) || math.IsNaN(value.Float()) || math.Trunc(value.Float()) != value.Float() {
		return value.Interface(), nil
	}
	original := strconv.FormatFloat(value.Float(), 'f', -1, value.Type().Bits())
	replacement, err := a.mapped(fmt.Sprintf("%s:%s", value.Type(), original), func() (string, error) {
		return a.randomDecimal(original)
	})
	if err != nil {
		return nil, err
	}
	parsed, err := strconv.ParseFloat(replacement, value.Type().Bits())
	if err != nil || math.IsInf(parsed, 0) {
		return nil, fmt.Errorf("anonymize: numeric identifier out of range")
	}
	result := reflect.New(value.Type()).Elem()
	result.SetFloat(parsed)
	return result.Interface(), nil
}
