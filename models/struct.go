package models

/*
"status": 1, // 状态 1为成功
    "src": "example/user/", // api地址，使用RESTful架构
    "op": {
        "unit": { // 行内操作按钮
            "edit": "form.html?id={{id}}", // 内置 edit
            "del": "id={{id}}" // 内置 del
        },
        "mult": { // 表格操作按钮
            "new": "form.html", // 内置 新建
            "refresh": "", // 内置 刷新
            "del": "id" // 内置 批量删除
        }
    },
    "struct": [
        { "key": "id", "text": "id", "hide": true },
        { "key": "name", "text": "姓名", "type": "text", "filter": "like" },
        { "key": "role", "text": "角色", "type": "select", "filter": "eq", "src": "json/select.json" },
        { "key": "ip", "text": "ip" },
        { "key": "time", "text": "时间", "type": "date" },
        { "key": "time", "text": "开始时间", "hide": true, "type": "date", "filter": "gt" },
        { "key": "time", "text": "结束时间", "hide": true, "type": "date", "filter": "lt" },
        { "key": "status", "text": "状态", "type": "bool", "filter": "eq" }
    ]
*/
type Struct struct {
	Status int                          `json:"status"`
	Src    string                       `json:"src"`
	Op     map[string]map[string]string `json:"op"`
	//Struct []StructItem                 `json:"struct"`
	Items []map[string]interface{} `json:"struct"`
}

type StructItem struct {
	Key     string `json:"key"`
	Text    string `json:"text"`
	Filter  string `json:"filter"`
	Src     string `json:"src"`
	Type    string `json:"type"`
	Require bool   `json:"require"`
	Edit    bool   `json:"edit"`
	Hide    bool   `json:"hide"`
	Maxlen  int    `json:"maxlen"`
}

func (s *Struct) GetStruct(t *Table) {
	// src和op由调用方法填写
	var (
		items = []map[string]interface{}{}
		i     map[string]interface{}
	)

	var getType = func(tp string) string {
		switch tp {
		case BOOL:
			return "bool"
		case PASSWORD:
			return "password"
		case DATE, DATETIME:
			return "date"
		case INT, BIGINT:
			return "integer"
		case FLOAT, DOUBLE:
			return "number"
		}
		return "text"
	}

	var getRequire = func(tp string, nn bool) bool {
		return tp != BOOL && nn
	}

	for _, c := range t.Columns {
		i = map[string]interface{}{}
		i["key"] = c.Name
		i["text"] = c.Alias
		i["edit"] = c.EditAble
		i["type"] = getType(c.Type)
		if getRequire(c.Type, c.NotNull) {
			i["require"] = true
		}
		if c.Size > 0 {
			i["maxlen"] = c.Size
		}
		items = append(items, i)

		/*
			if c.Type == PASSWORD {
				i = map[string]interface{}{}
				i["key"] = "re_" + c.Type
				i["text"] = "重复" + c.Alias
				i["edit"] = c.EditAble
				i["type"] = getType(c.Type)
				i["hide"] = true
				i["equal"] = c.Name
				i["tip"] = fmt.Sprintf("必须和%s相等", c.Alias)
				items = append(items, i)
			}
		*/
	}
	s.Status = 1
	s.Items = items
}
