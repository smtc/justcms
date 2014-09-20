package models

type Select struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

// 用,把数组t连接起来
// 使用 goutils.ToString(t, "") 代替！！

func conjectToString(t []int64) string {
	s := ""
	for i, v := range t {
		s += fmt.Sprint(v)
		if i != len(t)-1 {
			s += ","
		}
	}
	return s
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
