package admin

import (
	"net/http"

	"github.com/smtc/justcms/models"
	"github.com/smtc/justcms/utils"
	"github.com/zenazn/goji/web"
)

func ColumnList(c web.C, w http.ResponseWriter, r *http.Request) {
	param := utils.Getter(c.URLParams)
	id := param.GetInt64("table", 0)
	models, _ := models.ColumnList(id)
	list, _ := utils.ToMapList(models, []string{}, utils.FilterModeExclude)
	utils.Render(w).RenderPage(list)
}
