package models

import (
	"fmt"
	//"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
	"log"
	"net/http"
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

// 转化为int
func convertInt(key, value string) (i int, err error) {
	var i64 int64
	i64, err = convertInt64(key, value)
	i = int(i64)
	return
}

// 转化为int64
func convertInt64(key, value string) (i int64, err error) {
	vv := strings.TrimSpace(value)
	if i, err = strconv.ParseInt(vv, 10, 64); err == nil {
		return
	}
	return
}

// 转化为布尔值
func convertBool(key, value string) (b bool, err error) {
	vv := strings.ToLower(strings.TrimSpace(value))

	if vv == "true" {
		b = true
	}
	if vv != "false" {
		err = fmt.Errorf("param key %s value %s is NOT boolean.", key, value)
	}
	return
}

// 如果是字符串， 以逗号分割
func convertIntArray(key, value string) (ret []int, err error) {
	var ret64 []int64

	if ret64, err = convertInt64Array(key, value); err != nil {
		return
	}
	for _, i := range ret64 {
		ret = append(ret, int(i))
	}
	return
}

func convertInt64Array(key, value string) (ret []int64, err error) {
	var (
		sa  []string
		num int64
	)

	sa = strings.Split(value, ",")
	for _, ele := range sa {
		ele = strings.TrimSpace(ele)
		if num, err = strconv.ParseInt(ele, 10, 64); err == nil {
			ret = append(ret, num)
		}
	}
	return

	err = fmt.Errorf("Cannot convert param %s value %s to []int64", key, value)
	return
}

func buildWhereClause() {
	return
}

// split "1,2,3" to [1,2,3], [], nil
//    "-1, -2, -3" to [], [-1,-2,-3,], nil
//    "1, -2, 3" to [1,3], [-2], nil
func split_ids(key, ss string) (ia []int64, ea []int64, err error) {
	var i int64

	sa := strings.Split(ss, ",")
	ia = make([]int64, 0)
	ea = make([]int64, 0)
	for _, s := range sa {
		s = strings.TrimSpace(s)
		if i, err = strconv.ParseInt(s, 10, 64); err != nil {
			log.Printf("split %s id array [%s] failed:", key, ss, err.Error())
			continue
		}
		if i > 0 {
			ia = append(ia, i)
		} else if i < 0 {
			ea = append(ea, i)
		} else {
			log.Println("split_ids: get invlid id 0: ", s)
		}
	}

	if len(ia) == 0 && len(ea) == 0 {
		err = fmt.Errorf("No valid id found")
	}

	return
}

// split "a,b,c" to ["a", "b", "c"], rel is 0 (or)
// split "a+b+c" to ["a", "b", "c"], rel is 1 (and)
func split_sa(key, ss string) (sa []string, rel int, err error) {
	a := strings.Split(ss, ",")
	if len(a) == 0 {
		a := strings.Split(ss, "+")
		if len(a) == 0 {
			err = fmt.Errorf("param %s is empty: %s.", key, ss)
			return
		}
		rel = 1
	}
	sa = make([]string, 0)
	for _, s := range a {
		s = strings.TrimSpace(s)
		if s != "" {
			sa = append(sa, s)
		}
	}
	return
}

// 参数用逗号分割，为或的关系
// 参数用＋分割，为与的关系
// 简单参数列表：
//   id pid page_id
//   post_name
//   post__in post__not_in
//   author author_id author_name author__in author__not_in
//   cat category_name category__in category__not_in category__and category_name__in category_name__not_in
//   tag tag_id tag__in tag__not_in tag__and tag_slug__in tag_slug__not_in
//   post_parent post_parent__in post_parent__not_in
//   has_password
//   post_password
//   post_type
//   post_status
//   nopaging posts_per_page posts_per_archive_page offset page ignore_sticky_posts
//   date__before date__after

// 将category_name__in, category_name__and转化为category__in, category__and
// 将tag_slug__in, tag_slug__and转化为tag__in, tag__and
//
func parseQuery(r *http.Request) (opt map[string]interface{}, err error) {
	var (
		num int
		id  int64
		ids []int64
		eds []int64
		sa  []string
		rel int
		ret bool
	)

	r.ParseForm()
	opt = make(map[string]interface{})

	f := r.Form
	for key, _ := range f {
		value := f.Get(key)
		switch key {
		case "id":
			fallthrough
		case "pid":
			fallthrough
		case "page_id":
			if id, ok := convertInt64(key, value); ok != nil {
				opt["id"] = id
			} else {
				err = fmt.Errorf("Post Id not valid.")
			}
			return

		case "post_name":
			opt["post_name"] = value
			return

		// multi-rows

		// author
		case "author":
			fallthrough
		case "author_id":
			if ids, eds, err = split_ids(key, value); err != nil {
				return
			}
			if len(ids) == 1 {
				opt["author"] = ids[0]
			} else if len(ids) > 1 {
				opt["author__in"] = ids
			}
			if len(eds) >= 1 {
				opt["author__not_in"] = eds
			}
		case "author_name":
			name := strings.TrimSpace(f.Get(key))
			if name == "" {
				log.Printf("author_name is empty.\n")
			} else {
				opt["author_name"] = name
			}

		// category
		case "cat":
			if ids, eds, err = split_ids(key, value); err != nil {
				return
			}
			if len(ids) == 1 {
				opt["cat"] = ids[0]
			} else if len(ids) > 1 {
				opt["category__in"] = ids
			}
			if len(eds) > 0 {
				opt["category__not_in"] = eds
			}
		case "category_name":
			if sa, rel, err = split_sa(key, value); err != nil {
				return
			}
			if len(sa) == 1 {
				catname := sa[0]
				id, err = getCategoryIdByName(catname)
				if err != nil {
					err = fmt.Errorf("category name %s not exist.", catname)
					return
				} else {
					opt["cat"] = id
				}
			} else {
				ids = make([]int64, 0)
				for _, name := range sa {
					cid, err := getCategoryIdByName(name)
					if err != nil {
						log.Printf("category name %s is not exist.\n", name)
					}
					ids = append(ids, cid)
				}
				if len(ids) == 0 {
					err = fmt.Errorf("No valid category: %s", value)
					return
				}
				if len(ids) == 1 {
					opt["cat"] = ids[0]
				} else {
					if rel == 0 {
						opt["category__in"] = ids
					} else {
						opt["category__and"] = ids
					}
				}
			}

		// tag
		case "tag":
			if sa, rel, err = split_sa(key, value); err != nil {
				return
			}

			if len(sa) == 1 {
				tagname := sa[0]
				id, err = getTagIdByName(tagname)
				if err != nil {
					err = fmt.Errorf("tag name %s not exist.", tagname)
					return
				}
				opt["tag_id"] = id
			} else {
				ids = make([]int64, 0)
				for _, name := range sa {
					tid, err := getTagIdByName(name)
					if err != nil {
						log.Printf("tag name %s is not exist.\n", name)
					}
					ids = append(ids, tid)
				}
				if len(ids) == 0 {
					err = fmt.Errorf("No valid tag: %s", value)
					return
				} else if len(ids) == 1 {
					opt["tag_id"] = ids[0]
				} else if rel == 0 {
					opt["tag_id__in"] = ids
				} else {
					opt["tag_id__and"] = ids
				}
			}
		case "tag_id":
			if ids, eds, err = split_ids(key, value); err != nil {
				return
			}
			if len(ids) == 1 {
				opt["tag"] = ids[0]
			} else if len(ids) > 1 {
				opt["tag__in"] = ids
			}
			if len(eds) > 0 {
				opt["tag__not_in"] = eds
			}

		// Taxonomy

		// post_parent
		case "post_parent":
			if ids, eds, err = split_ids(key, value); err != nil {
				return
			}
			if len(ids) == 1 {
				opt["post_parent"] = ids[0]
			} else if len(ids) > 1 {
				opt["post_parent__in"] = ids
			}
			if len(eds) > 0 {
				opt["post_parent__not_in"] = eds
			}
		case "post":
			if ids, eds, err = split_ids(key, value); err != nil {
				return
			}
			if len(ids) == 1 {
				// 如果只有一个post被选择， 不再处理后面的选项，直接返回
				opt["id"] = ids[0]
				return
			} else if len(ids) > 1 {
				opt["post__id"] = ids
			}

			if len(eds) != 0 {
				opt["post__not_in"] = eds
			}

		// password
		case "has_password":
			if ret, err = convertBool(key, value); err != nil {
				log.Println(err)
				continue
			}
			opt["has_password"] = ret
		case "post_password":
			opt["post_password"] = value

		// post_type
		case "post_type":
			if sa, rel, err = split_sa(key, value); err != nil {
				return
			}
			if len(sa) == 1 {
				opt["post_type"] = sa[0]
			} else {
				if rel == 1 {
					log.Printf("param post_type cannot have relation AND: %s\n", value)
					return
				}
				opt["post_type__in"] = sa
			}

		// post_status
		case "post_status":
			if sa, rel, err = split_sa(key, value); err != nil {
				return
			}
			if len(sa) == 1 {
				opt["post_status"] = sa[0]
			} else {
				if rel == 1 {
					log.Printf("param post_status cannot have relation AND: %s\n", value)
					continue
				}
				opt["post_status__in"] = sa
			}

		// pagination
		case "nopaging":
			if ret, err = convertBool(key, value); err != nil {
				log.Println(err)
				continue
			}
			opt["nopaging"] = ret
		case "posts_per_page":
			if num, err = convertInt(key, value); err != nil {
				log.Println(err)
				continue
			}
			opt["posts_per_page"] = num
		case "posts_per_archive_page":
			if num, err = convertInt(key, value); err != nil {
				log.Println(err)
				continue
			}
			opt["posts_per_archive_page"] = num
		case "offset":
			if num, err = convertInt(key, value); err != nil {
				log.Println(err)
				continue
			}
			opt["offset"] = num
		case "paged":
			fallthrough
		case "page":
			num, err = convertInt(key, value)
			if err != nil || num < 1 {
				log.Println("param page invalid: %d %s\n", num, err)
				opt["page"] = 1
				continue
			}
			opt["page"] = num

		case "ignore_sticky_posts":
			if ret, err = convertBool(key, value); err != nil {
				log.Println(err)
				continue
			}
			opt["ignore_sticky_posts"] = ret

		// Order & Orderby Parameters
		case "order":
			value = strings.ToLower(strings.TrimSpace(value))
			if value == "asc" {
				opt["order"] = "asc"
			} else if value == "desc" {
				opt["order"] = "desc"
			} else {
				log.Printf("param order invalid: %s, set it to desc\n", value)
				opt["order"] = "desc"
			}
		case "orderby":
			opt["orderby"] = strings.TrimSpace(value)

		// Date Parameters
		case "date_before":
		case "date_after":
		}
	}

	if opt["post_type"] == nil {
		opt["post_type"] = "post"
	}
	if opt["post_status"] == nil {
		opt["post_staus"] = "published"
	}

	return
}

// get posts
// main query function for posts
// param:
//  opt
// opt keys:
//   wordperss WP_Query
//
func GetPosts(req *http.Request) (posts []*Post, err error) {
	var (
		post Post
		opt  map[string]interface{}
	)

	if opt, err = parseQuery(req); err != nil {
		return
	}

	db := database.GetDB("")

	// 查询单个post
	if opt["id"] != nil {
		if err = db.Where("id=?", opt["id"]).Find(&post).Error; err != nil {
			return
		}
		posts = append(posts, &post)
		return
	}
	if opt["post_name"] != nil {
		if err = db.Where("post_name=?", opt["post_name"]).Find(&post).Error; err != nil {
			return
		}
		posts = append(posts, &post)
		return
	}

	// 查询多个posts
	err = db.Find(&posts).Error

	return
}

func mergePosts() {

}
