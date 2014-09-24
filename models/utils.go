package models

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Select struct {
	Text  string `json:"text"`
	Value string `json:"value"`
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
	} else if vv != "false" {
		err = fmt.Errorf("param key %s value %s is NOT boolean.", key, value)
	}
	return
}

// 如果是字符串， 以逗号分割

// split "1,2,3" to [1,2,3], [], nil
//    "-1, -2, -3" to [], [-1,-2,-3,], nil
//    "1, -2, 3" to [1,3], [-2], nil
func splitIds(key, ss string) (ia []int64, ea []int64, err error) {
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
			ea = append(ea, -i)
		} else {
			log.Println("splitIds: get invlid id 0: ", s)
		}
	}

	if len(ia) == 0 && len(ea) == 0 {
		err = fmt.Errorf("No valid id found")
	}

	return
}

// split "a,b,c" to ["a", "b", "c"], rel is 0 (or)
// split "a+b+c" to ["a", "b", "c"], rel is 1 (and)
func splitSa(key, ss string) (sa []string, rel int, err error) {
	a := strings.Split(ss, ",")
	if len(a) == 1 {
		a = strings.Split(ss, "+")
		rel = 1
	}
	sa = make([]string, 0)
	for _, s := range a {
		s = strings.TrimSpace(s)
		if s != "" {
			sa = append(sa, s)
		}
	}
	if len(sa) == 0 {
		err = fmt.Errorf("param %s is empty: %s.", key, ss)
		return
	}

	return
}

func parseOrderBy(orderby string) string {
	allowed_keys := map[string]string{"post_name": "post_name",
		"post_author":   "author_name",
		"post_date":     "publish_at",
		"post_title":    "title",
		"post_modified": "modify_at",
		"post_parent":   "post_parent",
		"post_type":     "post_type",
		"name":          "post_name",
		"author":        "author_id",
		"date":          "publish_at",
		"title":         "title",
		"modified":      "modify_at",
		"parent":        "post_parent",
		"type":          "post_type",
		"id":            "id",
		"menu_order":    "menu_order",
	}
	// todo: meta_key

	return allowed_keys[orderby]
}

//filter map
// caution: interface{} must be simple values, or else it will panic!!!
func filterMap(m []map[string]interface{}, cond map[string]interface{}, op string) []map[string]interface{} {
	var (
		match = 0
		ret   = make([]map[string]interface{}, 0)
	)

	op = strings.ToUpper(op)
	if op != "AND" && op != "OR" && op != "NOT" {
		return ret
	}
	condCount := len(cond)

	for _, v := range m {
		match = 0
		for ck, cv := range cond {
			if v[ck] != nil && v[ck] == cv {
				match++
			}
		}
		if op == "AND" && match == condCount {
			ret = append(ret, v)
		}
		if op == "OR" && match > 0 {
			ret = append(ret, v)
		}
		if op == "NOT" && match == 0 {
			ret = append(ret, v)
		}
	}
	return ret
}

// merge src to dst
//   if src has entry A and dst has no entry A, then dst will has entry A
func mergeMap(src map[string]interface{}, dst map[string]interface{}) {
	for key, value := range src {
		if dst[key] == nil {
			dst[key] = value
		}
	}
}

/**
 * Sanitizes a string key.
 *
 * Keys are used as internal identifiers. Lowercase alphanumeric characters, dashes and underscores are allowed.
 *
 * @param string $key String key
 * @return string Sanitized key
 */
var sanitizeRe = regexp.MustCompile("[^a-z0-9_\\-]")

func sanitizeKey(key string) string {
	raw_key := key
	key = strings.ToLower(key)
	key = sanitizeRe.ReplaceAllString(key, "")

	/**
	 * Filter a sanitized key string.
	 *
	 * @since 3.0.0
	 *
	 * @param string $key     Sanitized key.
	 * @param string $raw_key The key prior to sanitization.
	 */
	_ = raw_key
	//apply_filters( 'sanitize_key', $key, $raw_key );
	return key
}

// 把value转换为sql表达式中的值
//   如下几种情况：
//      0 nil
//      1 string
//      2 []int64
//      3 []string
//      4 others
func sqlValue(v interface{}) (string, bool) {
	if v == nil {
		return "''", false
	}
	if s, ok := v.(string); ok {
		return "'" + s + "'", false
	}
	if sa, ok := v.([]string); ok {
		var res string = "("
		for i, s := range sa {
			res += "'" + s + "'"
			if i != len(sa)-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	}

	arrlen := 0
	switch v.(type) {
	case int:
		return fmt.Sprint(v), false
	case int8:
		return fmt.Sprint(v), false
	case int16:
		return fmt.Sprint(v), false
	case int32:
		return fmt.Sprint(v), false
	case int64:
		return fmt.Sprint(v), false
	case uint:
		return fmt.Sprint(v), false
	case uint8:
		return fmt.Sprint(v), false
	case uint16:
		return fmt.Sprint(v), false
	case uint32:

		return fmt.Sprint(v), false
	case uint64:
		return fmt.Sprint(v), false

	case []int:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	case []int8:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	case []int16:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	case []int32:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	case []int64:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	case []uint:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	case []uint8:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	case []uint16:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	case []uint32:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true
	case []uint64:
		arrlen = len(v.([]int)) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, ii := range v.([]int) {
			res += fmt.Sprint(ii)
			if i != arrlen-1 {
				res += ","
			}
		}
		res += ")"
		return res, true

	case []interface{}:
		vv := v.([]interface{})
		arrlen = len(vv) - 1
		if arrlen < 0 {
			return "''", false
		}
		res := "("
		for i, vi := range vv {
			res += "'" + fmt.Sprint(vi) + "'"
			if i != arrlen {
				res += ","
			}
		}
		res += ")"
		return res, true

	default:
		return fmt.Sprint(v), false
	}

	return "''", false
}
