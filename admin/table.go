package admin

import (
	"log"
	"net/http"

	"github.com/smtc/goutils"
	"github.com/smtc/justcms/models"
	"github.com/zenazn/goji/web"
)

func TableList(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	var tbls []models.Table
	if err := models.TableList(&tbls); err != nil {
		h.RenderError(err.Error())
		return
	}
	list, _ := goutils.ToMapList(tbls, []string{}, goutils.FilterModeExclude)
	h.RenderPage(list)
}

func TableEntity(c web.C, w http.ResponseWriter, r *http.Request) {
	h := goutils.HttpHandler(c, w, r)
	id := h.Param.GetInt64("id", 0)
	if id == 0 {
		h.RenderJson(nil, 0)
		return
	}

	var t models.Table
	if err := t.Get(id); err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(t, 1)
}

func TableSave(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		t   models.Table
		err error
		h   = goutils.HttpHandler(c, w, r)
	)

	//if err = json.NewDecoder(r.Body).Decode(&t); err != nil {
	if err = h.FormatBody(&t); err != nil {
		log.Println(err)
		h.RenderError(err.Error())
		return
	}

	err = t.Save()
	if err != nil {
		log.Println(err)
		h.RenderError(err.Error())
		return

	}

	h.RenderJson(t, 1)
}

func TableDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h      = goutils.HttpHandler(c, w, r)
		db     = models.GetDB(models.DEFAULT_DB)
		ids    []int64
		err    error
		tables []models.Table
	)

	err = h.FormatBody(&ids)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	err = db.Where("id in (?)", ids).Find(&tables).Error
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	for _, t := range tables {
		t.Delete()
	}

	h.RenderJson(nil, 1)
}
