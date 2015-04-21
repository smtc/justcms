package meta

import (
	"sync"
	//"time"
	"github.com/smtc/justcms/database"
	"strconv"
)

/*
--
-- 表的结构 `wp_options`
--

CREATE TABLE IF NOT EXISTS `wp_options` (
  `option_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `blog_id` int(11) NOT NULL DEFAULT '0',
  `option_name` varchar(64) NOT NULL DEFAULT '',
  `option_value` longtext NOT NULL,
  `autoload` varchar(20) NOT NULL DEFAULT 'yes',
  PRIMARY KEY (`option_id`),
  UNIQUE KEY `option_name` (`option_name`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 AUTO_INCREMENT=857 ;
*/

type Options struct {
	Id       int64  `json:"id"`
	Name     string `sql:"size:64" json:"name"`
	Value    string `sql:"size:100000" json:"value"`
	Autoload string `sql:"size:20" json:"autoload"`
}

var (
	// 缓存所有已经查询过的options
	_cacheOptions = struct {
		sync.Mutex
		options map[string]*Options
	}{options: make(map[string]*Options)}

	_defaultOptions = map[string]interface{}{
		"template": "default",
	}
)

// GetOptionByName
// params:
//   name: option name
//   def: if not found, the default value
// return:
//   opt: the option value
//   err: error
func GetOption(name string) (opt Options, err error) {
	_cacheOptions.Lock()
	if op, ok := _cacheOptions.options[name]; ok {
		_cacheOptions.Unlock()
		return *op, nil
	}
	_cacheOptions.Unlock()

	db := database.GetDB("")

	err = db.Where("name=?", name).Limit(1).Find(&opt).Error
	if err != nil {
		return
	}

	return
}

func GetStringOptions(name string) (val string, err error) {
	var opt Options
	if opt, err = GetOption(name); err == nil {
		return opt.Value, nil
	}
	return
}

func GetIntOption(name string) (val int, err error) {
	var (
		opt Options
	)
	opt, err = GetOption(name)

	val, err = strconv.Atoi(opt.Value)
	return

}

func GetInt64Option(name string) (val int64, err error) {
	var opt Options

	opt, err = GetOption(name)
	val, err = strconv.ParseInt(opt.Value, 10, 64)
	return
}

// update options
// 2014-10-10
func UpdateOptions(name, val string) (err error) {
	opt, err := GetOption(name)
	if err != nil {
		return
	}

	opt.Value = val

	db := database.GetDB("")
	if err = db.Save(&opt).Error; err != nil {
		return
	}

	_cacheOptions.options[opt.Name] = &opt
	return
}
