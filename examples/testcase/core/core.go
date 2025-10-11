package core

import (
	"reflect"
	"time"
)

type HouseUnitType uint8

const (
	HouseUnitTypeA HouseUnitType = iota + 22
	HouseUnitTypeB
	HouseUnitTypeC
	HouseUnitTypeD
	HouseUnitTypeE
)

// type Condition string

// const (
// 	Old        Condition = "OLD"
// 	New        Condition = "NEW"
// 	Ok, Normal Condition = "OK", "OK2"
// )

// type LogLevel string

// var (
// 	Debug LogLevel = "debug"
// 	Info  LogLevel = "info"
// )

type User struct {
	ID         int64 `sql:",pk,auto_increment"`
	No         uint
	JoinedTime time.Time
	Address    Address
	Kind       reflect.Kind
	Type       HouseUnitType
	Chan       reflect.ChanDir
	PostalCode *string
	ExtraInfo  struct {
		Flag bool
	}
	Nicknames [2]string
	Slice     []float64
	Map       map[string]float64
}

// +sql:ignore
type Address struct {
	Line1       string
	Line2       string
	CountryCode string
}

func main() {
	println("hello !")
}
