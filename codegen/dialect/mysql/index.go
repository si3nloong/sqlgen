package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"unsafe"
)

//go:generate stringer --type indexType --linecomment
type indexType uint8

const (
	bTree  indexType = iota // BTREE
	unique                  // UNIQUE
	pk
)

func indexName(columns []string, idxType indexType) string {
	str := strings.Join(columns, ",")
	hash := md5.Sum(unsafe.Slice(unsafe.StringData(str), len(str)))
	prefix := "IX"
	switch idxType {
	case unique:
		prefix = "UQ"
	case pk:
		prefix = "PK"
	}
	return prefix + "_" + hex.EncodeToString(hash[:])
}
