package models

import (
	"time"
)

// Post table
type Post struct {
	Id        int64     `json:"id"`
	ObjectId  string    `sql:"size:64" json:"object_id"`
	AuthorId  int64     `json:"author_id"`
	Title     string    `sql:"size:60" json:"title"`
	Content   string    `sql:"text" json:"content"`
	Status    string    `sql:"size:20" json:"status"`
	Excerpt   string    `sql:"size:500" json:"excerpt"`
	PostAt    time.Time `json:"post_at"`
	PublishAt time.Time `json:"publish_at"`
	ModifyAt  time.Time `json:"modify_at"`

	ReplyStatus     string `sql:"size:20;default:open" json:"reply_status"`
	PingStatus      string `sql:"size:20;default:open" json:"ping_status"`
	PostName        string `sql:"size:200" json:"post_name"`
	PostPassword    string `sql:"size:20" json:"post_password"`
	ToPing          string `sql:"type:text" json:"to_ping"`
	Pinged          string `sql:"type:text" json:"pinged"`
	ContentFiltered string `sql:"type:text" json:"content_filtered"`
	PostParent      int64  `json:"post_parent"`
	MenuOrder       int    `json:"menu_order"`
	PostType        string `sql:"size:20" json:"post_type"`
	PostMimeType    string `sql:"size:200" json:"post_mime_type"`
	ReplyCount      int64  `json:"reply_count"`
}

// PostMeta table
type PostMeta struct {
	Id           int64  `json:id`
	PostId       int64  `json:"post_id"`
	PostObjectId string `sql:"size:64" json:"object_id"`
	MetaKey      string `sql:"size:300" json:"meta_key"`
	MetaValue    string `sql:"size:100000" json:"meta_value"`
}
