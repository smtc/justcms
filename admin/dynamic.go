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
				"del":  "id={{id}}",
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
		total int
		err   error
		p     interface{}
	)

	p, total, err = models.GetDynamicPage(tn, 1, 20)
	if err != nil {
		h.RenderError(err.Error())
		return
	}

	_ = total
	h.RenderPage(p)
}
