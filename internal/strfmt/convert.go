package strfmt

import "unsafe"

func ToSnakeCase(s string) string {
	return toScreamingDelimited(s, '_', "", false)
}

func ToPascalCase(s string) string {
	return toCamelInitCase(s, true)
}

func ToCamelCase(s string) string {
	return toCamelInitCase(s, false)
}

func B2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
