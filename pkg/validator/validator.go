package validator

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func ValidJsonFields(v any) error {
	return validateFields("json", v)
}

func ValidYamlFields(v any) error {
	return validateFields("yaml", v)
}

func ValidEnvFields(v any) error {
	return validateFields("env", v)
}

func validateFields(structTagName string, v any) error {
	const omitempty = "omitempty"

	val := reflect.ValueOf(v).Elem()

	var missingRequiredFields []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := val.Type().Field(i)

		tag := structField.Tag.Get(structTagName)
		if !strings.Contains(tag, omitempty) && field.IsZero() {
			missingRequiredFields = append(missingRequiredFields, structField.Tag.Get(structTagName))
		}
	}

	if missingRequiredFields != nil {
		return errors.Errorf("missing reuqired fields: %s", strings.Join(missingRequiredFields, ", "))
	}

	return nil
}
