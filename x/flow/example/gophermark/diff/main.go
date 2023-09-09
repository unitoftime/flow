package main

import (
	"fmt"
	"reflect"
)

// Diff represents the differences between two structs.
type Diff map[string]interface{}

// FindDiff recursively finds the differences between two structs.
func FindDiff(left, right interface{}) Diff {
	diff := make(Diff)
	findDiff(reflect.ValueOf(left), reflect.ValueOf(right), "", diff)
	return diff
}

func findDiff(left, right reflect.Value, path string, diff Diff) {
	if left.Kind() != right.Kind() {
		diff[path] = right.Interface()
		return
	}

	switch left.Kind() {
	case reflect.Struct:
		for i := 0; i < left.NumField(); i++ {
			fieldName := left.Type().Field(i).Name
			findDiff(left.Field(i), right.Field(i), path+"."+fieldName, diff)
		}
	case reflect.Slice:
		for i := 0; i < left.Len(); i++ {
			findDiff(left.Field(i), right.Field(i), path+"."+"slice", diff)
		}
	case reflect.Ptr:
		if left.IsNil() && !right.IsNil() {
			diff[path] = right.Interface()
		} else if !left.IsNil() && right.IsNil() {
			diff[path] = nil
		} else if !left.IsNil() && !right.IsNil() {
			findDiff(left.Elem(), right.Elem(), path, diff)
		}
	default:
		if !reflect.DeepEqual(left.Interface(), right.Interface()) {
			diff[path] = right.Interface()
		}
	}
}

// // Merge applies a diff to a struct and returns the resulting struct.
// func Merge(original interface{}, diff Diff) interface{} {
// 	merged := reflect.ValueOf(original).Interface()
// 	merge(merged, diff)
// 	return merged
// }

// func merge(data interface{}, diff Diff) {
// 	val := reflect.ValueOf(data)
// 	typ := val.Type()

// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		fieldName := field.Name
// 		if newVal, ok := diff[fieldName]; ok {
// 			fieldValue := reflect.ValueOf(newVal)
// 			if val.Field(i).CanSet() && val.Field(i).Type().AssignableTo(fieldValue.Type()) {
// 				val.Field(i).Set(fieldValue)
// 			}
// 		}
// 	}
// }

func Merge(original interface{}, diff Diff) interface{} {
	originalValue := reflect.ValueOf(original).Elem()
	mergeStruct(originalValue, diff)
	return originalValue
}

func mergeStruct(val reflect.Value, diff Diff) {
	for fieldName, fieldValue := range diff {
		println(fieldName)
		field := val.FieldByName(fieldName)
		fmt.Println(field)
		if field.IsValid() {
			fieldType := field.Type()
			fieldValue := reflect.ValueOf(fieldValue)

			if fieldType == fieldValue.Type() {
				field.Set(fieldValue)
			} else if fieldType.Kind() == reflect.Struct && fieldValue.Kind() == reflect.Map {
				// Handle nested struct
				nestedDiff, ok := fieldValue.Interface().(Diff)
				if ok {
					nestedField := val.FieldByName(fieldName)
					if nestedField.IsValid() {
						mergeStruct(nestedField, nestedDiff)
					}
				}
			}
		}
	}
}

type MyStruct struct {
	Name    string
	Age     int
	Address struct {
		City    string
		ZipCode int
	}
}

func main() {
	// testStructToMap()
	testMapDiff()

	// left := MyStruct{
	// 	Name: "Alice",
	// 	Age:  30,
	// 	Address: struct {
	// 		City    string
	// 		ZipCode int
	// 	}{
	// 		City:    "New York",
	// 		ZipCode: 10001,
	// 	},
	// }

	// right := MyStruct{
	// 	Name: "Alice",
	// 	Age:  31,
	// 	Address: struct {
	// 		City    string
	// 		ZipCode int
	// 	}{
	// 		City:    "New York",
	// 		ZipCode: 10002,
	// 	},
	// }

	// diff := FindDiff(left, right)
	// fmt.Printf("Diff: %#v\n", diff)

	// merged := Merge(&left, diff)
	// fmt.Printf("Merged: %#v\n", merged)
}

// StructToMap recursively converts a struct into a map[string]interface{}.
func StructToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	structValue := reflect.ValueOf(input)
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Check if the field is exportable (starts with an uppercase letter)
		if fieldType.PkgPath == "" {
			fieldName := fieldType.Name

			// Check if the field is a struct and should be recursively converted
			if field.Kind() == reflect.Struct {
				result[fieldName] = StructToMap(field.Interface())
			} else {
				// For non-struct fields, convert to interface{}
				result[fieldName] = field.Interface()
			}
		}
	}

	return result
}

// Example struct
type Person struct {
	Name    string
	Age     int
	Address struct {
		City    string
		ZipCode int
	}
}

func testStructToMap() {
	person := Person{
		Name: "Alice",
		Age:  30,
		Address: struct {
			City    string
			ZipCode int
		}{
			City:    "New York",
			ZipCode: 10001,
		},
	}

	result := StructToMap(person)
	fmt.Printf("StructToMap: %#v\n", result)
}

func testMapDiff() {
	// Create two example maps with nested maps
	left := map[string]interface{}{
		"Name": "Alice",
		"Age":  30,
		"Address": map[string]interface{}{
			"City":    "New York",
			"ZipCode": 10001,
		},
	}

	right := map[string]interface{}{
		"Name": "Bob",
		"City": "San Francisco",
		"PhoneNumber": map[string]interface{}{
			"Home": "555-1234",
			"Work": "555-5678",
		},
	}

	// Compute the difference between the maps
	diff := DiffMaps(left, right)
	fmt.Printf("Diff: %#v\n", diff)

	// Merge the difference back into the original map
	merged := MergeDiff(left, diff)
	fmt.Printf("Merged: %#v\n", merged)
}

// DiffMaps recursively computes the difference between two maps.
func DiffMaps(left, right map[string]interface{}) map[string]interface{} {
	diff := make(map[string]interface{})

	// Find keys that are in left but not in right or have different values
	for key, rightValue := range right {
		leftValue, exists := left[key]

		if !exists {
			// Key exists in left but not in right
			diff[key] = leftValue
		} else {
			// Key exists in both maps, compare values
			switch rightValue.(type) {
			case map[string]any:
				// If leftValue is a map, recursively compute the difference
				rightValue, ok := rightValue.(map[string]interface{})
				if ok {
					nestedDiff := DiffMaps(leftValue.(map[string]interface{}), rightValue)
					if len(nestedDiff) > 0 {
						diff[key] = nestedDiff
					}
				} else {
					diff[key] = leftValue
				}
			default:
				if !reflect.DeepEqual(leftValue, rightValue) {
					diff[key] = leftValue
				}
			}
		}
	}
	return diff
}

// MergeDiff applies a diff to a map and returns the resulting map.
func MergeDiff(original map[string]interface{}, diff map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for key, value := range original {
		merged[key] = value
	}

	for key, value := range diff {
		if nestedDiff, ok := value.(map[string]interface{}); ok {
			// If the value is a map (nested diff), recursively merge it
			if originalValue, exists := merged[key]; exists {
				originalNested, ok := originalValue.(map[string]interface{})
				if ok {
					merged[key] = MergeDiff(originalNested, nestedDiff)
				} else {
					merged[key] = nestedDiff
				}
			} else {
				merged[key] = nestedDiff
			}
		} else {
			// For non-nested diffs, directly apply the value
			merged[key] = value
		}
	}

	return merged
}

