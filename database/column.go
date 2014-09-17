package database

const (
	AUTO_INCREMENT = "auto_increment"
	INT            = "int"
	BIGINT         = "bigint"
	FLOAT          = "float"
	DOUBLE         = "double"
	VARCHAR        = "varchar"
	TEXT           = "text"
	LONGTEXT       = "longtext"
	BOOL           = "bool"
	DATETIME       = "datetime"
	PICTURE        = "picture"
)

type columnType struct {
	Name     string
	SizeMust bool
	Des      string
}

var ColumnTypes = map[string]columnType{
	AUTO_INCREMENT: columnType{AUTO_INCREMENT, false, ""},
	INT:            columnType{INT, false, ""},
	BIGINT:         columnType{BIGINT, false, ""},
	FLOAT:          columnType{FLOAT, false, ""},
	DOUBLE:         columnType{DOUBLE, false, ""},
	VARCHAR:        columnType{VARCHAR, true, ""},
	TEXT:           columnType{TEXT, false, ""},
	LONGTEXT:       columnType{LONGTEXT, false, ""},
	BOOL:           columnType{BOOL, false, ""},
	DATETIME:       columnType{DATETIME, false, ""},
	PICTURE:        columnType{PICTURE, false, ""},
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
