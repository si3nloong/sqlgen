package encoding

import (
	"database/sql"
	"testing"
)

func TestUintScannerConversions(t *testing.T) {
	checkScanner(t, func(v *uint) sql.Scanner { return UintScanner[uint](v) }, func(v **uint) sql.Scanner { return UintScanner[uint](v) },
		[]scanCase[uint]{{int64(42), 42}, {[]byte("17"), 17}}, []any{"42", true, []byte("invalid"), []byte("18446744073709551616")})
}

func TestUint8ScannerConversions(t *testing.T) {
	checkScanner(t, func(v *uint8) sql.Scanner { return Uint8Scanner[uint8](v) }, func(v **uint8) sql.Scanner { return Uint8Scanner[uint8](v) },
		[]scanCase[uint8]{{int64(42), 42}, {[]byte("17"), 17}}, []any{"42", true, []byte("invalid"), []byte("18446744073709551616")})
}

func TestUint16ScannerConversions(t *testing.T) {
	checkScanner(t, func(v *uint16) sql.Scanner { return Uint16Scanner[uint16](v) }, func(v **uint16) sql.Scanner { return Uint16Scanner[uint16](v) },
		[]scanCase[uint16]{{int64(42), 42}, {[]byte("17"), 17}}, []any{"42", true, []byte("invalid"), []byte("18446744073709551616")})
}

func TestUint32ScannerConversions(t *testing.T) {
	checkScanner(t, func(v *uint32) sql.Scanner { return Uint32Scanner[uint32](v) }, func(v **uint32) sql.Scanner { return Uint32Scanner[uint32](v) },
		[]scanCase[uint32]{{int64(42), 42}, {[]byte("17"), 17}}, []any{"42", true, []byte("invalid"), []byte("18446744073709551616")})
}

func TestUint64ScannerConversions(t *testing.T) {
	checkScanner(t, func(v *uint64) sql.Scanner { return Uint64Scanner[uint64](v) }, func(v **uint64) sql.Scanner { return Uint64Scanner[uint64](v) },
		[]scanCase[uint64]{{int64(42), 42}, {[]byte("17"), 17}}, []any{"42", true, []byte("invalid"), []byte("18446744073709551616")})
}
