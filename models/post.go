package models

import (
	"fmt"
	//"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
	//"log"
	"time"
)

// Post table
type Post struct {
	Id         int64     `json:"id"`
	ObjectId   string    `sql:"size:64" json:"object_id"`
	AuthorId   int64     `json:"author_id"`
	AuthorName string    `sql:"size:40" json:"author_name"`
	Title      string    `sql:"size:60" json:"title"`
	Content    string    `sql:"size:100000" json:"content"`
	PostStatus string    `sql:"size:20" json:"post_status"`
	Excerpt    string    `sql:"size:500" json:"excerpt"`
	PostAt     time.Time `json:"post_at"`
	PublishAt  time.Time `json:"publish_at"`
	ModifyAt   time.Time `json:"modify_at"`
	ClosedAt   time.Time `json:"closed_at"`

	ReplyStatus     string    `sql:"size:20;default:'open'" json:"reply_status"`
	PingStatus      string    `sql:"size:20;default:'open'" json:"ping_status"`
	PostName        string    `sql:"size:200" json:"post_name"`
	PostPassword    string    `sql:"size:20" json:"post_password"`
	ToPing          string    `sql:"type:text" json:"to_ping"`
	Pinged          string    `sql:"type:text" json:"pinged"`
	ContentFiltered string    `sql:"type:text" json:"content_filtered"`
	PostParent      int64     `json:"post_parent"`
	MenuOrder       int       `json:"menu_order"`
	PostType        string    `sql:"size:20" json:"post_type"`
	PostMimeType    string    `sql:"size:200" json:"post_mime_type"`
	ReplyCount      int64     `json:"reply_count"`
	LastReplyAt     time.Time `json:"last_reply_at"`
	LikedCount      int       `json:"liked_count"`
	BookmarkCount   int       `json:"bookmark_count"`
	StarCount       int       `json:"star_count"`
	BlockCount      int       `json:"block_count"`
}

// PostMeta table
type PostMeta struct {
	Id        int64  `json:id`
	PostId    int64  `json:"post_id"`
	ObjectId  string `sql:"size:64" json:"object_id"`
	MetaKey   string `sql:"size:300" json:"meta_key"`
	MetaValue string `sql:"size:100000" json:"meta_value"`
}

// create new post
// param:
//   opt: map[string]interface{}
// return:
//   post: *Post
//   err: error
func NewPost(opt map[string]interface{}) (post *Post, err error) {

	return
}

// get post by objectid
func GetPostByObjectId(oid string) (post *Post, err error) {
	db := database.GetDB("")
	err = db.Where("object_id=?", oid).First(post).Error
	return
}

// register post type
func RegisterPostType(typ string, opts map[string]interface{}) (err error) {
	defaultOpts := map[string]interface{}{
		"labels":               []string{},
		"description":          "",
		"public":               false,
		"hierarchical":         false,
		"exclude_from_search":  "",
		"publicly_queryable":   "",
		"show_ui":              "",
		"show_in_menu":         "",
		"show_in_nav_menus":    "",
		"show_in_admin_bar":    "",
		"menu_position":        "",
		"menu_icon":            "",
		"capability_type":      "post",
		"capabilities":         []string{},
		"map_meta_cap":         "",
		"supports":             []string{},
		"register_meta_box_cb": "",
		"taxonomies":           []string{},
		"has_archive":          false,
		"rewrite":              true,
		"query_var":            true,
		"can_export":           true,
		"delete_with_user":     "",
		"_builtin":             false,
		"_edit_link":           "post.php?post=%d",
	}
	mergeMap(defaultOpts, opts)

	postType := sanitizeKey(typ)
	if len(postType) >= 20 {
		return fmt.Errorf("post type length should NOT exceed 20.")
	}

	return
}
