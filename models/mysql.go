package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
)

type mysql struct{}

func (m *mysql) HasTable(db *gorm.DB, t *Table) bool {
	scope := db.NewScope(nil)
	suc := gorm.NewDialect("mysql").HasTable(scope, t.Name)
	return suc
}

func (m *mysql) DropTable(db *gorm.DB, t *Table) {
	sql := fmt.Sprintf("DROP TABLE if exists `%v`;", t.Name)
	db.NewScope(nil).Raw(sql).Exec()
	log.Println(sql)
}

func (m *mysql) MigrateTable(db *gorm.DB, t, old *Table) {
	var (
		sql = ""
	)
	if t.Id != 0 && t.Name != old.Name {
		sql = fmt.Sprintf("ALTER TABLE `%v` RENAME `%v`;", old.Name, t.Name)
	}
	if sql != "" {
		db.NewScope(nil).Raw(sql).Exec()
		log.Println(sql)
	}
}

func (m *mysql) CreateTable(db *gorm.DB, t *Table) error {
	var (
		sql         []string
		primary_key []string
		err         error
		length      = len(t.Columns)
	)
	if length == 0 {
		return fmt.Errorf("Table have no columns!")
	}

	sql = append(sql, fmt.Sprintf("CREATE table `%v` (", t.Name))
	for _, c := range t.Columns {
		sql = append(sql, m.GetColumn(&c)+",")
		if c.PrimaryKey {
			primary_key = append(primary_key, "`"+c.Name+"`")
		}
	}
	if len(primary_key) == 0 {
		if t.Field("id") == nil {
			sql = append(sql, "`id` BIGINT(20) NOT NULL AUTO_INCREMENT,")
		}
		primary_key = append(primary_key, "`id`")
	}
	sql = append(sql, fmt.Sprintf("PRIMARY KEY (%v)", strings.Join(primary_key, ",")))
	sql = append(sql, ") COLLATE='utf8_general_ci' \nENGINE=MyISAM;")

	log.Println(strings.Join(sql, "\n"))
	_, err = db.NewScope(nil).DB().Query(strings.Join(sql, "\n"))

	return err
}

func (m *mysql) GetColumn(c *Column) string {
	var (
		size int
		sql  = ""
		ct   = ColumnTypes[c.Type]
	)
	sql += fmt.Sprintf("`%v` ", c.Name)
	switch c.Type {
	case AUTO_INCREMENT:
		sql += "BIGINT(20) NOT NULL AUTO_INCREMENT"
		return sql
	case BOOL:
		sql += "TINYINT(1) "
	case PICTURE, VARCHAR:
		size = c.Size
		if size == 0 {
			size = ct.Size
		}
		sql += fmt.Sprintf("VARCHAR(%v) ", size)
	default:
		sql += fmt.Sprintf("%v ", strings.ToUpper(c.Type))
	}
	if c.NotNull {
		sql += "NOT NULL "
	} else {
		sql += "NULL DEFAULT NULL "
	}
	return sql
}
