package main

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

type HouseUnit struct {
	No        uint
	BuildTime time.Time
	Address   Address
	Kind      reflect.Kind
	Type      HouseUnitType
	Chan      reflect.ChanDir
	// Lvl       LogLevel
	// Condition Condition
	Inner struct {
		Flag bool
	}
	Arr   [2]string
	Slice []float64
	Map   map[string]float64
}

type Address struct {
	Line1       string
	Line2       string
	CountryCode string
}

func main() {

}
