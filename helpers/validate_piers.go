package helpers

import (
	"reflect"

	goerrors "github.com/TudorHulban/go-errors"
)

func ValidatePiers(piers any) error {
	if piers == nil {
		return goerrors.ErrValidation{
			Caller: "ValidatePiers",
			Issue: goerrors.ErrNilInput{
				InputName: "piers",
			},
		}
	}

	piersType := reflect.TypeOf(piers)
	piersValue := reflect.ValueOf(piers)

	piersKind := piersType.Kind()

	if piersKind == reflect.Ptr {
		piersType = piersType.Elem()
		piersValue = piersValue.Elem()
		piersKind = piersType.Kind()
	}

	switch piersKind {
	case reflect.Struct:
		for fieldIndex := 0; fieldIndex < piersType.NumField(); fieldIndex++ {
			field := piersType.Field(fieldIndex)
			fieldValue := piersValue.Field(fieldIndex)

			switch field.Type.Kind() {
			case reflect.Ptr, reflect.Interface:
				if fieldValue.IsNil() {
					return goerrors.ErrValidation{
						Caller: "ValidatePiers",
						Issue: goerrors.ErrNilInput{
							InputName: field.Name,
						},
					}
				}

			default:
				continue
			}
		}

	default:
		return goerrors.ErrValidation{
			Caller: "ValidatePiers",
			Issue: goerrors.ErrInvalidInput{
				InputName: "piers",
			},
		}
	}

	return nil
}
