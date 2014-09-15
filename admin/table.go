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
	utils.Render(w).RenderPage(list)
}

func TableEntity(c web.C, w http.ResponseWriter, r *http.Request) {
	render := utils.Render(w)
	param := utils.Getter(c.URLParams)
	id := param.GetInt64("id", 0)
	if id == 0 {
		render.RenderJson(nil, 0)
		return
	}

	var t models.Table
	if err := t.Get(id); err != nil {
		render.RenderError(err.Error())
		return
	}

	render.RenderJson(t, 1)
}

func TableSave(c web.C, w http.ResponseWriter, r *http.Request) {
	render := utils.Render(w)
	var (
		t   models.Table
		err error
	)

	//if err = json.NewDecoder(r.Body).Decode(&t); err != nil {
	if err = utils.Request(r).FormatBody(&t); err != nil {
		log.Println(err)
		render.RenderError(err.Error())
		return
	}

	err = t.Save()
	if err != nil {
		log.Println(err)
		render.RenderError(err.Error())
		return

	}

	render.RenderJson(t, 1)
}
