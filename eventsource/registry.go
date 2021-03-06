package eventsource

import (
	"fmt"
	"reflect"
	"strings"
)

var registry = make(map[string]reflect.Type)

func RegisterEventType(source EventData) {
	rawType, name := GetTypeName(source)
	registry[name] = rawType
}

func GetEventOfType(name string) (EventData, error) {
	rawType, ok := registry[name]

	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}

	return reflect.New(rawType).Interface().(EventData), nil
}

// GetTypeName of given struct
func GetTypeName(source interface{}) (reflect.Type, string) {
	rawType := reflect.TypeOf(source)

	// source is a pointer, convert to its value
	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}

	name := rawType.String()
	// we need to extract only the name without the package
	// name currently follows the format `package.StructName`
	parts := strings.Split(name, ".")
	return rawType, parts[1]
}
