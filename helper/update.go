package helper

import (
	"context"
	"errors"
	"go/ast"
	"reflect"
	"strings"

	"github.com/zzzhr1990/go-mongo-util/name"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateMany Update it auto
func UpdateMany(ctx context.Context, collection *mongo.Collection, filter interface{}, data interface{}, controlMap map[string]bool, opts ...*options.UpdateOptions) (int64, error) {
	if data == nil {
		return 0, errors.New("data is nil")
	}
	reflectType := reflect.ValueOf(data).Type()
	val := reflect.ValueOf(data)
	needElm := false
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		needElm = true

	}

	// Scope value need to be a struct
	if reflectType.Kind() != reflect.Struct {
		return 0, errors.New("cannot support this input struct")
	}
	if needElm {
		if val.Kind() != reflect.Ptr && val.Kind() != reflect.Interface {
			return 0, errors.New("cannot support this input struct")
		}
		val = val.Elem()
	}
	proj := bson.D{}
	change := 0
	for i := 0; i < reflectType.NumField(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {
			tags := parseSimpleTagSetting(fieldStruct.Tag)
			fieldName := ""
			if len(tags) > 0 && len(tags[0]) > 0 {
				fieldName = tags[0]
			} else {
				fieldName = name.ToColumnName(fieldStruct.Name)
			}
			fVal := val.Field(i)
			forceUpdate := false
			forceIgnore := false
			if controlMap != nil {
				fVal, ok := controlMap[fieldName]
				if ok {
					if fVal {
						forceUpdate = true
					} else {
						forceIgnore = true
					}
				}
			}

			if !forceIgnore && (forceUpdate || !fVal.IsZero()) {
				proj = append(proj, bson.E{Key: fieldName, Value: fVal.Interface()})
				change++
			}

		}
	}
	if change < 1 {
		return 0, nil
	}
	updlang := bson.D{bson.E{Key: "$set", Value: proj}}
	if collection == nil {
		return 0, errors.New("connection not available")
	}
	res, err := collection.UpdateMany(ctx, filter, bson.D{bson.E{Key: "$set", Value: updlang}}, opts...)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return res.ModifiedCount, nil
}

func parseTagSetting(tags reflect.StructTag) map[string]string {
	setting := map[string]string{}
	for _, str := range []string{tags.Get("qzmongo"), tags.Get("bson")} {
		if str == "" {
			continue
		}
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}

func parseSimpleTagSetting(tags reflect.StructTag) []string {
	setting := []string{}
	for _, str := range []string{tags.Get("bson")} {
		if str == "" {
			continue
		}
		tags := strings.Split(str, ",")
		for _, value := range tags {
			setting = append(setting, value)
		}
	}
	return setting
}

/*
func Entity(ctx context.Context, collection *mongo.Collection, data interface{}) error {
	val := reflect.ValueOf(data)
	kd := val.Kind()

	if kd == reflect.Ptr {

		typeOfCat := reflect.TypeOf(data)
		// 取类型的元素
		typeOfCat = typeOfCat.Elem()
		// 显示反射类型对象的名称和种类
		// fields := typeOfCat.NumField()
		//获取字段的类型
		val := reflect.ValueOf(data).Elem()
		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			fmt.Printf("---------------->%v: %v\n", typeOfCat.Field(i).Name, valueField.Interface())
		}
	}
	if kd == reflect.Struct {
		typeOfCat := reflect.TypeOf(data)
		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			fmt.Printf("================>%v: %v\n", typeOfCat.Field(i).Name, valueField.Interface())
		}
	}
	return nil
}
*/
