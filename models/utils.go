package models

type Select struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

// 用,把数组t连接起来
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
