package models

func GetDynamicPage(tn string, page, size int) (interface{}, int, error) {
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

	return d.GetPage(ndb, t, page, size)
}
