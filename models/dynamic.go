package models

import "fmt"

func GetDynamicPage(tn, where string, page, size int) (interface{}, int, error) {
	var (
		ndb = GetDB(DYNAMIC_DB)
		d   = GetDriver()
		t   *Table
		err error
	)

	t, err = GetTableByName(tn)
	if err != nil {
		return nil, 0, err
	}

	return d.GetPage(ndb, t, where, page, size)
}

func GetDynamicEntity(tn string, id int64) (interface{}, error) {
	where := fmt.Sprintf("id=%v", id)
	v, _, err := GetDynamicPage(tn, where, 1, 1)
	if err != nil {
		return nil, err
	}
	if result, ok := v.([]map[string]interface{}); ok {
		return result[0], nil
	}
	return nil, fmt.Errorf("have no result.")
}

func DynamicSave(tn string, entity map[string]interface{}) error {
	var (
		err   error
		table *Table
		ndb   = GetDB(DYNAMIC_DB)
		d     = GetDriver()
	)

	table, err = GetTableByName(tn)
	if err != nil {
		return err
	}

	return d.SaveEntity(ndb, table, entity)
}

func DynamicDelete(tn string, ids []int64) error {
	var (
		ndb = GetDB(DYNAMIC_DB)
		d   = GetDriver()
	)

	return d.RemoveEntities(ndb, tn, ids)
}
