package models

import (
	//"time"
	"fmt"
	"github.com/smtc/justcms/database"
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

var defaultOptions = map[string]interface{}{
	"template": "default",
}

// GetOptionByName
// params:
//   name: option name
//   def: if not found, the default value
// return:
//   opt: the option value
//   err: error
func GetOptionByName(name string) (opt interface{}, err error) {
	var row Options

	db := database.GetDB("")

	err = db.Where("name=?", name).Limit(1).Find(&row).Error
	if err != nil {
		opt = defaultOptions[name]
		return
	}
	opt = row.Value
	return
}

func GetStringOptions(name string) (val string, err error) {
	opt, err := GetOptionByName(name)
	var ok bool
	if val, ok = opt.(string); !ok {
		err = fmt.Errorf("Cannot convert option %s to type string", name)
	}
	return
}

func GetIntOption(name string) (val int, err error) {
	opt, err := GetOptionByName(name)
	var (
		ok    bool
		val64 int64
	)
	if val64, ok = opt.(int64); !ok {
		err = fmt.Errorf("Cannot convert option %s to type string", name)
	}
	val = int(val64)
	return

}

func GetInt64Option(name string) (val int64, err error) {

	opt, err := GetOptionByName(name)
	var (
		ok bool
	)
	if val, ok = opt.(int64); !ok {
		err = fmt.Errorf("Cannot convert option %s to type string", name)
	}
	return
}
