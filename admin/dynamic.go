package admin

import (
	"fmt"
	"net/http"

	"github.com/smtc/goutils"
	"github.com/smtc/justcms/models"
	"github.com/zenazn/goji/web"
)

func DynamicStruct(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h      = goutils.HttpHandler(c, w, r)
		tn     = h.Param.GetString("table", "")
		method = h.Param.GetString("method", "list")
		table  *models.Table
		err    error
	)

	if tn == "" {
		h.RenderError("数据表不存在")
		return
	}

	table, err = models.GetTableByName(tn)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	srt := models.Struct{}
	srt.GetStruct(table)
	srt.Src = fmt.Sprintf("dynamic/api/%s/", tn)
	if method == "list" {
		srt.Op = map[string]map[string]string{
			"unit": map[string]string{
				"edit": fmt.Sprintf("dynamic.edit:%s#{{id}}", tn),
				"del":  "[{{id}}]",
			},
			"mult": map[string]string{
				"del":     "id",
				"new":     fmt.Sprintf("dynamic.edit:%s#0", tn),
				"refresh": "",
			},
		}
	}

	h.RenderJsonNoWrap(srt)
}

func DynamicList(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h     = goutils.HttpHandler(c, w, r)
		tn    = h.Param.GetString("table", "")
		where = ""
		total int
		err   error
		p     interface{}
	)

	p, total, err = models.GetDynamicPage(tn, where, 1, 20)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	_ = total
	h.RenderPage(p)
}

func DynamicEntity(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h      = goutils.HttpHandler(c, w, r)
		tn     = h.Param.GetString("table", "")
		id     = h.Param.GetInt64("id", 0)
		entity interface{}
		err    error
	)

	if id == 0 {
		h.RenderJson(nil, 1)
		return
	}

	entity, err = models.GetDynamicEntity(tn, id)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(entity, 1)
}

func DynamicSave(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h      = goutils.HttpHandler(c, w, r)
		tn     = h.Param.GetString("table", "")
		entity map[string]interface{}
		err    error
	)

	err = h.FormatBody(&entity)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	err = models.DynamicSave(tn, entity)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(nil, 1)
}

func DynamicDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		h   = goutils.HttpHandler(c, w, r)
		tn  = h.Param.GetString("table", "")
		ids []int64
		err error
	)

	err = h.FormatBody(&ids)
	if err != nil {
		h.RenderError(err.Error())
	}

	err = models.DynamicDelete(tn, ids)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	h.RenderJson(nil, 1)
}
