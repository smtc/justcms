package models

import (
	//"time"
	"github.com/smtc/justcms/database"
)

/*
--
-- 表的结构 `wp_terms`
--
CREATE TABLE IF NOT EXISTS `wp_terms` (
  `term_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL DEFAULT '',
  `slug` varchar(200) NOT NULL DEFAULT '',
  `term_group` bigint(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`term_id`),
  UNIQUE KEY `slug` (`slug`),
  KEY `name` (`name`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 AUTO_INCREMENT=233 ;

--
-- 表的结构 `wp_term_relationships`
--
CREATE TABLE IF NOT EXISTS `wp_term_relationships` (
  `object_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `term_taxonomy_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `term_order` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`object_id`,`term_taxonomy_id`),
  KEY `term_taxonomy_id` (`term_taxonomy_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

--
-- 表的结构 `wp_term_taxonomy`
--
CREATE TABLE IF NOT EXISTS `wp_term_taxonomy` (
  `term_taxonomy_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `term_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `taxonomy` varchar(32) NOT NULL DEFAULT '',
  `description` longtext NOT NULL,
  `parent` bigint(20) unsigned NOT NULL DEFAULT '0',
  `count` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`term_taxonomy_id`),
  UNIQUE KEY `term_id_taxonomy` (`term_id`,`taxonomy`),
  KEY `taxonomy` (`taxonomy`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 AUTO_INCREMENT=245 ;
*/

type TermRelation struct {
	ObjectId  int64 `json:"object_id"`
	TermId    int64 `json:"term_id"`
	TermOrder int   `json:"term_order"`
}

type TermTaxonomy struct {
	Id          int64  `json:"id"`
	Name        string `sql:"size:200" json:"name"`
	Slug        string `sql:"size:200" json:"slug"`
	TermGroup   int    `json:"term_group"`
	Taxonomy    string `sql:"size:60" json:"taxonomy"`
	Description string `sql:"size:100000" json:"description"`
	Parent      int64  `json:"parent`
	Count       int64  `json:"count"`
}

// getCatNameById
// todo: 增加进程category缓存
func getCategoryIdById(cid int64) (id int64, err error) {
	var term TermTaxonomy

	db := database.GetDB("")

	if err = db.Where("taxonomy=category").Where("id=?", cid).First(&term).Error; err != nil {
		return
	}
	id = term.Id
	return
}

// getCategoryIdByName
// 根据category name获取category ID
func getCategoryIdByName(cname string) (id int64, err error) {
	var term TermTaxonomy

	db := database.GetDB("")

	if err = db.Where("taxonomy=category").Where("name=?", cname).First(&term).Error; err != nil {
		return
	}
	id = term.Id
	return
}

// getTagIdByName
// 根据tag name获取tag ID
func getTagIdByName(tag string) (id int64, err error) {
	var term TermTaxonomy

	db := database.GetDB("")

	if err = db.Where("taxonomy=tag").Where("name=?", tag).First(&term).Error; err != nil {
		return
	}
	id = term.Id

	return
}
