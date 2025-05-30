// Ref: https://github.com/erni27/mapper/
// https://github.com/alexsem80/go-mapper

// Package mapper provides a mapper for structs.
package mapper

import (
	"fmt"
	"reflect"

	"emperror.dev/errors"
	"github.com/iancoleman/strcase"

	linq "github.com/ahmetb/go-linq/v3"

	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	reflectionHelper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/reflectionhelper"
)

var (
	// ErrNilFunction is the error returned by CreateCustomMap or CreateMapWith
	// if a nil function is passed to the method.
	ErrNilFunction = errors.New("mapper: nil function")
	// ErrMapNotExist is the error returned by the Map method
	// if a map for given types does not exist.
	ErrMapNotExist = errors.New("mapper: map does not exist")
	// ErrMapAlreadyExists is the error returned by one of the CreateMap method
	// if a given map already exists. Mapper does not allow to override MapFunc.
	ErrMapAlreadyExists = errors.New("mapper: map already exists")
	// ErrUnsupportedMap is the error returned by CreateMap or CreateMapWith
	// if a source - destination mapping is not supported. The mapping is supported only for
	// structs to structs with at least one exported field by a src side which corresponds to a dst field.
	ErrUnsupportedMap = errors.New("mapper: unsupported map")
)

const (
	// SrcKeyIndex is the index of the source key.
	SrcKeyIndex = iota
	// DestKeyIndex is the index of the destination key.
	DestKeyIndex
)

// mappingsEntry is a struct that contains the source and destination types.
type mappingsEntry struct {
	SourceType      reflect.Type
	DestinationType reflect.Type
}

// typeMeta is a struct that contains the keys to tags and tags to keys.
type typeMeta struct {
	keysToTags map[string]string
	tagsToKeys map[string]string
}

// MapFunc is a function that maps a source type to a destination type.
type MapFunc[TSrc any, TDst any] func(TSrc) TDst

var (
	profiles     = map[string][][2]string{}
	maps         = map[mappingsEntry]interface{}{}
	mapperConfig = &MapperConfig{
		MapUnexportedFields: false,
	}
)

// Configure is a function that configures the mapper.
func Configure(config *MapperConfig) {
	mapperConfig = config
}

// ClearMappings clears all registered mappings and profiles.
// This is primarily used for testing purposes.
func ClearMappings() {
	maps = make(map[mappingsEntry]interface{})
	profiles = make(map[string][][2]string)
}

// validateTypes validates that both source and destination types are structs or pointers to structs.
func validateTypes(srcType, desType reflect.Type) error {
	if (srcType.Kind() != reflect.Struct && (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() != reflect.Struct)) ||
		(desType.Kind() != reflect.Struct && (desType.Kind() == reflect.Ptr && desType.Elem().Kind() != reflect.Struct)) {
		return ErrUnsupportedMap
	}

	return nil
}

// registerPointerStructMapping registers mappings for pointer struct types.
func registerPointerStructMapping(srcType, desType reflect.Type) error {
	pointerStructTypeKey := mappingsEntry{
		SourceType:      srcType,
		DestinationType: desType,
	}
	nonePointerStructTypeKey := mappingsEntry{
		SourceType:      srcType.Elem(),
		DestinationType: desType.Elem(),
	}

	if _, ok := maps[nonePointerStructTypeKey]; ok {
		return ErrMapAlreadyExists
	}
	if _, ok := maps[pointerStructTypeKey]; ok {
		return ErrMapAlreadyExists
	}

	maps[nonePointerStructTypeKey] = nil
	maps[pointerStructTypeKey] = nil

	return nil
}

// registerStructMapping registers mappings for non-pointer struct types.
func registerStructMapping(srcType, desType reflect.Type) error {
	nonePointerStructTypeKey := mappingsEntry{SourceType: srcType, DestinationType: desType}
	pointerStructTypeKey := mappingsEntry{
		SourceType:      reflect.New(srcType).Type(),
		DestinationType: reflect.New(desType).Type(),
	}

	if _, ok := maps[nonePointerStructTypeKey]; ok {
		return ErrMapAlreadyExists
	}
	if _, ok := maps[pointerStructTypeKey]; ok {
		return ErrMapAlreadyExists
	}

	maps[nonePointerStructTypeKey] = nil
	maps[pointerStructTypeKey] = nil

	return nil
}

// getBaseTypes returns the base types (dereferenced if pointer) for source and destination.
func getBaseTypes(srcType, desType reflect.Type) (reflect.Type, reflect.Type) {
	if srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Struct {
		srcType = srcType.Elem()
	}
	if desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Struct {
		desType = desType.Elem()
	}

	return srcType, desType
}

// CreateMap is a function that creates a map.
func CreateMap[TSrc any, TDst any]() error {
	var src TSrc
	var dst TDst
	srcType := reflect.TypeOf(&src).Elem()
	desType := reflect.TypeOf(&dst).Elem()

	if err := validateTypes(srcType, desType); err != nil {
		return err
	}

	if srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Struct {
		if err := registerPointerStructMapping(srcType, desType); err != nil {
			return err
		}
	} else {
		if err := registerStructMapping(srcType, desType); err != nil {
			return err
		}
	}

	srcType, desType = getBaseTypes(srcType, desType)
	configProfile(srcType, desType)

	return nil
}

// CreateCustomMap is a function that creates a custom map.
func CreateCustomMap[TSrc any, TDst any](fn MapFunc[TSrc, TDst]) error {
	if fn == nil {
		return ErrNilFunction
	}
	var src TSrc
	var dst TDst
	srcType := reflect.TypeOf(&src).Elem()
	desType := reflect.TypeOf(&dst).Elem()

	if (srcType.Kind() != reflect.Struct && (srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() != reflect.Struct)) ||
		(desType.Kind() != reflect.Struct && (desType.Kind() == reflect.Ptr && desType.Elem().Kind() != reflect.Struct)) {
		return ErrUnsupportedMap
	}

	k := mappingsEntry{SourceType: srcType, DestinationType: desType}
	if _, ok := maps[k]; ok {
		return ErrMapAlreadyExists
	}
	maps[k] = fn

	// if srcType.Kind() == reflect.Ptr && srcType.Elem().Kind() == reflect.Struct {
	//	srcType = srcType.Elem()
	//}
	//
	// if desType.Kind() == reflect.Ptr && desType.Elem().Kind() == reflect.Struct {
	//	desType = desType.Elem()
	//}

	return nil
}

// handleArrayTypes processes array/slice types and returns the base type.
func handleArrayTypes(t reflect.Type) (reflect.Type, bool) {
	isArray := t.Kind() == reflect.Array ||
		(t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Array) ||
		t.Kind() == reflect.Slice ||
		(t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice)

	if isArray {
		return t.Elem(), true
	}

	return t, false
}

// handleCustomMapping processes custom mapping function results.
func handleCustomMapping[TDes any](results []reflect.Value) (TDes, error) {
	if len(results) != 2 {
		return *new(TDes), fmt.Errorf(
			"expected 2 return values from mapping function, got %d",
			len(results),
		)
	}

	if !results[1].IsNil() {
		if err, ok := results[1].Interface().(error); ok && err != nil {
			return *new(TDes), err
		}

		return *new(TDes), errors.New("expected error return value from mapping function")
	}

	if val, ok := results[0].Interface().(TDes); ok {
		return val, nil
	}

	return *new(TDes), errors.New("type assertion failed for mapping result")
}

// Map is a function that maps a source type to a destination type.
func Map[TDes any, TSrc any](src TSrc) (TDes, error) {
	if reflect.ValueOf(src).IsZero() {
		return *new(TDes), nil
	}

	var des TDes
	srcType := reflect.TypeOf(src)
	desType := reflect.TypeOf(des)

	srcType, srcIsArray := handleArrayTypes(srcType)
	desType, desIsArray := handleArrayTypes(desType)

	k := mappingsEntry{SourceType: srcType, DestinationType: desType}
	fn, ok := maps[k]
	if !ok {
		return *new(TDes), ErrMapNotExist
	}

	if fn != nil {
		fnReflect := reflect.ValueOf(fn)
		if desIsArray && srcIsArray {
			linq.From(src).Select(func(x interface{}) interface{} {
				return fnReflect.Call([]reflect.Value{reflect.ValueOf(x)})[0].Interface()
			}).ToSlice(&des)

			return des, nil
		}

		return handleCustomMapping[TDes](fnReflect.Call([]reflect.Value{reflect.ValueOf(src)}))
	}

	desTypeValue := reflect.ValueOf(&des).Elem()
	if err := processValues[TDes, TSrc](reflect.ValueOf(src), desTypeValue); err != nil {
		return *new(TDes), err
	}

	return des, nil
}

// buildProfile builds the mapping profile for source and destination types.
func buildProfile(srcMeta, destMeta typeMeta, srcMethods []string) [][2]string {
	var profile [][2]string

	for srcKey, srcTag := range srcMeta.keysToTags {
		// Check camel case match
		if _, ok := destMeta.keysToTags[strcase.ToCamel(srcKey)]; ok {
			profile = append(profile, [2]string{srcKey, strcase.ToCamel(srcKey)})
		}

		// Check direct key match
		if _, ok := destMeta.keysToTags[srcKey]; ok {
			profile = append(profile, [2]string{srcKey, srcKey})

			continue
		}

		// Check tag to key match
		if destKey, ok := destMeta.tagsToKeys[srcKey]; ok {
			profile = append(profile, [2]string{srcKey, destKey})

			continue
		}

		// Check tag to tag match
		if destKey, ok := destMeta.tagsToKeys[srcTag]; ok {
			profile = append(profile, [2]string{srcKey, destKey})

			continue
		}
	}

	// Add method matches
	for _, method := range srcMethods {
		if _, ok := destMeta.keysToTags[method]; ok {
			profile = append(profile, [2]string{method, method})
		}
	}

	return profile
}

// configProfile is a function that configures the profile.
func configProfile(srcType reflect.Type, destType reflect.Type) {
	if srcType.Kind() != reflect.Struct {
		defaultLogger.GetLogger().Errorf(
			"expected reflect.Struct kind for type %s, but got %s",
			srcType.String(),
			srcType.Kind().String(),
		)
	}

	if destType.Kind() != reflect.Struct {
		defaultLogger.GetLogger().Errorf(
			"expected reflect.Struct kind for type %s, but got %s",
			destType.String(),
			destType.Kind().String(),
		)
	}

	srcMeta := getTypeMeta(srcType)
	destMeta := getTypeMeta(destType)
	srcMethods := getTypeMethods(srcType)

	profile := buildProfile(srcMeta, destMeta, srcMethods)
	profiles[getProfileKey(srcType, destType)] = profile
}

// getProfileKey is a function that gets the profile key.
func getProfileKey(srcType reflect.Type, destType reflect.Type) string {
	return fmt.Sprintf("%s_%s", srcType.Name(), destType.Name())
}

// getTypeMeta is a function that gets the type meta.
func getTypeMeta(val reflect.Type) typeMeta {
	fieldsNum := val.NumField()

	keysToTags := make(map[string]string)
	tagsToKeys := make(map[string]string)

	for i := 0; i < fieldsNum; i++ {
		field := val.Field(i)
		fieldName := field.Name
		fieldTag := field.Tag.Get("mapper")

		keysToTags[fieldName] = fieldTag
		if fieldTag != "" {
			tagsToKeys[fieldTag] = fieldName
		}
	}

	return typeMeta{
		keysToTags: keysToTags,
		tagsToKeys: tagsToKeys,
	}
}

// getTypeMethods is a function that gets the type methods.
func getTypeMethods(val reflect.Type) []string {
	methodsNum := val.NumMethod()
	var keys []string

	for i := 0; i < methodsNum; i++ {
		methodName := val.Method(i).Name
		keys = append(keys, methodName)
	}

	return keys
}

// getSourceFieldValue gets the source field value, handling both exported and unexported fields.
func getSourceFieldValue(
	src reflect.Value,
	sourceField reflect.Value,
	fieldName string,
) reflect.Value {
	if sourceField.Kind() == reflect.Invalid {
		return reflectionHelper.GetFieldValueFromMethodAndReflectValue(
			src.Addr(),
			strcase.ToCamel(fieldName),
		)
	}

	if !sourceField.CanInterface() {
		if mapperConfig.MapUnexportedFields {
			return reflectionHelper.GetFieldValue(sourceField)
		}

		return reflectionHelper.GetFieldValueFromMethodAndReflectValue(
			src.Addr(),
			strcase.ToCamel(fieldName),
		)
	}

	if mapperConfig.MapUnexportedFields {
		return reflectionHelper.GetFieldValue(sourceField)
	}

	return sourceField
}

// mapStructs is a function that maps structs.
func mapStructs[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	profile, ok := profiles[getProfileKey(src.Type(), dest.Type())]

	if !ok {
		// @TODO: Fix this unit tests
		defaultLogger.GetLogger().Warnf(
			"no conversion specified for types %s and %s",
			src.Type().String(),
			dest.Type().String(),
		)

		return
	}

	for _, keys := range profile {
		destinationField := dest.FieldByName(keys[DestKeyIndex])
		sourceField := src.FieldByName(keys[SrcKeyIndex])
		sourceFieldValue := getSourceFieldValue(src, sourceField, keys[SrcKeyIndex])

		if err := processValues[TDes, TSrc](sourceFieldValue, destinationField); err != nil {
			defaultLogger.GetLogger().Errorf("error processing values: %v", err)

			return
		}
	}
}

// mapSlices is a function that maps slices.
func mapSlices[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	// Make dest slice
	dest.Set(reflect.MakeSlice(dest.Type(), src.Len(), src.Cap()))

	// Get each element of slice
	// process values mapping
	for i := 0; i < src.Len(); i++ {
		srcVal := src.Index(i)
		destVal := dest.Index(i)

		if err := processValues[TDes, TSrc](srcVal, destVal); err != nil {
			defaultLogger.GetLogger().Errorf("error processing values: %v", err)

			return
		}
	}
}

// mapPointers is a function that maps pointers.
func mapPointers[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	// create new struct from provided dest type
	val := reflect.New(dest.Type().Elem()).Elem()

	if err := processValues[TDes, TSrc](src.Elem(), val); err != nil {
		defaultLogger.GetLogger().Errorf("error processing values: %v", err)

		return
	}

	// assign address of instantiated struct to destination
	dest.Set(val.Addr())
}

// mapMaps is a function that maps maps.
func mapMaps[TDes any, TSrc any](src reflect.Value, dest reflect.Value) {
	// Make dest map
	dest.Set(reflect.MakeMapWithSize(dest.Type(), src.Len()))

	// Get each element of map as key-values
	// process keys and values mapping and update dest map
	srcMapIter := src.MapRange()
	destMapIter := dest.MapRange()

	for destMapIter.Next() && srcMapIter.Next() {
		destKey := reflect.New(destMapIter.Key().Type()).Elem()
		destValue := reflect.New(destMapIter.Value().Type()).Elem()
		if err := processValues[TDes, TSrc](srcMapIter.Key(), destKey); err != nil {
			defaultLogger.GetLogger().Errorf("error processing values: %v", err)

			return
		}
		if err := processValues[TDes, TSrc](srcMapIter.Value(), destValue); err != nil {
			defaultLogger.GetLogger().Errorf("error processing values: %v", err)

			return
		}

		dest.SetMapIndex(destKey, destValue)
	}
}

// handleInvalidKinds checks and handles invalid kinds.
func handleInvalidKinds(srcKind, destKind reflect.Kind) bool {
	return srcKind == reflect.Invalid || destKind == reflect.Invalid
}

// handleEqualTypes handles the case when source and destination types are equal.
func handleEqualTypes(src, dest reflect.Value) error {
	if src.Type() == dest.Type() {
		reflectionHelper.SetFieldValue(dest, src.Interface())

		return nil
	}

	return nil
}

// handleKindMapping maps values based on their kind.
func handleKindMapping[TDes any, TSrc any](src, dest reflect.Value) error {
	//nolint:exhaustive // missing cases in switch of type reflect.Kind: reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.String, reflect.UnsafePointer
	switch src.Kind() {
	case reflect.Struct:
		mapStructs[TDes, TSrc](src, dest)
	case reflect.Slice:
		mapSlices[TDes, TSrc](src, dest)
	case reflect.Map:
		mapMaps[TDes, TSrc](src, dest)
	case reflect.Ptr:
		mapPointers[TDes, TSrc](src, dest)
	default:
		dest.Set(src)
	}

	return nil
}

// processValues is a function that processes the values.
func processValues[TDes any, TSrc any](src reflect.Value, dest reflect.Value) error {
	if src.Kind() == reflect.Interface {
		src = src.Elem()
	}
	if dest.Kind() == reflect.Interface {
		dest = dest.Elem()
	}

	srcKind := src.Kind()
	destKind := dest.Kind()

	if handleInvalidKinds(srcKind, destKind) {
		return nil
	}

	if srcKind != destKind {
		return nil
	}

	if err := handleEqualTypes(src, dest); err != nil {
		return err
	}

	return handleKindMapping[TDes, TSrc](src, dest)
}
