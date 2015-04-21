package plugins

import (
	"github.com/guotie/deferinit"
	"github.com/smtc/justcms/database"
	"log"
	"time"
)

type PluginModel struct {
	ObjectId    string    `sql:"size:64" json:"object_id"`
	Name        string    `sql:"size:120" json:"name"`
	Author      string    `sql:"size:100" json:"author"`
	License     string    `sql:"size:40" json:"license"`
	Version     string    `sql:"size:24" json:"version"`
	Site        string    `sql:"size:255" json:"site"`
	Description string    `sql:"size:65536" json:"description"`
	Status      string    `sql:"status" json:"status"`
	InstalledAt time.Time `json:"installed_at"`
	ActivatedAt time.Time `json:"activated_at"`
}

type Plugin interface {
	Name() string
	Activate()
	Deactivate()
	Config()
}

var _plugins map[string]Plugin = make(map[string]Plugin)

// plugin模块初始化
// 调用defer init 初始化plugins
func init() {
	deferinit.AddInit(activateAll, nil, -10)
}

// 激活所有处于acitvate状态的plugin
func activateAll() {
	var (
		err    error
		models []PluginModel
		plugin Plugin
	)
	db := database.GetDB("")
	if err = db.Where("status=active").Find(&models).Error; err != nil {
		log.Printf("Fetch active plugin failed: %s\n", err.Error())
		return
	}
	for _, m := range models {
		plugin = _plugins[m.Name]
		if plugin != nil {
			plugin.Activate()
		}
	}
}

// Install
// 将Plugin Model保持到数据库中
// Param:
//
// Return:
//   err: error. if error occurs, err is not nil; or else err is nil
func Install(name, author, license, version, site, descp string) error {
	var (
		model PluginModel
		err   error
	)
	db := database.GetDB("")
	if err = db.Where("name=?", name).Find(&model).Error; err == nil {
		log.Printf("Plugin %s has been installed.\n", name)
		return err
	}
	model.Name = name
	model.Author = author
	model.License = license
	model.Version = version
	model.Site = site
	model.Description = descp
	model.InstalledAt = time.Now()

	if err = db.Save(&model).Error; err != nil {
		log.Printf("Install Plugin %s failed: %s\n", name, err.Error())
		return err
	}
	return nil
}

// Uninstall
func Uninstall(name string) (err error) {
	db := database.GetDB("")
	err = db.Where("name=?", name).Delete(PluginModel{}).Error

	return
}

// 将plugin加入到全局plugins map中
func Set(p Plugin) {
	_plugins[p.Name()] = p
}

// 根据name获取plugin
func Get(name string) Plugin {
	return _plugins[name]
}
