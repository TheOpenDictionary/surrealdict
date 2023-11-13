package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/TheOpenDictionary/odict/lib/core"
	"github.com/iancoleman/strcase"
)

type Identifiable struct {
	ID string
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func createQuery(s interface{}) string {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return ""
	}

	tableName := strcase.ToLowerCamel(strings.Replace(val.Type().Name(), "Representable", "", -1))
	columns := []string{}
	values := []string{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := strcase.ToLowerCamel(val.Type().Field(i).Name)

		if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct {
			for j := 0; j < field.Len(); j++ {
				subInsert := createQuery(field.Index(j).Interface())
				columns = append(columns, fieldName)
				values = append(values, fmt.Sprintf("(%s).id", subInsert))
			}
		} else if field.Kind() == reflect.Map && field.Type().Elem().Kind() == reflect.Struct {
			iter := field.MapRange()

			for iter.Next() {
				subInsert := createQuery(iter.Value().Interface())
				columns = append(columns, fieldName)
				values = append(values, fmt.Sprintf("(%s).id", subInsert))
			}
		} else if fieldName != "id" && fieldName != "xmlname" {
			columns = append(columns, fieldName)
			values = append(values, fmt.Sprintf("'%v'", field.Interface()))
		}
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURN id", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))
}

func main() {
	dict, _ := core.ReadDictionaryFromPath("/Users/tjnickerson/.linguistic/dictionaries/eng-eng.odict")

	d := dict.AsRepresentable()

	print(createQuery(d.Entries["create"]))
	// writeData(db, d)

	fmt.Printf("All done!")
}
