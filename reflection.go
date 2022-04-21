// Copyright 2022 PerimeterX. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package marshmallow

import (
	"encoding/json"
	"reflect"
	"strings"
)

var unmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

type reflectionInfo struct {
	i int
	t reflect.Type
}

func mapStructFields(target interface{}) map[string]reflectionInfo {
	t := reflectStructType(target)
	result := cacheLookup(t)
	if result != nil {
		return result
	}
	num := t.NumField()
	result = make(map[string]reflectionInfo, num)
	for i := 0; i < num; i++ {
		field := t.Field(i)
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			continue
		}
		if index := strings.Index(name, ","); index > -1 {
			name = name[:index]
		}
		result[name] = reflectionInfo{
			i: i,
			t: field.Type,
		}
	}
	cacheStore(t, result)
	return result
}

func reflectStructValue(target interface{}) reflect.Value {
	v := reflect.ValueOf(target)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func reflectStructType(target interface{}) reflect.Type {
	t := reflect.TypeOf(target)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

var primitiveConverters = map[reflect.Kind]func(v interface{}) (interface{}, bool){
	reflect.Bool: func(v interface{}) (interface{}, bool) {
		res, ok := v.(bool)
		return res, ok
	},
	reflect.Int: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return int(res), true
		}
		return v, false
	},
	reflect.Int8: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return int8(res), true
		}
		return v, false
	},
	reflect.Int16: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return int16(res), true
		}
		return v, false
	},
	reflect.Int32: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return int32(res), true
		}
		return v, false
	},
	reflect.Int64: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return int64(res), true
		}
		return v, false
	},
	reflect.Uint: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return uint(res), true
		}
		return v, false
	},
	reflect.Uint8: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return uint8(res), true
		}
		return v, false
	},
	reflect.Uint16: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return uint16(res), true
		}
		return v, false
	},
	reflect.Uint32: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return uint32(res), true
		}
		return v, false
	},
	reflect.Uint64: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return uint64(res), true
		}
		return v, false
	},
	reflect.Float32: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return float32(res), true
		}
		return v, false
	},
	reflect.Float64: func(v interface{}) (interface{}, bool) {
		res, ok := v.(float64)
		if ok {
			return res, true
		}
		return v, false
	},
	reflect.Interface: func(v interface{}) (interface{}, bool) {
		return v, true
	},
	reflect.String: func(v interface{}) (interface{}, bool) {
		res, ok := v.(string)
		return res, ok
	},
}

func assignValue(field reflect.Value, value interface{}) {
	if value == nil {
		return
	}
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Type().AssignableTo(field.Type()) {
		field.Set(reflectValue)
	}
}

func isValidValue(v interface{}) bool {
	value := reflect.ValueOf(v)
	return value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct && !value.IsNil()
}

func safeReflectValue(t reflect.Type, v interface{}) reflect.Value {
	if v == nil {
		return reflect.Zero(t)
	}
	return reflect.ValueOf(v)
}