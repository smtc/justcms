package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/smtc/goutils"
	"github.com/smtc/justcms/database"
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

	if iid <= 0 && sid == "" {
		err = fmt.Errorf("id is invalid")
		return nil, err
	}

	mvalue := sanitizeMeta(typ, key, value)
	// todo: apply_filter: "add_{{typ}}_metadata"
	db := database.GetDB("")
	if override {
		if iid > 0 {
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

	if iid > 0 {
		err = db.Where("target_id=?", iid).Where("meta_typ=?", typ).Where("meta_key=?", key).Find(&meta).Error
	} else {
		// ObjectId 是全局唯一
		err = db.Where("objetc_id=?", sid).Where("meta_key=?", key).Find(&meta).Error

	}
	if err != nil && err != gorm.RecordNotFound {
		return nil, err
	}
	meta.TargetId = iid
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
func UpdateMetaData(id interface{}, typ, key, value string) (err error) {
	var (
		iid  = goutils.ToInt64(id, 0)
		sid  = goutils.ToString(id, "")
		meta Meta
	)
	if typ == "" || key == "" {
		err = fmt.Errorf("type or key invalid")
		return
	}

	if iid <= 0 && sid == "" {
		err = fmt.Errorf("id is invalid")
		return
	}
	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")
	if iid > 0 {
		err = db.Where("target_id=?", iid).Where("meta_typ=?", typ).Where("meta_key=?", key).Find(&meta).Error
	} else {
		// ObjectId 是全局唯一, 不需要type作为条件
		err = db.Where("objetc_id=?", sid).Where("meta_key=?", key).Find(&meta).Error
	}
	if err != nil {
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
func DelMetaData(id interface{}, typ, key string) (err error) {
	var (
		iid  = goutils.ToInt64(id, 0)
		sid  = goutils.ToString(id, "")
		meta Meta
	)
	if typ == "" || key == "" {
		err = fmt.Errorf("type or key invalid")
		return
	}

	if iid <= 0 && sid == "" {
		err = fmt.Errorf("id is invalid")
		return
	}

	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")
	if iid > 0 {
		err = db.Where("target_id=?", iid).Where("meta_typ=?", typ).Where("meta_key=?", key).Find(&meta).Error
	} else {
		// ObjectId 是全局唯一, 不需要type作为条件
		err = db.Where("objetc_id=?", sid).Where("meta_key=?", key).Find(&meta).Error
	}
	if err != nil {
		return err
	}
	err = db.Delete(&meta).Error
	return err
}

// get all metas of id & typ
func GetMetas(id interface{}, typ string) (metas []Meta, err error) {
	var (
		iid = goutils.ToInt64(id, 0)
		sid = goutils.ToString(id, "")
	)
	if typ == "" {
		err = fmt.Errorf("type or key invalid")
		return
	}

	if iid <= 0 && sid == "" {
		err = fmt.Errorf("id is invalid")
		return
	}

	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")
	if iid > 0 {
		err = db.Where("meta_typ=?", typ).Where("target_id=?", iid).Find(&metas).Error
	} else {
		err = db.Where("meta_typ=?", typ).Where("object_id=?", sid).Find(&metas).Error
	}

	return
}

// get meta
func GetMetaData(id interface{}, typ, key string) (value string, err error) {
	var (
		iid  = goutils.ToInt64(id, 0)
		sid  = goutils.ToString(id, "")
		meta Meta
	)
	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")
	if iid > 0 {
		err = db.Where("target_id=?", iid).Where("meta_typ=?", typ).Where("meta_key=?", key).Find(&meta).Error
	} else {
		// ObjectId 是全局唯一, 不需要type作为条件
		err = db.Where("objetc_id=?", sid).Where("meta_key=?", key).Find(&meta).Error
	}

	if err != nil {
		return
	}
	value = meta.MetaValue

	return
}

// 是否存在meta key
func HasMetaData(id interface{}, typ, key string) bool {
	var (
		iid   = goutils.ToInt64(id, 0)
		sid   = goutils.ToString(id, "")
		err   error
		count int
	)
	// todo: apply_filter, do_action, etc...

	db := database.GetDB("")
	if iid > 0 {
		err = db.Table("metas").Where("target_id=?", iid).Where("meta_typ=?", typ).Where("meta_key=?", key).Count(&count).Error
	} else {
		// ObjectId 是全局唯一, 不需要type作为条件
		err = db.Table("metas").Where("objetc_id=?", sid).Where("meta_key=?", key).Count(&count).Error
	}
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
func getMetaSql() {

}
