package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/smtc/goutils"
	"github.com/smtc/justcms/database"
)

var (
	_metaTables = map[string]string{
		"user":    "account_metas",
		"accout":  "account_metas",
		"post":    "post_metas",
		"reply":   "reply_metas",
		"comment": "comment_metas",
	}
)

// meta_typ is the table name, such as account, post, reply
type Meta struct {
	Id        int64
	MetaTyp   string `sql:"size:64" json:"meta_typ"`
	TargetId  int64  `json:"target_id"`
	ObjectId  string `sql:"size:64" json:"object_id"`
	MetaKey   string `sql:"size:300" json:"meta_key"`
	MetaValue string `sql:"size:100000" json:"meta_value"`
}

// deal with xxx_meta table
// api:
//		AddMetaData
//		UpdateMetaData
//		DelMetaData
//		GetMetaData
//		HasMetaData
//
//		GetMetaDataById
//		UpdateMetaDataById
//		DelMetaDataById
//
//		UpdateMetaCache
//
//		getMetaSql
//		getMetaTable
//
//		isProtectedMeta
//		sanitizeMeta
//		registerMeta

// Add metadata for the specified object.
//
func AddMetaData(id interface{}, typ, key, value string, override bool) (*Meta, error) {
	var (
		iid   = goutils.ToInt64(id, 0)
		sid   = goutils.ToString(id, "")
		err   error
		meta  Meta
		count int
	)
	if typ == "" || key == "" {
		err = fmt.Errorf("type or key invalid")
		return nil, err
	}

	if iid == 0 && sid == "" {
		err = fmt.Errorf("id is invalid")
		return nil, err
	}

	mvalue := sanitizeMeta(typ, key, value)
	// todo: apply_filter: "add_{{typ}}_metadata"
	db := database.GetDB("")
	if override {
		if iid != 0 {
			err = db.Table("metas").Where("meta_key=?", key).Where("target_id=?", iid).Where("meta_typ=?", typ).Count(&count).Error
		} else {
			err = db.Table("metas").Where("meta_key=?", key).Where("object_id=%d", sid).Where("meta_typ=?", typ).Count(&count).Error
		}
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("meta key %s for %s %d has exist: %d", key, typ, id, count)
		}
	}
	// todo: do_action: "add_{{typ}}_meta"

	if iid != 0 {
		meta.TargetId = iid
		meta.MetaTyp = typ
		meta.MetaKey = key
		meta.MetaValue = mvalue
	} else {
		// ObjectId 是全局唯一
		if err = db.Where("objetc_id=?", sid).Where("meta_key=?", key).Find(&meta).Error; err != nil && err != gorm.RecordNotFound {
			return nil, err
		}
		meta.ObjectId = sid
		meta.MetaTyp = typ
		meta.MetaKey = key
		meta.MetaValue = mvalue
	}
	if err = db.Save(&meta).Error; err != nil {
		return nil, err
	}

	return &meta, nil
}

func getMetaTable(typ string) string {
	return _metaTables[typ]
}

func sanitizeMeta(typ, key, value string) string {
	return value
}
