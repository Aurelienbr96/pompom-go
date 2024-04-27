package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateInsertSQL(t interface{}, tableName string) string {
	val := reflect.ValueOf(t)
	typeOfT := val.Type()

	fields := []string{}
	values := []string{}

	for i := 0; i < val.NumField(); i++ {
		field := typeOfT.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			fields = append(fields, dbTag)
			values = append(values, fmt.Sprintf(":%s", dbTag))
		}
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(fields, ", "),
		strings.Join(values, ", "),
	)
}
