package models

import (
	"time"
)

/*
--
-- 表的结构 `wp_commentmeta`
--
--
CREATE TABLE IF NOT EXISTS `wp_commentmeta` (
  `meta_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `comment_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `meta_key` varchar(255) DEFAULT NULL,
  `meta_value` longtext,
  PRIMARY KEY (`meta_id`),
  KEY `comment_id` (`comment_id`),
  KEY `meta_key` (`meta_key`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 AUTO_INCREMENT=6673 ;

--
-- 表的结构 `wp_comments`
--
CREATE TABLE IF NOT EXISTS `wp_comments` (
  `comment_ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `comment_post_ID` bigint(20) unsigned NOT NULL DEFAULT '0',
  `comment_author` tinytext NOT NULL,
  `comment_author_email` varchar(100) NOT NULL DEFAULT '',
  `comment_author_url` varchar(200) NOT NULL DEFAULT '',
  `comment_author_IP` varchar(100) NOT NULL DEFAULT '',
  `comment_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `comment_date_gmt` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `comment_content` text NOT NULL,
  `comment_karma` int(11) NOT NULL DEFAULT '0',
  `comment_approved` varchar(20) NOT NULL DEFAULT '1',
  `comment_agent` varchar(255) NOT NULL DEFAULT '',
  `comment_type` varchar(20) NOT NULL DEFAULT '',
  `comment_parent` bigint(20) unsigned NOT NULL DEFAULT '0',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`comment_ID`),
  KEY `comment_approved` (`comment_approved`),
  KEY `comment_post_ID` (`comment_post_ID`),
  KEY `comment_approved_date_gmt` (`comment_approved`,`comment_date_gmt`),
  KEY `comment_date_gmt` (`comment_date_gmt`),
  KEY `comment_parent` (`comment_parent`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 AUTO_INCREMENT=3237 ;
*/

type Comment struct {
	Id          int64     `json:"id"`
	ObjectId    string    `sql:"size:64" json:"object_id"`
	PostId      int64     `json:"post_id"`
	Author      string    `sql:"size:40" json:"author"`
	Email       string    `sql:"size:100" json:"email"`
	AuthorUrl   string    `sql:"size:40" json:"author"`
	AuthorIP    string    `sql:"size:40" json:"author"`
	Content     string    `sql:"size:60000" json:"content"`
	Karma       int       `json:"karma"`
	Approved    string    `sql:"size:20" json:"approved"`
	Agent       string    `sql:"size:255" json:"agent"`
	ReplyType   string    `sql:"size:20" json:"reply_type"`
	ReplyParent string    `sql:"size:64" json:"reply_parent"`
	AccountId   int64     `json:"account_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CommentMeta struct {
	Id              int64  `json:"id"`
	CommentId       int64  `json:"comment_id"`
	CommentObejctId string `sql:"size:64" json:"comment_objectid"`
	MetaKey         string `sql:"size:300" json:"meta_key"`
	MetaValue       string `sql:"size:100000" json:"meta_value"`
}
