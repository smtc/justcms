package models

import (
	"fmt"
	//"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
	"log"
	"strconv"
	"strings"
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
	Status     string    `sql:"size:20" json:"status"`
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
	Id           int64  `json:id`
	PostId       int64  `json:"post_id"`
	PostObjectId string `sql:"size:64" json:"object_id"`
	MetaKey      string `sql:"size:300" json:"meta_key"`
	MetaValue    string `sql:"size:100000" json:"meta_value"`
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

func convertInt(value interface{}) (i int, ok bool) {
	i, ok = value.(int)
	if !ok {
		var i64 int64
		if i64, ok = value.(int64); ok {
			i = int(i64)
			return
		}
		var f float64
		if f, ok = value.(float64); ok {
			i = int(f)
		}
	}
	return
}

func convertInt64(value interface{}) (i64 int64, ok bool) {
	i64, ok = value.(int64)
	if !ok {
		var i int
		if i, ok = value.(int); ok {
			i64 = int64(i)
			return
		}
		var f float64
		if f, ok = value.(float64); ok {
			i64 = int64(f)
		}
	}
	return
}

func convertString(value interface{}) (s string, ok bool) {
	s, ok = value.(string)
	if !ok {
		var b []byte
		if b, ok = value.([]byte); ok {
			s = string(b)
			return
		}
		var bt byte
		if bt, ok = value.(byte); ok {
			s = string(bt)
			return
		}
	}
	return
}

func convertBool(value interface{}) (b bool, ok bool) {
	b, ok = value.(bool)
	if !ok {
		var s string
		if s, ok = convertString(value); ok {
			if s != "" {
				b = true
			}
			return
		}
		var i int
		if i, ok = convertInt(value); ok {
			if i != 0 {
				b = true
			}
			return
		}
	}
	return
}

// 如果是字符串， 以逗号分割
func convertIntArray(value interface{}) (ret []int, err error) {
	var ret64 []int64

	if ret64, err = convertInt64Array(value); err != nil {
		return
	}
	for _, i := range ret64 {
		ret = append(ret, int(i))
	}
	return
}

func convertInt64Array(value interface{}) (ret []int64, err error) {
	var (
		s   string
		sa  []string
		num int64
		ok  bool
		ia  []int
		fa  []float64
	)

	if s, ok = convertString(value); ok {
		sa = strings.Split(s, ",")
		for _, ele := range sa {
			ele = strings.TrimSpace(ele)
			if num, err = strconv.ParseInt(ele, 10, 64); err == nil {
				ret = append(ret, num)
			}
		}
		return
	}
	if ret, ok = value.([]int64); ok {
		return
	}
	if ia, ok = value.([]int); ok {
		for _, i := range ia {
			ret = append(ret, int64(i))
		}
		return
	}
	if fa, ok = value.([]float64); ok {
		for _, f := range fa {
			ret = append(ret, int64(f))
		}
		return
	}
	err = fmt.Errorf("Cannot convert value to []int64")
	return
}

func buildWhereClause() {
	return
}

// get posts
// main query function for posts
// param:
//  opt
// opt keys:
//   wordperss WP_Query
//
func GetPosts(opt map[string]interface{}) (posts []*Post, err error) {
	var post Post
	db := database.GetDB("")

	//
	if opt["id"] != nil || opt["pid"] != nil || opt["page_id"] != nil || opt["post_name"] != nil {
		var (
			id   int64
			name string
			ok   bool
		)
		if opt["id"] != nil {
			id, ok = convertInt64(opt["id"])
		}
		if !ok && opt["pid"] != nil {
			id, ok = convertInt64(opt["pid"])
		}
		if !ok && opt["page_id"] != nil {
			id, ok = convertInt64(opt["page_id"])
		}
		if ok {
			post.Id = id
			err = db.First(&post).Error
			posts = append(posts, &post)
			return
		}
		if opt["post_name"] != nil {
			if name, ok = convertString(opt["post_name"]); ok {
				err = db.Where("title=?", name).Find(posts).Error
				return
			}
		}
		log.Printf("param id or pid or page_id type error\n")
	}

	for key, value := range opt {
		switch key {
		// Author Parameters
		case "author":
			fallthrough
		case "author_id":
			author_id, ok := convertInt(value)
			if !ok {
				log.Printf("param author or author_id type error: should be int\n")
				continue
			}
			db = db.Where("author_id=?", author_id)

		case "author_name":
			if name, ok := convertString(value); ok {
				db = db.Where("author_name=?", name)
			} else {
				log.Println("param author_name type error.")
			}
		case "author__in":
		case "author__not_in":

		// Category Parameters
		case "cat": // int
			if cat, ok := convertInt64(value); ok {
			} else {
				log.Println("param cat type error.")
			}
		case "category_name":
		case "category__and":
		case "category__in":
		case "category__not_in":

		// Tag Parameters
		case "tag": // string
		case "tag_id":
		case "tag__and":
		case "tag__in":
		case "tag__not_in":
		case "tag_slug__and":
		case "tag_slug__in":

		// Taxonomy Parameters
		// Post & Page Parameters
		case "post_parent": // int
			if post_parent, ok := convertInt64(value); ok {
				db = db.Where("post_parent = ?", post_parent)
			} else {
				log.Println("param post_parent type error")
			}
		case "post_parent__in":
		case "post_parent__not_in":
		case "post__in":
		case "post__not_in":

		// Password Parameters
		case "has_password":
			if has_pwd, ok := convertBool(value); ok {
				if has_pwd {
					db = db.Where("post_password != ''")
				} else {
					db = db.Where("post_password == ''")
				}
			} else {
				log.Println("param has_password type error")
			}
		case "post_password":
			if post_pwd, ok := convertString(value); ok {
				db = db.Where("post_password=?", post_pwd)
			} else {
				log.Println("param post_password type error")
			}

		// Type Parameters
		case "post_type":
			if post_type, ok := convertString(value); ok {
				db = db.Where("post_type=?", post_type)
			} else {
				log.Println("param post_type type error")
			}

		// Status Parameters
		case "post_status":
			if post_status, ok := convertString(value); ok {
				db = db.Where("post_status=?", post_status)
			} else {
				log.Println("param post_status type error")
			}

		// Pagination Parameters
		case "nopaging":
		case "posts_per_page":
		case "posts_per_archive_page":
		case "offset":
			if offset, ok := convertInt(value); ok {
				db = db.Offset(offset)
			} else {
				log.Println("param offset type error")
			}
		case "page":
		case "ignore_sticky_posts":

		// Order & Orderby Parameters
		case "order":
		case "orderby":

		// Date Parameters
		case "date_before":
		case "date_after":
		}
	}

	err = db.Find(&posts).Error
	return
}
