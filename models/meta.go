package models

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	//"github.com/smtc/goutils"
	"log"
	"regexp"
	"strings"
	"github.com/smtc/justcms/database"
)

// 本文有一个地方处理与wordpress不同，有可能造成sql语句查询结果错误，需要进一步验证，见parseMetaVar函数 ———— guotie 2014-09-24
// meta_typ is the table name, such as account, post, reply
type Meta struct {
	Id        int64
	SiteId    int64  `json:"site_id"`
	MetaTyp   string `sql:"size:64" json:"meta_typ"`
	ObjectId  string `sql:"size:64" json:"object_id"`
	MetaKey   string `sql:"size:300" json:"meta_key"`
	MetaValue string `sql:"type:text" json:"meta_value"`
}

// deal with xxx_meta table
// api:
//	 以下函数中的第一个参数id均不是Meta table中的id字段，而是target_id或object_id
//		AddMetaData
//		UpdateMetaData
//		DelMetaData
//		GetMetaData
//		HasMetaData
//
//		GetMetaDataById    -- has not implement
//		UpdateMetaDataById -- has not implement
//		DelMetaDataById    -- has not implement
//
//		UpdateMetaCache
//
//		getMetaSql
//
//		isProtectedMeta
//		sanitizeMeta
//		registerMeta

// Add metadata for the specified object.
//
func AddMetaData(sid string, typ, key, value string, override bool) (*Meta, error) {
	var (
		err   error
		meta  Meta
		count int
	)

	if key == "" {
		err = fmt.Errorf("key invalid")
		return nil, err
	}

	if typ == "" {
		return nil, fmt.Errorf("type invalid")
	}

	if sid == "" {
		err = fmt.Errorf("id is invalid")
		return nil, err
	}

	mvalue := sanitizeMeta(typ, key, value)
	// todo: apply_filter: "add_{{typ}}_metadata"
	db := database.GetDB("")
	if override == false {
		err = db.Table("metas").Where("meta_key=?", key).Where("object_id=%d", sid).Where("meta_typ=?", typ).Count(&count).Error

		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("meta key %s for %s has exist: %d", key, typ, count)
		}
	}
	// todo: do_action: "add_{{typ}}_meta"

	// ObjectId 是全局唯一
	err = db.Where("object_id=?", sid).Where("meta_key=?", key).Find(&meta).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil, err
	}
	meta.ObjectId = sid
	meta.MetaTyp = typ
	meta.MetaKey = key
	meta.MetaValue = mvalue
	if err = db.Save(&meta).Error; err != nil {
		return nil, err
	}

	return &meta, nil
}

// 更新meta
func UpdateMetaData(sid string, typ, key, value string, insert bool) (err error) {
	var meta Meta

	if sid == "" {
		err = fmt.Errorf("sid is invalid")
		return
	}
	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")

	// ObjectId 是全局唯一, 不需要type作为条件
	err = db.Where("object_id=?", sid).Where("meta_key=?", key).Find(&meta).Error
	if err != nil {
		if err == gorm.RecordNotFound && insert {
			meta.ObjectId = sid
			meta.MetaTyp = typ
			meta.MetaKey = key
			meta.MetaValue = value
			err = db.Save(&meta).Error
			return
		}
		// 不存在时, 返回错误
		return err
	}

	meta.MetaValue = value
	if err = db.Save(&meta).Error; err != nil {
		return err
	}

	return nil
}

// 删除一个meta
func DelMetaData(sid string, key string) (err error) {
	var meta Meta

	if key == "" {
		return fmt.Errorf("key should not empty")
	}

	if sid == "" {
		err = fmt.Errorf("sid is invalid")
		return
	}

	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")

	// ObjectId 是全局唯一, 不需要type作为条件
	err = db.Where("object_id=?", sid).Where("meta_key=?", key).Find(&meta).Error

	if err != nil {
		return err
	}
	err = db.Delete(&meta).Error
	return err
}

// get all metas of id & typ
func GetMetas(sid string, key string) (metas []Meta, err error) {
	if key == "" {
		err = fmt.Errorf("key invalid")
		return
	}

	if sid == "" {
		err = fmt.Errorf("sid is invalid")
		return
	}

	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")
	err = db.Where("object_id=?", sid).Find(&metas).Error

	return
}

// get meta
func GetMetaData(sid string, key string) (value string, err error) {
	var meta Meta
	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")

	// ObjectId 是全局唯一, 不需要type作为条件
	err = db.Where("object_id=?", sid).Where("meta_key=?", key).Find(&meta).Error

	if err != nil {
		return
	}
	value = meta.MetaValue

	return
}

func ScanMetaData(sid string, key string, value interface{}) error {
	var row *sql.Row

	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")
	// ObjectId 是全局唯一, 不需要type作为条件
	row = db.Model(Meta{}).Where("object_id=?", sid).Where("meta_key=?", key).Select("meta_value").Row()

	err := row.Scan(value)
	return err
}

// 是否存在meta key
func HasMetaData(sid string, key string) bool {
	var (
		err   error
		count int
	)
	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")

	// ObjectId 是全局唯一, 不需要type作为条件
	err = db.Table("metas").Where("object_id=?", sid).Where("meta_key=?", key).Count(&count).Error

	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}

	return false
}

func isProtectedMeta(key string, typ ...string) (res bool) {
	if key[0] == '_' {
		res = true
	}
	// todo: apply_filter
	return
}
func sanitizeMeta(typ, key, value string) string {
	// todo: apply_filter
	return value
}

// todo: RegisterMeta
func RegisterMeta() {

}

// get meta sql
// 目前仅考虑最简单的情况
/*
	opts:
	{
		"relation": xx, // "OR", "AND"
		"queries":[
		 	{
				"meta_key":     xx,
				"meta_cast_type":     xx,
				"meta_value":   xx,
				"meta_compare": xx
			},
			{
				"meta_key":xxx,
				"meta_value":xxx,
				"meta_compare":xx
			}
		],
		"onlyKeyQueries": ["key1", "key2"]
*/
// typ 用于查找属于哪个表的meta，例如account，post，reply等
// tn  主表名称，例如posts, accounts, replies
// id  字段，例如post_id, account_id, 或object_id
func getMetaSql(typ, tn string, sid string, opts map[string]interface{}) (qc queryClause, err error) {
	var (
		ok          bool
		joins       []string
		wheres      []string
		join, where string
	)

	typ = strings.TrimSpace(typ)
	tn = strings.TrimSpace(tn)
	if typ == "" || tn == "" {
		err = fmt.Errorf("getMetaSql: param typ & tn should NOT be empty.")
		return
	}
	queryOpt, err := parseMetaVar(opts)
	if err != nil {
		return
	}

	if sid == "" {
		err = fmt.Errorf("getMetaSql: param Id should not be empty.")
		return
	}

	onlyKeyQueries := queryOpt["onlyKeyQueries"].([]string)
	if len(onlyKeyQueries) > 0 {
		/*
			$join[]  = "INNER JOIN $meta_table ON $primary_table.$primary_id_column = $meta_table.$meta_id_column";

			foreach ( $key_only_queries as $key => $q )
				$where["key-only-$key"] = $wpdb->prepare( "$meta_table.meta_key = %s", trim( $q['key'] ) );
		*/

		// sid is unique, then typ is not needed
		join = fmt.Sprintf(" INNER JOIN metas ON %s.id = meta.object_id", tn)

		joins = append(joins, join)
		for _, k := range onlyKeyQueries {
			where = "meta.meta_key=" + k
			wheres = append(wheres, where)
		}
	}

	queries := queryOpt["queries"].([]map[string]interface{})
	for _, query := range queries {
		var (
			cnt                    int
			isArray                bool
			value, meta_values     string
			key, compare, castType string
		)
		if key, ok = query["meta_key"].(string); !ok {
			log.Println("getMetaSql: meta_key type should be string.")
			continue
		}
		value, isArray = sqlValue(query["meta_value"])
		compare = sqlCompare(query["meta_compare"], isArray)
		castType = getCastType(query["meta_cast_type"])

		alias := "metas"
		if cnt = len(joins); cnt != 0 {
			alias = "mt" + fmt.Sprint(cnt)
		}
		if "NOT EXISTS" == compare {
			join = "LEFT JOIN metas"
			if cnt != 0 {
				join += " AS " + alias
			}

			join += fmt.Sprintf(" ON (%s.object_id = %s.object_id AND %s.meta_key = '%s')",
				tn, alias, alias, key)

			joins = append(joins, join)
			wheres = append(wheres, " "+alias+".meta_id IS NULL")
			continue
		}

		join = "LEFT JOIN metas"
		if cnt != 0 {
			join += " AS " + alias
		}

		join += fmt.Sprintf(" ON (%s.object_id = %s.object_id)",
			tn, alias)

		joins = append(joins, join)

		// build where condition
		where = ""
		if key != "" {
			where = alias + ".meta_key = " + key
		}
		if compare == "IN" || compare == "NOT IN" {
			meta_values = "(" + value + ")"
		} else if compare == "BETWEEN" || compare == "NOT BETWEEN" {
			segs := strings.Split(value, ",")
			if len(segs) != 2 {
				err = fmt.Errorf("MetaSql: compare condition is BETWEEN, but the value is NOT 2")
				return
			}
			meta_values = segs[0] + " AND " + segs[1]
		} else if compare == "LIKE" || compare == "NOT LIKE" {
			meta_values = "%" + value + "%"
		} else {
			meta_values = value
		}
		if where != "" {
			where = where + " AND "
		}
		where = " (" + where + fmt.Sprintf("CAST(%s.meta_valueAS %s) %s %s)", alias, castType, compare, meta_values)

		wheres = append(wheres, where)
	}

	if len(wheres) == 0 {
		qc.where = ""
	} else {
		qc.where = " AND (" + strings.Join(wheres, "\n"+fmt.Sprint(queryOpt["relation"])) + " )"
	}
	qc.join = " " + strings.Join(joins, "\n")

	/*
		TODO: wordpress apply filter

		apply_filters_ref_array( 'get_meta_sql', array( compact( 'join', 'where' ),
					 $this->queries, $type, $primary_table, $primary_id_column, $context ) );
	*/
	return
}

// 把compare比较转换为string
// 如下几种情况：
//    1 nil
//    2 string
//    3 others
var _compares = map[string]string{"=": "=", "!=": "!=", ">": ">", ">=": ">=", "<": "<", "<=": "<=",
	"LIKE": "LIKE", "NOT LIKE": "NOT LIKE",
	"IN": "IN", "NOT IN": "NOT IN",
	"BETWEEN": "BETWEEN", "NOT BETWEEN": "NOT BETWEEN",
	"NOT EXISTS": "NOT EXISTS",
	"REGEXP":     "REGEXP", "NOT REGEXP": "NOT REGEXP", "RLIKE": "RLIKE"}

func sqlCompare(c interface{}, isArray bool) string {
	if c == nil {
		if isArray {
			return "IN"
		}
		return "="
	}
	if s, ok := c.(string); ok {
		s = strings.ToUpper(strings.TrimSpace(s))
		if _compares[s] != "" {
			return s
		}
		return "="
	}
	return "="
}

var castTypeReg = regexp.MustCompile("^(?:BINARY|CHAR|DATE|DATETIME|SIGNED|UNSIGNED|TIME|NUMERIC(?:\\(\\d+(?:,\\s?\\d+)?\\))?|DECIMAL(?:\\(\\d+(?:,\\s?\\d+)?\\))?)$")

func getCastType(i interface{}) string {
	s, ok := i.(string)
	if !ok {
		return "CHAR"
	}
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return "CHAR"
	}
	//if ( ! preg_match( '/^(?:BINARY|CHAR|DATE|DATETIME|SIGNED|UNSIGNED|TIME|NUMERIC(?:\(\d+(?:,\s?\d+)?\))?|DECIMAL(?:\(\d+(?:,\s?\d+)?\))?)$/', $meta_type ) )
	//	return 'CHAR';

	if "NUMERIC" == s {
		return "SIGNED"
	}
	return s
}

// 格式化opts参数, 传入getMetaSql
//  opts["meta_relation"] 覆盖默认的relation(AND)
func parseMetaVar(opts map[string]interface{}) (metaOpt map[string]interface{}, err error) {
	var query = map[string]interface{}{}

	metaOpt["relation"] = "AND"
	if opts["meta_relation"] != nil {
		if rel, ok := opts["meta_relation"].(string); ok {
			metaOpt["relation"] = rel
		} else {
			log.Println("parseMetaVar: opts[\"meta_relation\"] type is NOT string, use default \"AND\" relation")
		}
	}

	queries := make([]map[string]interface{}, 0)
	onlyKeyQueries := make([]string, 0)

	if opts["meta_key"] != nil {
		query["meta_key"] = opts["meta_key"]
	}
	if opts["meta_compare"] != nil {
		query["meta_compare"] = opts["meta_compare"]
	}
	if opts["meta_value"] != nil {
		query["meta_value"] = opts["meta_value"]
	}
	if opts["meta_cast_type"] != nil {
		query["meta_cast_type"] = opts["meta_cast_type"]
	}

	// FIX Me!!!
	// 这里与wordpress meta.php中的处理不同，wordprss只有当query["value"]为数组且为空数组时，才将query放入到onlyKeyQuery中
	// 而query["value"]为空字符串时，还是在query数组中处理，但后面的似乎是一样的

	if isEmpty(query["meta_value"]) {
		if s, ok := query["meta_key"].(string); ok {
			s = strings.TrimSpace(s)
			if s != "" {
				onlyKeyQueries = append(onlyKeyQueries, s)
			} else {
				log.Println("parseMetaVar: meta_key should NOT empty when meta_value is empty.")
			}
		} else {
			log.Println("parseMetaVar: meta_key should be string")
		}
	} else {
		queries = append(queries, query)
	}

	if opts["meta_query"] != nil {
		var (
			ok bool
			q  map[string]interface{}
			qa []map[string]interface{}
		)
		// 有可能是map[string]interface{}类型或[]map[string]interface{}类型
		if q, ok = opts["meta_query"].(map[string]interface{}); ok {
			qa = append(qa, q)
		} else {
			qa, ok = opts["meta_query"].([]map[string]interface{})
		}

		if ok {
			for _, q := range qa {
				if isEmpty(q["meta_value"]) {
					if s, ok := query["meta_key"].(string); ok {
						s = strings.TrimSpace(s)
						if s != "" {
							onlyKeyQueries = append(onlyKeyQueries, s)
						} else {
							log.Println("parseMetaVar: meta_key should NOT empty when meta_value is empty.")
						}
					} else {
						log.Println("parseMetaVar: meta_key should be string")
					}
				} else {
					queries = append(queries, q)
				}
			}
		} else {
			// 非法的meta_query
			log.Println("parseMetaVar: opts[\"meta_query\"] type is NOT map[string]interface{} or []map[string]interface{}, ignored")
		}
	}

	if len(queries) == 0 && len(onlyKeyQueries) == 0 {
		err = fmt.Errorf("parseMetaVar: no valid queries")
		return
	}

	metaOpt["queries"] = queries
	metaOpt["onlyKeyQueries"] = onlyKeyQueries

	return
}
