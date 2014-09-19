package models

import (
	"fmt"
	"strings"
)

type mysql struct{}

func (m *mysql) ExistTable(t *Table) bool {
	return false
}

func (m *mysql) CreateTable(t *Table) error {
	var (
		sql         []string
		primary_key []string
		length      = len(t.Columns)
	)
	if length == 0 {
		return fmt.Errorf("Table have no columns!")
	}

	sql = append(sql, fmt.Sprintf("CREATE table `%v` (", t.Name))
	for _, c := range t.Columns {
		sql = append(sql, c.GetCreate()+",")
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

	println(strings.Join(sql, "\n"))

	return nil
}

func (m *mysql) MigrateTable(t *Table) error {
	return nil
}
