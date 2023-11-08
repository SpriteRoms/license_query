//go:build js && wasm

package main

import (
	"embed"
	_ "embed"
	"reflect"
	"syscall/js"
)

func main() {
	wait := make(chan struct{})
	js.Global().Set("queryLicense", js.FuncOf(queryLicense))
	if !js.Global().Get("INITED_EVENT").IsUndefined() {
		js.Global().Call("INITED_EVENT")
	}
	<-wait
}

func structToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Struct {
		typ := reflect.TypeOf(data)
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			fieldName := typ.Field(i).Name
			fieldValue := field.Interface()
			if field.Kind() == reflect.Struct {
				result[fieldName] = structToMap(fieldValue)
			} else if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
				arr := make([]interface{}, field.Len())
				for i := 0; i < field.Len(); i++ {
					subField := field.Index(i)
					if subField.Kind() == reflect.Struct {
						arr[i] = structToMap(subField.Interface())
					} else {
						arr[i] = subField.Interface()
					}
				}
				result[fieldName] = arr
			} else {
				result[fieldName] = fieldValue
			}
		}
	}

	return result
}

//go:embed web_licenses/*.dat
var webFs embed.FS

type JsLicense struct {
	Contact  string
	Licenses []WebLicense
}

func queryLicense(this js.Value, args []js.Value) interface{} {
	contact := HexHash256([]byte(args[0].String()))
	data, err := webFs.ReadFile("web_licenses/" + contact + ".dat")
	if err != nil {
		// return js.ValueOf(fmt.Sprintf("webFs.ReadFile: %v", err))
		return js.ValueOf(nil)
	}
	license, err := UnmarshalLicense(data)
	if err != nil {
		// return js.ValueOf(fmt.Sprintf("UnmarshalLicense: %v", err))
		return js.ValueOf(nil)
	}

	jsData := structToMap(*license)

	return js.ValueOf(jsData)
}
