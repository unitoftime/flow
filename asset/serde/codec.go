package serde

// // A (int): 555
// // B (string): STRINGBBBBB
// // [ 0 (int): 1, 1 (int): 2, 2 (string): a, 3 (string): b, 4 (serde.MyData2): {MyData2Struct}, ]
// // D . @T (string): serde.MyData
// // @V . A (int): 0
// // B (string): MyData - D
// // [ ]
// // D . @T (string): serde.MyData3
// // @V . Number (int): 12345
// // A (int): 555
// // B (string): STRINGBBBBB
// // [ 0 (int): 1, 1 (int): 2, 2 (string): a, 3 (string): b, 4 (serde.MyData2): {MyData2Struct}, ]
// // D (serde.MyData): {0 MyData - D [] {12345}}

// // func Encode(input any) any {
// // 	m := make(map[string]any)
// // 	config := &mapstructure.DecoderConfig{
// // 		// WeaklyTypedInput: true,
// // 		DecodeHook: encodeHookFunc(),
// // 		Result: &m,
// // 	}
// // 	decoder, err := mapstructure.NewDecoder(config)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	err = decoder.Decode(input)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	return m
// // }

// // func encodeHookFunc() mapstructure.DecodeHookFunc {
// // 	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
// // 		fmt.Printf("encodeHookFunc: from: %s, to: %s, data: %T\n", f.Kind(), t.Kind(), data)
// // 		return data, nil
// // 	}
// // }

// // func Decode(output any, input any) {
// // 	config := &mapstructure.DecoderConfig{
// // 		WeaklyTypedInput: true,
// // 		DecodeHook: decodeHookFunc(),
// // 		Result: output,
// // 	}
// // 	decoder, err := mapstructure.NewDecoder(config)
// // 	if err != nil {
// // 		panic(err)
// // 		// return err
// // 	}
// // 	err = decoder.Decode(input)
// // 	if err != nil {
// // 		panic(err)
// // 		// return err
// // 	}
// // 	// fmt.Printf("Decoded: %T: %v\n", targ.Get(), targ.Get())

// // 	// return nil
// // }

// // func decodeHookFunc() mapstructure.DecodeHookFunc {
// // 	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
// // 		fmt.Printf("decodeHookFunc: from: %s, to: %s, data: %T\n", f.Kind(), t.Kind(), data)
// // 		// return data, nil
// // 		if f.Kind() == reflect.Slice || f.Kind() == reflect.Array {
// // 			fmt.Println("SLICE")
// // 		}

// // 		if f.Kind() != reflect.Map {
// // 			return data, nil
// // 		}
// // 		if t.Kind() != reflect.Interface {
// // 			return data, nil
// // 		}
// // 		m, ok := data.(map[string]any)
// // 		if !ok {
// // 			return data, nil
// // 		}
// // 		typeName, ok := m["@T"]
// // 		if !ok {
// // 			return data, nil
// // 		}
// // 		wrappedValue, ok := m["@V"]
// // 		if !ok {
// // 			return data, nil
// // 		}

// // 		fmt.Println("decodeHookFunc: using type ", typeName)

// // 		targ, ok := registry.Get(typeName.(string))
// // 		if !ok {
// // 			panic(fmt.Sprintf("decodeHookFunc: Could not find typeName in registry: %s", typeName)) // TODO: Dont panic
// // 			return data, nil // TODO: warn? Unregistered type?
// // 		}

// // 		newTarg := targ.New()
// // 		Decode(newTarg.Ptr(), wrappedValue)
// // 		// TODO: check errors
// // 		ret := newTarg.Get()

// // 		return ret, nil

// // // 			if len(m) != 1 {
// // // 				return data, nil
// // // 			}
// // // for k, v := range m {
// // // 				targ, ok := registry.Get(k)
// // // 				if !ok {
// // // 					panic(fmt.Sprintf("Couldnt find: %s", k)) // TODO: Dont panic
// // // 					return data, nil // TODO: warn? Unregistered type?
// // // 				}
// // // 				var ret any
// // // 				switch t := v.(type) {
// // // 				case map[string]any:
// // // 					newTarg := targ.New()
// // // 					err := Decode(t, newTarg.Ptr())
// // // 					if err != nil {
// // // 						return data, err // TODO: return nil? or err?
// // // 					}
// // // 					ret = newTarg.Get()
// // 	}
// // }

// //--------------------------------------------------------------------------------
// // type wrappedType struct {
// // 	T string // TODO: uint64 hash?
// // 	V any
// // }

// func Encode(input any) any {
// 	if input == nil { return nil }
// 	rv := reflect.ValueOf(input)
// 	// if rv.IsZero() { return input }
// 	fmt.Printf("Encode: %T\n", rv.Interface())

// 	switch rv.Kind() {
// 	case reflect.Array, reflect.Slice:
// 		numElems := rv.Len()
// 		output := make([]any, numElems)
// 		for i := 0; i < numElems; i++ {
// 			idxValue := rv.Index(i)
// 			idxKind := idxValue.Kind()

// 			idxOutput := Encode(idxValue.Interface())
// 			if idxKind == reflect.Interface {
// 				idxTypeName, ok := registry.fromType(idxValue.Interface())
// 				if !ok {
// 					fmt.Printf("Missing: %T\n", idxValue.Interface())
// 					panic(fmt.Sprintf("serde: Unregistered type in slice index: %d: %T", i, idxValue.Interface()))
// 				}
// 				fmt.Println("SLICE OF INTERFACES", idxTypeName)
// 				wrapped := map[string]any{
// 					"@T": idxTypeName,
// 					"@V": idxOutput,
// 				}
// 				output[i] = wrapped
// 			} else {
// 				output[i] = idxOutput
// 			}
// 		}
// 		return output
// 	case reflect.Struct:
// 		output := make(map[string]any)
// 		structType := rv.Type()
// 		for i := 0; i < rv.NumField(); i++ {
// 			structField := structType.Field(i)
// 			if !structField.IsExported() { continue } // Skip non exported fields
// 			sfValue := rv.Field(i)
// 			sfKind := structField.Type.Kind()

// 			iface := sfValue.Interface()
// 			if iface == nil { continue }
// 			fmt.Printf("field %s: %T\n", structField.Name, iface)
// 			sfOutput := Encode(sfValue.Interface())

// 			if sfKind == reflect.Interface {
// 				sfTypeName, ok := registry.fromType(sfValue.Interface())
// 				if !ok {
// 					fmt.Printf("Missing: %T\n", sfValue.Interface())
// 					panic(fmt.Sprintf("serde: Unregistered type: %s: %T", structField.Name, sfValue.Interface()))
// 				}
// 				fmt.Println("INTERFACE", structField.Name, sfTypeName)
// 				// map{
// 				// 	"@T": type,
// 				// 	"@V": val,
// 				// }
// 				wrapped := map[string]any{
// 					"@T": sfTypeName,
// 					"@V": sfOutput,
// 				}
// 				// wrapped := wrappedType{
// 				// 	T: sfTypeName,
// 				// 	V: sfOutput,
// 				// }
// 				output[structField.Name] = wrapped
// 			} else {
// 				output[structField.Name] = sfOutput
// 			}
// 		}
// 		return output
// 	// case reflect.Array, reflect.Slice:
// 	// case reflect.Map:
// 	// default:
// 	case reflect.Pointer:
// 		return Encode(rv.Elem())
// 	case reflect.Interface:
// 		fmt.Println("AAAA")
// 	}
// 	return input
// }

// func Decode(output any, input any) {
// 	rv := reflect.ValueOf(output)

// 	rval, isReflectVal := output.(reflect.Value)
// 	if isReflectVal {
// 		fmt.Println("IsAlreadyReflected")
// 		rv = rval
// 	}

// 	inputSlice, isSlice := input.([]any)
// 	if isSlice {
// 		if rv.Kind() == reflect.Pointer {
// 			rv = reflect.Indirect(rv)
// 		}

// 		sliceType := rv.Type().Elem()
// 		for i := range inputSlice {
// 			appendVal := reflect.Zero(sliceType)
// 			fmt.Println("Decode Slice", i, inputSlice[i])
// 			rv.Set(reflect.Append(rv, appendVal))
// 			idxVal := rv.Index(i)
// 			Decode(idxVal, inputSlice[i])
// 		}
// 		return
// 	}

// 	str, isString := input.(string)
// 	if isString {
// 		rv.Set(reflect.ValueOf(str))
// 		return
// 	}

// 	inputMap, isMap := input.(map[string]any)
// 	if !isMap {
// 		fmt.Printf("OUT, IN: %T, %T\n", output, input)
// 		rv.Elem().Set(reflect.ValueOf(input))
// 		return
// 	}

// 	typeName, isTyped := inputMap["@T"]
// 	if isTyped {
// 		// TODO: Do typed map decoding
// 		fmt.Println("Typed Map")
// 		wrappedData, ok := inputMap["@V"]
// 		if !ok { panic("ERROR MISSING WRAPPED VALUE") }

// 		decoder, ok := registry.Get(typeName.(string))
// 		if !ok {
// 			panic(fmt.Sprintf("serde: Unregistered type: %s", typeName))
// 		}

// 		decoder = decoder.New()
// 		decPtr := decoder.Ptr()
// 		Decode(decPtr, wrappedData)

// 		decVal := decoder.Get()

// 		// I think you need to do something special like pass in an addressible interface here then set the field value to that or something. look at how you did imgui stuff
// 		// fmt.Println("value.CanSet:", typeName, wrappedData, decVal, decPtr, rv.CanSet())
// 		// fmt.Println("---")
// 		// fmt.Println(reflect.Indirect(rv), reflect.Indirect(rv).CanSet())
// 		rv.Set(reflect.ValueOf(decVal))

// 	} else {
// 		if rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface  {
// 			fmt.Printf("ISPTR %T\n", rv.Interface())
// 			rv = rv.Elem()
// 		}
// 		if rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
// 			fmt.Printf("ISPTR222: %T\n", rv.Interface())
// 			rv = rv.Elem()
// 		}
// 		if rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
// 			fmt.Printf("ISPTR333: %T\n", rv.Interface())
// 			rv = rv.Elem()
// 		}
// 		for k, v := range inputMap {
// 			fmt.Println("InputMap:", k, v)
// 			fv := rv.FieldByName(k)

// 			switch underlying := v.(type) {
// 			case map[string]any:
// 				// TODO Need the pointer to the field value
// 				fmt.Println("Recursive Decode: map")
// 				Decode(fv, underlying) // Note: we pass the struct rval inside b/c decode can handle that
// 			case []any:
// 				sliceType := fv.Type().Elem()
// 				for i := range underlying {
// 					appendVal := reflect.Zero(sliceType)
// 					fmt.Println("Decode Slice", i, underlying[i])
// 					fv.Set(reflect.Append(fv, appendVal))
// 					idxVal := fv.Index(i)
// 					Decode(idxVal, underlying[i])
// 				}
// 			default:
// 				fmt.Printf("Explicit Set %s: %T, (kind: %s)\n", k, underlying, fv.Kind())
// 				switch fv.Kind() {
// 				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
// 					fmt.Println("INT")
// 					// fv.SetInt(reflect.ValueOf(underlying).Int())
// 					fv.SetInt(decodeInt(underlying))
// 				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 					fmt.Println("UINT")
// 					fv.SetUint(decodeUint(underlying))
// 					// fv.SetUint(reflect.ValueOf(underlying).Uint())
// 				case reflect.Ptr:
// 					fmt.Println("PTR")
// 					if !fv.IsZero() {
// 						fv.Set(reflect.ValueOf(v)) // TODO: I think this is wrong
// 					}
// 				default:
// 					fmt.Println("DEFAULT")
// 					fv.Set(reflect.ValueOf(v))
// 				}
// 				// fv.Set(reflect.ValueOf(v))
// 			}

// 			// switch underlying := v.(type) {
// 			// case map[string]any:
// 			// 	// TODO Need the pointer to the field value
// 			// 	fmt.Println("Recursive Decode: map")
// 			// 	Decode(fv, underlying) // Note: we pass the struct rval inside b/c decode can handle that
// 			// case []any:
// 			// 	sliceType := fv.Type().Elem()
// 			// 	for i := range underlying {

// 			// 		// decodedUnderlying := Decode(
// 			// 		// reflect.ValueOf(underlying[i])

// 			// 		appendVal := reflect.Zero(sliceType)
// 			// 		fmt.Println("Decode Slice", i, underlying[i])
// 			// 		fv.Set(reflect.Append(fv, appendVal))
// 			// 		idxVal := fv.Index(i)
// 			// 		Decode(idxVal, underlying[i])
// 			// 		// idxVal.Set(one)
// 			// 	}
// 			// default:
// 			// 	fmt.Printf("Explicit Set %s: %T\n", k, underlying)
// 			// 	fv.Set(reflect.ValueOf(v))
// 			// }
// 		}
// 	}
// }

// func decodeInt(data any) int64 {
// 	val := reflect.ValueOf(data)
// 	switch val.Kind() {
// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
// 		return int64(val.Int())
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 		return int64(val.Uint())
// 	case reflect.Float32, reflect.Float64:
// 		return int64(val.Float())
// 	}
// 	panic(fmt.Sprintf("Unhandled decodeInt for type %T (kind: %s)", data, val.Kind()))
// }

// func decodeUint(data any) uint64 {
// 	val := reflect.ValueOf(data)
// 	switch val.Kind() {
// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
// 		return uint64(val.Int()) // TODO: negative check?
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 		return uint64(val.Uint())
// 	case reflect.Float32, reflect.Float64:
// 		return uint64(val.Float()) // TODO: negative check?
// 	}
// 	panic(fmt.Sprintf("Unhandled decodeInt for type %T (kind: %s) (val: %v)", data, val.Kind(), data))
// }

// // --------------------------------------------------------------------------------
// // func Encode(output any, input any) {
// // 	rv := reflect.ValueOf(input)

// // 	switch rv.Kind() {
// // 	case reflect.Struct:
// // 		st := rv.Type()
// // 		for i := 0; i < rv.NumField(); i++ {
// // 			sf := st.Field(i)
// // 			sfOutput := make(map[string]any)
// // 			sfValue := rv.Field(i)
// // 			Encode(sfOutput, sfValue.Interface())
// // 			output[sf.Name] = sfOutput
// // 		}

// // 	// case reflect.Array, reflect.Slice:
// // 	// case reflect.Map:
// // 	// default:
// // 	}

// // }

// // type typeWrapped[T any] struct{
// // 	TYPE string // TODO: Hash? uint64
// // 	Data T
// // }

// // func (t typeWrapped[T]) Type() string {
// // 	return t.TYPE
// // }

// // type Typer interface {
// // 	Type() string
// // }

// // func Encode(v any, m map[string]any) (error) {
// // 	// m := make(map[string]any)
// // 	config := &mapstructure.DecoderConfig{
// // 		// WeaklyTypedInput: true,
// // 		DecodeHook: encodeHook(),
// // 		Result: &m,
// // 	}
// // 	decoder, err := mapstructure.NewDecoder(config)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	err = decoder.Decode(v)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	return nil
// // }

// // func encodeHook() mapstructure.DecodeHookFunc {
// // 	// return func(from reflect.Value, to reflect.Value) (interface{}, error) {
// // 	// 	// fmt.Printf("TYPE   | %T\n", data)
// // 	// 	fmt.Printf("STRING | From: %s To: %s\n", from, to)
// // 	// 	// fmt.Printf("KIND   | From: %s To: %s\n", f.Kind(), t.Kind())
// // 	// 	return from.Interface(), nil
// // 	// }
// // 	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
// // 		fmt.Printf("TYPE   | %T\n", data)
// // 		fmt.Printf("STRING | From: %s To: %s\n", f.String(), t.String())
// // 		fmt.Printf("KIND   | From: %s To: %s\n", f.Kind(), t.Kind())

// // 		// m := make(map[string]any)
// // 		// m["TYPE"] = f.String()

// // 		// var aa map[string]any
// // 		// err := Encode(t, aa)
// // 		// if err != nil {
// // 		// 	return data, err // TODO: return nil? or err?
// // 		// }
// // 		// m["VALUE"] = aa

// // 		// // if t.Kind() != reflect.Interface {
// // 		// // 	return data, nil
// // 		// // }

// // 		return data, nil
// // 	}
// // 			// registry.fromType(f)

// // 			// if f.Kind() != reflect.Map {
// // 			// 	return data, nil
// // 			// }
// // 			// // fmt.Printf("Type: %T\n", data)
// // 			// // fmt.Printf("From: %s To: %s\n", f.String(), t.String())
// // 			// if t.Kind() != reflect.Interface {
// // 			// 	return data, nil
// // 			// }

// // 			// m, ok := data.(map[string]any)
// // 			// if !ok {
// // 			// 	return data, nil
// // 			// }
// // 			// if len(m) != 1 {
// // 			// 	return data, nil
// // 			// }
// // 			// for k, v := range m {
// // 			// 	targ, ok := registry.Get(k)
// // 			// 	if !ok {
// // 			// 		panic(fmt.Sprintf("Couldnt find: %s", k)) // TODO: Dont panic
// // 			// 		return data, nil // TODO: warn? Unregistered type?
// // 			// 	}

// // 			// 	var ret any
// // 			// 	switch t := v.(type) {
// // 			// 	case map[string]any:
// // 			// 		newTarg := targ.New()
// // 			// 		err := Decode(t, newTarg.Ptr())
// // 			// 		if err != nil {
// // 			// 			return data, err // TODO: return nil? or err?
// // 			// 		}
// // 			// 		ret = newTarg.Get()
// // 			// 	case string:
// // 			// 		ret, ok = singleTypeConvert(t, targ)
// // 			// 		if !ok { return data, nil }
// // 			// 	case bool:
// // 			// 		ret, ok = singleTypeConvert(t, targ)
// // 			// 		if !ok { return data, nil }
// // 			// 	case int:
// // 			// 		ret, ok = singleTypeConvert(t, targ)
// // 			// 		if !ok { return data, nil }
// // 			// 	case float32:
// // 			// 		ret, ok = singleTypeConvert(t, targ)
// // 			// 		if !ok { return data, nil }
// // 			// 	case float64:
// // 			// 		ret, ok = singleTypeConvert(t, targ)
// // 			// 		if !ok { return data, nil }
// // 			// 	default:
// // 			// 		return data, nil
// // 			// 	}
// // 			// 	// fmt.Println("Special Ret: ", k, ret)
// // 			// 	return ret, nil
// // 			// }
// // 			// return data, nil
// // 		// }
// // }

// // //--------------------------------------------------------------------------------

// // func Decode(m map[string]any, v any) error {
// // 	config := &mapstructure.DecoderConfig{
// // 		WeaklyTypedInput: true,
// // 		// DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
// // 		DecodeHook: mapstructure.ComposeDecodeHookFunc(
// // 			mapstructure.StringToTimeDurationHookFunc(),
// // 			registeredMapHookFunc(),
// // 		),
// // 		Result: v,
// // 	}
// // 	decoder, err := mapstructure.NewDecoder(config)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	err = decoder.Decode(m)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	// fmt.Printf("Decoded: %T: %v\n", targ.Get(), targ.Get())

// // 	return nil
// // }

// // func singleTypeConvert[T any](src T, targ decodeTargeter) (any, bool) {
// // 	value := reflect.ValueOf(src)
// // 	dstType := targ.Type()
// // 	if !value.CanConvert(dstType) {
// // 		return nil, false
// // 	}
// // 	retValue := value.Convert(dstType)
// // 	return retValue.Interface(), true
// // }

// // func registeredMapHookFunc() mapstructure.DecodeHookFunc {
// // 	return func(f reflect.Type, t reflect.Type,
// // 		data interface{}) (interface{}, error) {
// // 			if f.Kind() != reflect.Map {
// // 				return data, nil
// // 			}
// // 			// fmt.Printf("Type: %T\n", data)
// // 			// fmt.Printf("From: %s To: %s\n", f.String(), t.String())
// // 			if t.Kind() != reflect.Interface {
// // 				return data, nil
// // 			}

// // 			m, ok := data.(map[string]any)
// // 			if !ok {
// // 				return data, nil
// // 			}
// // 			if len(m) != 1 {
// // 				return data, nil
// // 			}
// // 			for k, v := range m {
// // 				targ, ok := registry.Get(k)
// // 				if !ok {
// // 					panic(fmt.Sprintf("Couldnt find: %s", k)) // TODO: Dont panic
// // 					return data, nil // TODO: warn? Unregistered type?
// // 				}

// // 				var ret any
// // 				switch t := v.(type) {
// // 				case map[string]any:
// // 					newTarg := targ.New()
// // 					err := Decode(t, newTarg.Ptr())
// // 					if err != nil {
// // 						return data, err // TODO: return nil? or err?
// // 					}
// // 					ret = newTarg.Get()
// // 				case string:
// // 					ret, ok = singleTypeConvert(t, targ)
// // 					if !ok { return data, nil }
// // 				case bool:
// // 					ret, ok = singleTypeConvert(t, targ)
// // 					if !ok { return data, nil }
// // 				case int:
// // 					ret, ok = singleTypeConvert(t, targ)
// // 					if !ok { return data, nil }
// // 				case float32:
// // 					ret, ok = singleTypeConvert(t, targ)
// // 					if !ok { return data, nil }
// // 				case float64:
// // 					ret, ok = singleTypeConvert(t, targ)
// // 					if !ok { return data, nil }
// // 				default:
// // 					return data, nil
// // 				}
// // 				// fmt.Println("Special Ret: ", k, ret)
// // 				return ret, nil
// // 			}
// // 			return data, nil
// // 		}
// // }
