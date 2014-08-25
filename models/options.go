package models

import (
//"time"
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
	BlogId   int64  `json:"blog_id"`
	Name     string `sql:"size:64" json:"name"`
	Value    string `sql:"size:100000" json:"value"`
	Autoload string `sql:"size:20" json:"autoload"`
}
