package models

import (
	"time"
)

/*
--
-- 表的结构 `wp_links`
--

CREATE TABLE IF NOT EXISTS `wp_links` (
  `link_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `link_url` varchar(255) NOT NULL DEFAULT '',
  `link_name` varchar(255) NOT NULL DEFAULT '',
  `link_image` varchar(255) NOT NULL DEFAULT '',
  `link_target` varchar(25) NOT NULL DEFAULT '',
  `link_description` varchar(255) NOT NULL DEFAULT '',
  `link_visible` varchar(20) NOT NULL DEFAULT 'Y',
  `link_owner` bigint(20) unsigned NOT NULL DEFAULT '1',
  `link_rating` int(11) NOT NULL DEFAULT '0',
  `link_updated` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `link_rel` varchar(255) NOT NULL DEFAULT '',
  `link_notes` mediumtext NOT NULL,
  `link_rss` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`link_id`),
  KEY `link_visible` (`link_visible`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 AUTO_INCREMENT=6 ;
*/

type Link struct {
	Id          int64     `json:"id"`
	Url         string    `sql:"size:255" json:"url"`
	Name        string    `sql:"size:255" json:"name"`
	Image       string    `sql:"size:255" json:"image"`
	Target      string    `sql:"size:255" json:"target"`
	Description string    `sql:"size:255" json:"description"`
	Visible     string    `sql:"size:20" json:"visible"`
	LinkOwner   int64     `json:"link_owner"`
	Rating      int       `json:"rating"`
	UpdatedAt   time.Time `json:"updated_at"`
	Relation    string    `sql:"size:255" json:"relation"`
	Notes       string    `sql:"type:text" json:"notes"`
	Rss         string    `sql:"size:255" json:"rss"`
}
