package postgres

import "github.com/si3nloong/sqlgen/codegen/dialect"

func (s *postgresDriver) ColumnDataTypes() map[string]*dialect.ColumnType {
	return map[string]*dialect.ColumnType{
		// "rune": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "CHAR(1)",
		// 		Encoder:  "string({{goPath}})",
		// 	}
		// },
		// "[]byte": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "BLOB",
		// 	}
		// },
		// "string": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "VARCHAR(255)",
		// 	}
		// },
		// "[]string": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "text[]",
		// 		Encoder:  "github.com/lib/pq.StringArray({{goPath}})",
		// 		Decoder:  "github.com/lib/pq.StringArray({{addrOfGoPath}})",
		// 	}
		// },
		// "[...]string": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "text[{{len}}]",
		// 		Encoder:  "github.com/lib/pq.StringArray({{goPath}})",
		// 		Decoder:  "github.com/lib/pq.StringArray({{addrOfGoPath}})",
		// 	}
		// },
		// "[]bool": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "bool[]",
		// 		Encoder:  "github.com/lib/pq.BoolArray({{goPath}})",
		// 		Decoder:  "github.com/lib/pq.BoolArray({{addrOfGoPath}})",
		// 	}
		// },
		// "[][]byte": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "bytea",
		// 		Encoder:  "github.com/lib/pq.ByteaArray({{goPath}})",
		// 		Decoder:  "github.com/lib/pq.ByteaArray({{addrOfGoPath}})",
		// 	}
		// },
		// "[]float32": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "double[]",
		// 		Encoder:  "github.com/lib/pq.Float32Array({{goPath}})",
		// 		Decoder:  "github.com/lib/pq.Float32Array({{addrOfGoPath}})",
		// 	}
		// },
		// "[]float64": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "double[]",
		// 		Encoder:  "github.com/lib/pq.Float64Array({{goPath}})",
		// 		Decoder:  "github.com/lib/pq.Float64Array({{addrOfGoPath}})",
		// 	}
		// },
		// "[]int32": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "int4[]",
		// 		Encoder:  "github.com/lib/pq.Int32Array({{goPath}})",
		// 		Decoder:  "github.com/lib/pq.Int32Array({{addrOfGoPath}})",
		// 	}
		// },
		// "[]int64": func(c dialect.Column) *dialect.Mapping {
		// 	return &dialect.Mapping{
		// 		DataType: "int8[]",
		// 		Encoder:  "github.com/lib/pq.Int64Array({{goPath}})",
		// 		Decoder:  "github.com/lib/pq.Int64Array({{addrOfGoPath}})",
		// 	}
		// },
	}
}

// func dataType(f sequel.GoColumnSchema) *columnDefinition {
// 	var (
// 		ptrs = make([]types.Type, 0)
// 		t    = f.Type()
// 		col  = new(columnDefinition)
// 		prev types.Type
// 	)

// 	col.name = f.ColumnName()
// 	if dataType, ok := f.DataType(); ok {
// 		col.dataType = dataType
// 		return col
// 	}
// 	defer func() {
// 		if !col.nullable {
// 			col.nullable = len(ptrs) > 0
// 		}
// 	}()

// 	for t != nil {
// 		switch v := t.(type) {
// 		case *types.Pointer:
// 			prev = v.Elem()
// 			ptrs = append(ptrs, v)
// 		case *types.Basic:
// 			prev = t.Underlying()
// 		case *types.Named:
// 			prev = t.Underlying()
// 		default:
// 			break
// 		}

// 		switch t.String() {
// 		case "rune":
// 			col.dataType = "CHAR(1)"
// 			col.length = 1
// 			return col
// 		case "int8", "int16":
// 			col.dataType = serialOrInt(f, "INT2")
// 			col.defaultValue = int64(0)
// 			return col
// 		case "int32", "int":
// 			col.dataType = serialOrInt(f, "INT4") // Or INT4
// 			col.defaultValue = int64(0)
// 			return col
// 		case "int64":
// 			col.dataType = serialOrInt(f, "INT8")
// 			col.defaultValue = int64(0)
// 			return col
// 		case "bool":
// 			col.dataType = "BOOL"
// 			col.defaultValue = false
// 			return col
// 		case "uint8", "uint16", "byte":
// 			col.dataType = serialOrInt(f, "INT2")
// 			col.defaultValue = uint64(0)
// 			col.check = sql.RawBytes(`CHECK(` + f.ColumnName() + ` >= 0)`)
// 			return col
// 		case "uint32", "uint":
// 			col.dataType = serialOrInt(f, "INT4")
// 			col.defaultValue = uint64(0)
// 			col.check = sql.RawBytes(`CHECK(` + f.ColumnName() + ` >= 0)`)
// 			return col
// 		case "uint64":
// 			col.dataType = serialOrInt(f, "INT8")
// 			col.defaultValue = uint64(0)
// 			col.check = sql.RawBytes(`CHECK(` + f.ColumnName() + ` >= 0)`)
// 			return col
// 		case "float32":
// 			col.dataType = "DOUBLE PRECISION"
// 			col.defaultValue = float64(0.0)
// 			return col
// 		case "float64":
// 			col.dataType = "DOUBLE PRECISION"
// 			col.defaultValue = float64(0.0)
// 			return col
// 		case "cloud.google.com/go/civil.Time":
// 			col.dataType = "TIME"
// 			col.defaultValue = sql.RawBytes(`CURRENT_TIME`)
// 			return col
// 		case "cloud.google.com/go/civil.Date":
// 			col.dataType = "DATE"
// 			col.defaultValue = sql.RawBytes(`CURRENT_DATE`)
// 			return col
// 		case "cloud.google.com/go/civil.DateTime":
// 			col.defaultValue = sql.RawBytes(`NOW()`)
// 			if size := f.Size(); size > 0 {
// 				col.dataType = fmt.Sprintf("TIMESTAMP(%d)", size)
// 				col.length = size
// 			} else {
// 				col.dataType = "TIMESTAMP"
// 			}
// 			return col
// 		case "time.Time":
// 			col.defaultValue = sql.RawBytes(`NOW()`)
// 			if size := f.Size(); size > 0 {
// 				col.length = size
// 				col.dataType = fmt.Sprintf("TIMESTAMP(%d) WITH TIME ZONE", size)
// 			} else {
// 				col.dataType = "TIMESTAMP WITH TIME ZONE"
// 			}
// 			return col
// 		case "string":
// 			col.defaultValue = ""
// 			col.length = 255
// 			if size := f.Size(); size > 0 {
// 				col.length = size
// 			}
// 			col.dataType = fmt.Sprintf("VARCHAR(%d)", col.length)
// 			return col
// 		case "[]rune":
// 			col.dataType = "VARCHAR(255)"
// 			return col
// 		case "[]byte":
// 			col.dataType = "BYTEA"
// 			return col
// 		case "[16]byte":
// 			// if f.IsBinary {
// 			// 	return "BIT(16)"
// 			// }
// 			col.dataType = "VARBIT(36)"
// 			return col
// 		case "encoding/json.RawMessage":
// 			col.dataType = "JSONB"
// 			return col
// 		case "database/sql.NullBool":
// 			col.dataType = "BOOL"
// 			col.defaultValue = false
// 			col.nullable = true
// 			return col
// 		case "database/sql.NullString":
// 			col.dataType = "VARCHAR(255)"
// 			col.defaultValue = ""
// 			col.length = 255
// 			col.nullable = true
// 			return col
// 		case "database/sql.NullInt16":
// 			col.dataType = serialOrInt(f, "INT2")
// 			col.defaultValue = int64(0)
// 			col.nullable = true
// 			return col
// 		case "database/sql.NullInt32":
// 			col.dataType = serialOrInt(f, "INT4")
// 			col.defaultValue = int64(0)
// 			col.nullable = true
// 			return col
// 		case "database/sql.NullInt64":
// 			col.dataType = serialOrInt(f, "INT8")
// 			col.defaultValue = int64(0)
// 			col.nullable = true
// 			return col
// 		case "database/sql.NullFloat64":
// 			col.dataType = "DOUBLE PRECISION"
// 			col.defaultValue = float64(0.0)
// 			col.nullable = true
// 			return col
// 		case "database/sql.NullTime":
// 			col.dataType = "TIMESTAMP WITH TIMEZONE"
// 			col.nullable = true
// 			col.defaultValue = sql.RawBytes(`NOW()`)
// 			return col
// 		default:
// 			switch {
// 			case strings.HasPrefix(t.String(), "[]"):
// 				col.dataType = "JSON"
// 				return col
// 			case strings.HasPrefix(t.String(), "map"):
// 				col.dataType = "JSON"
// 				return col
// 			}
// 		}
// 		if prev == t {
// 			break
// 		}
// 		t = prev
// 	}
// 	col.dataType = "VARCHAR(255)"
// 	return col
// }

// func serialOrInt(f sequel.GoColumnSchema, dataType string) string {
// 	if f.AutoIncr() {
// 		return strings.ReplaceAll(dataType, "INT", "SERIAL")
// 	}
// 	return dataType
// }

// func format(v any) string {
// 	switch vi := v.(type) {
// 	case string:
// 		return "'" + vi + "'"
// 	case bool:
// 		return strconv.FormatBool(vi)
// 	case int:
// 		return strconv.Itoa(vi)
// 	case int64:
// 		return strconv.FormatInt(vi, 10)
// 	case uint64:
// 		return strconv.FormatUint(vi, 10)
// 	case float32:
// 		return strconv.FormatFloat(float64(vi), 'f', -1, 64)
// 	case float64:
// 		return strconv.FormatFloat(vi, 'f', -1, 64)
// 	case sql.RawBytes:
// 		return unsafe.String(unsafe.SliceData(vi), len(vi))
// 	default:
// 		panic(fmt.Sprintf("unsupported data type %T", vi))
// 	}
// }
