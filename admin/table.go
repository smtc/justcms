package admin

import (
	"log"
	"net/http"

	"github.com/smtc/justcms/models"
	"github.com/smtc/justcms/utils"
	"github.com/zenazn/goji/web"
)

func TableList(w http.ResponseWriter, r *http.Request) {
	models, _ := models.TableList()
	list, _ := utils.ToMapList(models, []string{}, utils.FilterModeExclude)
	utils.RenderPage(list, w, r)
}

func TableEntity(c web.C, w http.ResponseWriter, r *http.Request) {
	id := utils.GetInt64(c.URLParams, "id", 0)
	if id == 0 {
		utils.RenderJson(nil, 0, w)
		return
	}

	var t models.Table
	t.Get(id)

	utils.RenderJson(t, 1, w)
}

func TableSave(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		t   models.Table
		err error
	)

	//if err = json.NewDecoder(r.Body).Decode(&t); err != nil {
	if err = utils.RequestStruct(r, &t); err != nil {
		log.Println(err)
		utils.RenderError(err.Error(), w)
		return
	}

	err = t.Save()
	if err != nil {
		log.Println(err)
		utils.RenderError(err.Error(), w)
		return

	}

	utils.RenderJson(t, 1, w)
}
