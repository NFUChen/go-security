package entity

import (
	"github.com/labstack/echo/v4"
	"reflect"
)

type ISerializable interface {
	Serialize() any
}

type ResponseEntity[T any] struct {
	StatusCode int `json:"status_code"`
	Message    T   `json:"message"`
}

// serializeSlice serializes a slice, checking each element for the ISerializable interface.
func serializeSlice(slice reflect.Value) []any {
	serializedSlice := make([]any, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		element := slice.Index(i).Interface()
		if serializableElement, ok := element.(ISerializable); ok {
			serializedSlice[i] = serializableElement.Serialize()
		} else if reflect.ValueOf(element).Kind() == reflect.Slice {
			serializedSlice[i] = serializeSlice(reflect.ValueOf(element))
		} else {
			serializedSlice[i] = element
		}
	}
	return serializedSlice
}

func (entity *ResponseEntity[T]) serialize() any {
	// Check if Message implements ISerializable
	if serializable, ok := any(entity.Message).(ISerializable); ok {
		return map[string]any{
			"status_code": entity.StatusCode,
			"message":     serializable.Serialize(),
		}
	}

	value := reflect.ValueOf(entity.Message)
	switch value.Kind() {
	case reflect.Slice:
		return map[string]any{
			"status_code": entity.StatusCode,
			"message":     serializeSlice(value),
		}

	case reflect.Map:
		serializedMap := make(map[any]any)
		for _, key := range value.MapKeys() {
			originalValue := value.MapIndex(key).Interface()
			if serializableValue, ok := originalValue.(ISerializable); ok {
				serializedMap[key.Interface()] = serializableValue.Serialize()
			} else if reflect.ValueOf(originalValue).Kind() == reflect.Slice {
				serializedMap[key.Interface()] = serializeSlice(reflect.ValueOf(originalValue))
			} else {
				serializedMap[key.Interface()] = originalValue
			}
		}
		return map[string]any{
			"status_code": entity.StatusCode,
			"message":     serializedMap,
		}
	default:
		return map[string]any{
			"status_code": entity.StatusCode,
			"message":     entity.Message,
		}
	}

}

func NewResponse[T any](ctx echo.Context, statusCode int, message T) error {
	entity := &ResponseEntity[T]{
		StatusCode: statusCode,
		Message:    message,
	}
	return ctx.JSON(statusCode, entity.serialize())
}
