package admin

import (
	"fmt"
	"net/http"

	"github.com/smtc/goutils"
	"github.com/smtc/justcms/models"
	"github.com/zenazn/goji/web"
)

func ColumnList(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)

	var (
		id    = h.Param.GetInt64("table_id", 0)
		err   error
		table *models.Table
	)

	table, err = models.GetTable(id)
	if err != nil {
		h.RenderError(err.Error())
	}
	table.Refresh()

	fmt.Printf("%v", table.Columns)

	list, _ := goutils.ToMapList(table.Columns, []string{}, goutils.FilterModeExclude)
	h.RenderPage(list, 0)
}

func ColumnEntity(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)

	var (
		column = models.Column{TableId: h.Param.GetInt64("table_id", 0)}
		id     = h.Param.GetInt64("id", 0)
	)

	if id == 0 {
		h.RenderJson(nil, 0, "")
	} else {
		column.Get(id)
		h.RenderJson(column, 1, "")
	}
}

func ColumnSave(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)

	var (
		column   models.Column
		err      error
		table_id = h.Param.GetInt64("table_id", 0)
	)

	if err = h.FormatBody(&column); err != nil {
		h.RenderError(err.Error())
		return
	}

	if column.Id == 0 {
		column.TableId = table_id
	}

	if err = column.Save(); err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(nil, 1, "")
}

func ColumnDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)

	var data []int64
	if err := h.FormatBody(&data); err != nil {
		h.RenderError(err.Error())
		return
	}

	if err := models.ColumnDelete("id in (?)", data); err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(data, 1, "")
}

func ColumnType(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)

	var types []models.Select
	for _, t := range models.ColumnTypes {
		types = append(types, models.Select{t.Name, t.Name})
	}

	h.RenderJson(types, 1, "")
}

func ColumnFilter(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	filters := map[string]string{}
	for k, v := range models.Filters {
		filters[v] = k
	}
	h.RenderJson(filters, 1, "")
}
