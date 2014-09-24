package models

import (
	"fmt"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/goutils"
)

type Table struct {
	Id        int64
	Name      string    `sql:"size:45;not null;unique"`
	Alias     string    `sql:"size:45"`
	Des       string    `sql:"size:512"`
	CreatedAt time.Time `json:"created_at"`
	EditAt    time.Time `json:"edit_at"`
	Engine    string    `sql:"size:45"`
	Columns   []Column
}

var (
	// reserve table name
	NotAllowTables = []string{}
	tables         = make(map[int64]*Table)
)

func getTableDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func (t *Table) Exist() bool {
	for _, nat := range NotAllowTables {
		if t.Name == nat {
			return true
		}
	}
	db := getTableDB()
	var count int
	db.Model(Table{}).Where("id!=? and name=?", t.Id, t.Name).Count(&count)
	return count > 0
}

func GetTable(id int64) (*Table, error) {
	var (
		db  = getTableDB()
		t   = tables[id]
		err error
	)
	if t != nil {
		return t, nil
	}
	t = &Table{}
	err = db.First(t, id).Error
	if err != nil {
		return t, err
	}
	tables[id] = t
	t.Refresh()
	return t, nil
}

func GetTableByName(name string) (*Table, error) {
	for _, t := range tables {
		if t.Name == name {
			return t, nil
		}
	}

	var (
		t   = &Table{}
		db  = getTableDB()
		err error
	)

	t = &Table{}
	err = db.Where("name=?", name).First(t).Error
	if err != nil {
		return t, err
	}
	tables[t.Id] = t
	t.Refresh()
	return t, nil
}

func (t *Table) Get(id int64) error {
	db := getTableDB()
	err := db.First(t, id).Error
	return err
}

func (t *Table) Refresh() error {
	db := getTableDB()
	err := db.First(t, t.Id).Error
	if err != nil {
		return err
	}
	t.Columns = nil
	err = db.Model(t).Related(&t.Columns).Order("order_index").Error
	sort.Sort(ColumnSort(t.Columns))
	return err
}

func (t *Table) Save() error {
	var (
		db     = getTableDB()
		isNew  = false
		d      = GetDriver()
		column Column
		old    Table
	)
	if t.Exist() {
		return fmt.Errorf("table '%v' is existed", t.Name)
	}
	if t.Id == 0 {
		t.CreatedAt = time.Now()
		isNew = true
	} else {
		db.Where("id = ?", t.Id).Find(&old)
		t.CreatedAt = old.CreatedAt
	}
	t.EditAt = time.Now()

	if err := db.Save(t).Error; err != nil {
		return err
	}

	if isNew {
		if t.Field("id") == nil {
			column.Name = "id"
			column.Alias = "id"
			column.Type = AUTO_INCREMENT
			column.PrimaryKey = true
			column.TableId = t.Id
			column.NotNull = true
			column.Save()
		}

		t.CreateTable()
	} else {
		d.MigrateTable(db, t, &old)
	}
	tables[t.Id] = t
	return nil
}

func (t *Table) Delete() error {
	var (
		db  = getTableDB()
		d   = GetDriver()
		err error
	)
	if err = db.Where("table_id = ?", t.Id).Delete(Column{}).Error; err != nil {
		return err
	}

	tables[t.Id] = nil
	err = db.Delete(t).Error
	if err != nil {
		return err
	}

	d.DropTable(GetDB(DYNAMIC_DB), t)
	return nil
}

func TableList(tbls *[]Table) error {
	db := getTableDB()
	err := db.Find(tbls).Error
	return err
}

func (t *Table) Field(v interface{}) *Column {
	for _, c := range t.Columns {
		if id, ok := v.(int64); ok && c.Id == id {
			return &c
		}
		if name, ok := v.(string); ok && c.Name == name {
			return &c
		}
	}
	return nil
}

func (t *Table) PrimaryKey() *[]Column {
	keys := []Column{}
	for _, c := range t.Columns {
		if c.PrimaryKey {
			keys = append(keys, c)
		}
	}
	return &keys
}

func (t Table) MarshalJSON() ([]byte, error) {
	j, _ := goutils.ToJsonOnly(t)
	return []byte(j), nil
}

func (t *Table) CreateTable() error {
	var (
		d   = GetDriver()
		db  = GetDB(DYNAMIC_DB)
		err error
	)

	if d.HasTable(db, t) {
		return fmt.Errorf("Table '%v' already exists.", t.Name)
	}

	t.Refresh()
	err = d.CreateTable(db, t)
	return err
}
