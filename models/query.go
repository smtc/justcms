package models

// 查询posts
import (
	"fmt"
	//"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"net/url"
	"strings"
	//"time"

	"github.com/smtc/justcms/database"
)

type queryClause struct {
	join     string
	where    string
	limits   string
	groupby  string
	orderby  string
	distinct string
}

func sqlIn(ar []int64) (s string) {
	for i, v := range ar {
		if i != len(ar) {
			s += fmt.Sprintf("%d,", v)
		} else {
			s += fmt.Sprintf("%d", v)
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
//   cat_children
//   tag tag_id tag__in tag__not_in tag__and tag_slug__in tag_slug__not_in
//   post_parent post_parent__in post_parent__not_in
//   has_password
//   post_password
//   post_type
//   post_status
//   nopaging posts_per_page posts_per_archive_page offset page ignore_sticky_posts
//   date__before date__after
//
// 将author_name 转化为author__in, author__and
// 将category_name__in, category_name__and转化为category__in, category__and
// 将tag_slug__in, tag_slug__and转化为tag__in, tag__and
//
func parseQuery(r string) (opt map[string]interface{}, err error) {
	var (
		num   int
		id    int64
		ids   []int64
		eds   []int64
		sa    []string
		rel   int
		ret   bool
		query url.Values
	)

	// ParseQuery has do QueryUnescape
	if query, err = url.ParseQuery(r); err != nil {
		return
	}

	opt = make(map[string]interface{})

	for key, _ := range query {
		value := query.Get(key)
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
			if ids, eds, err = splitIds(key, value); err != nil {
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
			if sa, rel, err = splitSa(key, value); err != nil {
				return
			}
			if len(sa) == 1 {
				id, err = getAuthorIdByName(sa[0])
				if err != nil {
					return
				}
				opt["author"] = id
			} else {
				ids = make([]int64, 0)
				for _, name := range sa {
					aid, err := getAuthorIdByName(name)
					if err != nil {
						log.Printf("author name %s is not exist.\n", name)
					}
					ids = append(ids, aid)
				}
				if len(ids) == 0 {
					err = fmt.Errorf("No valid author name: %s", value)
					return
				}
				if len(ids) == 1 {
					opt["author"] = ids[0]
				} else {
					if rel == 0 {
						opt["author__in"] = ids
					} else {
						err = fmt.Errorf("no post with multi author currently.")
						return
					}
				}
			}

		// category
		case "cat":
			if ids, eds, err = splitIds(key, value); err != nil {
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
		case "cat_children":
			cc := strings.ToLower(value)
			if cc == "true" || cc == "1" {
				opt["cat_children"] = true
			} else if cc == "false" || cc == "0" {
				opt["cat_children"] = false
			} else {
				log.Println("param cat_children invalid " + value)
			}
		case "category_name":
			if sa, rel, err = splitSa(key, value); err != nil {
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
			if sa, rel, err = splitSa(key, value); err != nil {
				return
			}

			if len(sa) == 1 {
				tagname := sa[0]
				id, err = getTagIdByName(tagname)
				if err != nil {
					err = fmt.Errorf("tag name %s not exist.", tagname)
					return
				}
				opt["tag"] = id
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
					opt["tag"] = ids[0]
				} else if rel == 0 {
					opt["tag__in"] = ids
				} else {
					opt["tag__and"] = ids
				}
			}
		case "tag_id":
			if ids, eds, err = splitIds(key, value); err != nil {
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
			if ids, eds, err = splitIds(key, value); err != nil {
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
			if ids, eds, err = splitIds(key, value); err != nil {
				return
			}
			if len(ids) == 1 {
				// 如果只有一个post被选择， 不再处理后面的选项，直接返回
				opt["id"] = ids[0]
				return
			} else if len(ids) > 1 {
				opt["post__in"] = ids
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
			if sa, rel, err = splitSa(key, value); err != nil {
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
			if sa, rel, err = splitSa(key, value); err != nil {
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
				opt["order"] = "ASC"
			} else if value == "desc" {
				opt["order"] = "DESC"
			} else {
				log.Printf("param order invalid: %s, set it to desc\n", value)
				opt["order"] = "DESC"
			}
		case "orderby":
			opt["orderby"] = strings.TrimSpace(value)

		// Date Parameters
		case "date_before":
		case "date_after":
		case "menu_order":
			opt["menu_order"] = value
		}
	}

	if opt["post_type"] == nil && opt["post_type__in"] == nil {
		opt["post_type"] = "post"
	}
	if opt["post_status"] == nil && opt["post_status__in"] == nil {
		opt["post_staus"] = "published"
	}

	return
}

// author, author__in, author__not_in
// cat, category__in, category__not_in, category__and
// tag, tag__in, tag__not_in, tag__and
// post_parent, post_parent__in, post_parent__not_in
// post__in, post__not_in
// has_password
// post_password
// post_type, post_type__in
// post_status, post_status__in
// nopaging
// posts_per_page
// posts_per_archive_page
// offset
// page
// ignore_sticky_posts
// order
// orderby
// date_before, date_after
// menu_order
func buildWhereClause(opt map[string]interface{}) (clause []string, err error) {
	var (
		join    = ""
		qc, cqc queryClause
	)

	_ = join

	if opt["menu_order"] != nil {
		qc.where += " And menu_order = " + opt["menu_order"].(string)
	}

	// taxonomy
	// 目前getTaxSql不会返回错误
	cqc, err = getTaxSql(buildTaxQuery(opt), "AND", "posts", "id")
	qc.where += cqc.where
	qc.join += cqc.join

	// author, user stuff
	cqc, err = getAuthorSql(opt)
	qc.where += cqc.where

	// order, order by
	if opt["order"] == nil {
		opt["order"] = "DESC"
	}
	if opt["orderby"] != nil {
		qc.orderby = parseOrderBy(opt["orderby"].(string))
	}

	// post status

	// post type

	return
}

// 构建taxQuery数组
func buildTaxQuery(opt map[string]interface{}) []taxQuery {
	// category has children
	cc := false
	ta := make([]taxQuery, 0)
	if opt["cat_children"] != nil && opt["cat_children"].(bool) {
		cc = true
	}

	for key, value := range opt {
		switch key {
		case "tag__in":
			q := taxQuery{"tag", value.([]int64), "IN", false}
			ta = append(ta, q)
		case "tag__not_in":
			q := taxQuery{"tag", value.([]int64), "NOT IN", false}
			ta = append(ta, q)
		case "tag__and":
			q := taxQuery{"tag", value.([]int64), "AND", false}
			ta = append(ta, q)

		// category
		case "category__in":
			q := taxQuery{"category", value.([]int64), "IN", cc}
			ta = append(ta, q)
		case "category__not_in":
			q := taxQuery{"category", value.([]int64), "NOT IN", cc}
			ta = append(ta, q)
		case "category__and":
			q := taxQuery{"category", value.([]int64), "AND", cc}
			ta = append(ta, q)
		}
	}

	return ta
}

// get posts
// main query function for posts
// param:
//  opt
// opt keys:
//   wordperss WP_Query
//
func GetPosts(req *http.Request) (posts []*Post, err error) {
	var opt map[string]interface{}

	if opt, err = parseQuery(req.URL.RawQuery); err != nil {
		return
	}

	return getPosts(opt)
}

func getPosts(opt map[string]interface{}) (posts []*Post, err error) {
	var (
		post  Post
		where []string
	)

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

	where, err = buildWhereClause(opt)
	_ = where
	// 查询多个posts
	err = db.Where(where).Find(&posts).Error

	return
}
