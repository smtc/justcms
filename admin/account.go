package admin

import (
	"net/http"

	"github.com/smtc/JustCms/models"
	"github.com/smtc/JustCms/utils"
)

func AccountList(w http.ResponseWriter, r *http.Request) {
	models, _ := models.AccountList(0, 20, map[string]interface{}{})
	list, _ := utils.ToMapList(models, []string{"email", "name"}, utils.FilterModeInclude)
	utils.RenderPage(list, w, r)
}
