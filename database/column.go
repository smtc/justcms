package database

type columnType struct {
	Name     string
	SizeMust bool
	Des      string
}

var ColumnTypes = []columnType{
	columnType{"int", false, ""},
	columnType{"bigint", false, ""},
	columnType{"float", false, ""},
	columnType{"double", false, ""},
	columnType{"varchar", true, ""},
	columnType{"text", false, ""},
	columnType{"longtext", false, ""},
	columnType{"bool", false, ""},
	columnType{"datetime", false, ""},
	columnType{"picture", false, ""},
}

var Filters = map[string]string{
	"eq":   "=",
	"neq":  "!=",
	"gt":   ">",
	"egt":  ">=",
	"lt":   "<",
	"elt":  "<=",
	"like": "like",
	"in":   "in",
}
