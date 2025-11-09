package goutil

import (
	"go/types"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var genericRegexp = regexp.MustCompile(`(.*)\[\w+\]$`)

func GetTypeName(i any, args ...string) string {
	typeOf := reflect.TypeOf(i)
	submatches := genericRegexp.FindStringSubmatch(typeOf.Name())
	return typeOf.PkgPath() + "." + submatches[1]
}

func GetFunctionName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GenericFunc(i any, args ...string) string {
	str := strings.TrimSuffix(GetFunctionName(i), "[...]")
	return str + "(" + strings.Join(args, ",") + ")"
}

func GenericFuncName(i any, genericType string, args ...string) string {
	str := GetFunctionName(i)
	if strings.HasSuffix(str, "[...]") && genericType != "" {
		str = str[:len(str)-5] + "[" + genericType + "]"
	}
	str = str + "(" + strings.Join(args, ",") + ")"
	return str
}

func PointerUnderlyingType(t types.Type) (ut types.Type) {
	v, ok := t.(*types.Pointer)
	if ok {
		t = PointerUnderlyingType(v.Elem())
		return t
	}
	return t
}
