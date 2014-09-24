package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/goutils"
)

type mysql struct{}

func (m *mysql) exec(db *gorm.DB, sqlStr string) error {
	if sqlStr == "" {
		return nil
	}
	_, err := db.NewScope(nil).DB().Query(sqlStr)
	if err != nil {
		log.Println(sqlStr, err.Error())
	}
	return err
}

func (m *mysql) HasTable(db *gorm.DB, t *Table) bool {
	scope := db.NewScope(nil)
	suc := gorm.NewDialect("mysql").HasTable(scope, t.Name)
	return suc
}

func (m *mysql) DropTable(db *gorm.DB, t *Table) error {
	sqlStr := fmt.Sprintf("DROP TABLE if exists `%v`;", t.Name)
	return m.exec(db, sqlStr)
}

func (m *mysql) MigrateTable(db *gorm.DB, t, old *Table) error {
	var (
		sqlStr = ""
	)
	if t.Id != 0 && t.Name != old.Name {
		sqlStr = fmt.Sprintf("ALTER TABLE `%v` RENAME `%v`;", old.Name, t.Name)
	}

	return m.exec(db, sqlStr)
}

func (m *mysql) CreateTable(db *gorm.DB, t *Table) error {
	var (
		sqlStr      []string
		primary_key []string
		length      = len(t.Columns)
	)
	if length == 0 {
		log.Println("Table have no columns!")
		return fmt.Errorf("Table have no columns!")
	}

	sqlStr = append(sqlStr, fmt.Sprintf("CREATE table `%v` (", t.Name))
	for _, c := range t.Columns {
		sqlStr = append(sqlStr, m.GetColumn(&c)+",")
		if c.PrimaryKey {
			primary_key = append(primary_key, "`"+c.Name+"`")
		}
	}
	if len(primary_key) == 0 {
		if t.Field("id") == nil {
			sqlStr = append(sqlStr, "`id` BIGINT(20) NOT NULL AUTO_INCREMENT,")
		}
		primary_key = append(primary_key, "`id`")
	}
	sqlStr = append(sqlStr, fmt.Sprintf("PRIMARY KEY (%v)", strings.Join(primary_key, ",")))
	sqlStr = append(sqlStr, ") COLLATE='utf8_general_ci' \nENGINE=MyISAM;")

	return m.exec(db, strings.Join(sqlStr, "\n"))
}

func (m *mysql) GetColumn(c *Column) string {
	var (
		size   int
		sqlStr = ""
		ct     = ColumnTypes[c.Type]
	)
	sqlStr += fmt.Sprintf("`%v` ", c.Name)
	switch c.Type {
	case AUTO_INCREMENT:
		sqlStr += "BIGINT(20) NOT NULL AUTO_INCREMENT"
		return sqlStr
	case BOOL:
		sqlStr += "TINYINT(1) "
	case PICTURE, VARCHAR:
		size = c.Size
		if size == 0 {
			size = ct.Size
		}
		sqlStr += fmt.Sprintf("VARCHAR(%v) ", size)
	case DATE, DATETIME:
		sqlStr += "BIGINT "
	default:
		sqlStr += fmt.Sprintf("%v ", strings.ToUpper(c.Type))
	}
	if c.NotNull {
		sqlStr += "NOT NULL "
	} else {
		sqlStr += "NULL "
	}
	if c.DefaultValue != "" {
		sqlStr += fmt.Sprintf("DEFAULT '%s'", c.DefaultValue)
	}
	return sqlStr
}

func (m *mysql) DropColumn(db *gorm.DB, columns []Column, table string) error {
	cs := []string{}
	for _, c := range columns {
		cs = append(cs, fmt.Sprintf("DROP COLUMN `%s`", c.Name))
	}
	sqlStr := fmt.Sprintf("ALTER TABLE `%s` %s;", table, strings.Join(cs, ","))
	return m.exec(db, sqlStr)
}

func (m *mysql) AddColumn(db *gorm.DB, c *Column, table string) error {
	sqlStr := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", table, m.GetColumn(c))
	return m.exec(db, sqlStr)
}

func (m *mysql) ChangeColumn(db *gorm.DB, c *Column, old, table string) error {
	sqlStr := fmt.Sprintf("ALTER TABLE `%s` CHANGE COLUMN `%s` %s", table, old, m.GetColumn(c))
	return m.exec(db, sqlStr)
}

func (m *mysql) GetPage(db *gorm.DB, t *Table, where string, page, size int) (interface{}, int, error) {
	sqlStr := "SELECT * FROM `%s` %s limit %d, %d;"
	if where != "" {
		where = "where " + where
	}
	sqlStr = fmt.Sprintf(sqlStr, t.Name, where, (page-1)*size, size)

	rows, err := db.NewScope(nil).DB().Query(sqlStr)
	defer rows.Close()
	if err != nil {
		return nil, 0, err
	}

	var (
		columns, _ = rows.Columns()
		values     = make([]sql.RawBytes, len(columns))
		scanArgs   = make([]interface{}, len(values))
		results    = []map[string]interface{}{}
	)

	var getValue = func(c *Column, s string) interface{} {
		switch c.Type {
		case INT, BIGINT, AUTO_INCREMENT:
			return goutils.ToInt64(s, 0)
		case BOOL:
			return goutils.ToBool(s, false)
		case FLOAT, DOUBLE:
			return goutils.ToFloat64(s, 0)
		case DATE, DATETIME:
			t := time.Unix(goutils.ToInt64(s, 0), 0)
			if c.Type == DATE {
				return t.Format("2006-01-02")
			} else {
				return t.Format("2006-01-02 15:04:05")
			}
		}
		return s
	}

	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		result := map[string]interface{}{}
		err = rows.Scan(scanArgs...)
		if err != nil {
			//panic(err.Error()) // proper error handling instead of panic in your app
		}
		for i, v := range values {
			result[columns[i]] = getValue(t.Field(columns[i]), string(v))
		}
		results = append(results, result)
	}
	return results, 0, nil
}

func (m *mysql) SaveEntity(db *gorm.DB, t *Table, entity map[string]interface{}) error {
	var (
		id     int64
		keys   = []string{}
		vals   = []string{}
		sqlStr = ""
		getter = goutils.Getter(entity)
	)

	for _, c := range t.Columns {
		if (c.Name == "id" || entity[c.Name] == nil) && c.Type != BOOL {
			continue
		}

		keys = append(keys, "`"+c.Name+"`")
		switch c.Type {
		case DATE, DATETIME:
			var ft string
			if c.Type == DATE {
				ft = "2006-01-02"
			} else {
				ft = goutils.TIMEFORMAT
			}
			d := getter.GetTime(c.Name, goutils.TIMEDEFAULT, ft)
			vals = append(vals, fmt.Sprintf("%d", d.Unix()))
		case BOOL:
			b := getter.GetBool(c.Name, false)
			vals = append(vals, fmt.Sprintf("%v", b))
		default:
			vals = append(vals, "'"+getter.GetString(c.Name, "")+"'")
		}
	}

	id = getter.GetInt64("id", 0)
	if id == 0 {
		// insert
		sqlStr = fmt.Sprintf("INSERT INTO `%s` (%s) VALUE (%s);", t.Name, strings.Join(keys, ","), strings.Join(vals, ","))
	} else {
		// update
		sets := []string{}
		for i, k := range keys {
			sets = append(sets, k+"="+vals[i])
		}

		sqlStr = fmt.Sprintf("UPDATE `%s` SET %s WHERE `id`=%d;", t.Name, strings.Join(sets, ","), id)
	}

	return m.exec(db, sqlStr)
}

func (m *mysql) RemoveEntities(db *gorm.DB, tn string, ids []int64) error {
	sqlStr := fmt.Sprintf("DELETE FROM `%s` WHERE `id` in (%v);", tn, goutils.ToString(ids, ""))
	return m.exec(db, sqlStr)
}
