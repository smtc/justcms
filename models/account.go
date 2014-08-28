package models

import (
	"time"

	"github.com/smtc/JustCms/database"
)

// 账号管理

type Account struct {
	Id         int64  `json:"id"`
	ObjectId   string `sql:"size:64" json:"object_id"`
	Name       string `sql:"size:40" json:"name"`
	Email      string `sql:"size:100" json:"email"`
	Avatar     string `sql:"size:120" json:"avatar"`
	Msisdn     string `sql:"size:20" json:"msisdn"`
	Password   string `sql:"size:80" json:"password"`
	Roles      string `sql:"type:text" json:"roles"` // 这是一个string数组, 以,分割
	City       string `sql:"size:40" json:"city"`
	MainUser   bool   `json:"main_user"` // 是否是主用户
	MainId     int64  `json:"main_id"`   // 如果不是主用户, 主用户id；否则为0
	Approved   bool   `json:"approved"`
	Activing   bool   `json:"acitiving"`
	ApprovedBy string `sql:"size:20" json:"approved_by"`
	IpAddr     string `sql:"size:30" json:"ipaddr"`
	DaysLogin  int    `json:"days_login"`

	Birthday   time.Time `json:"birthday"`
	BannedAt   time.Time `json:"banned_at"`
	BannedTill time.Time `json:"banned_till"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ApprovedAt time.Time `json:"approved_at"`
	LastLogin  time.Time `json:"last_login"`
	LastPost   time.Time `json:"last_post"`

	Goldcoins     int64 `json:"gold_coins"`   // 金币
	Coppercoins   int64 `json:"copper_coins"` // 铜币
	Reputation    int   `json:"reputation"`   // 声望值
	Credits       int   `json:"credits"`      // 信用等级
	Experience    int   `json:"experience"`   // 经验值
	Activation    int   `json:"activation"`   // 活跃度
	Redcards      int   `json:"red_cards"`
	Yellowcards   int   `json:"yellow_cards"`
	Notifications int   `json:"notifications"`
	Messages      int   `json:"messages"`
}

type AccountMeta struct {
	Id           int64
	Aid          int64
	PostObjectId string `sql:"size:64" json:"object_id"`
	MetaKey      string `sql:"size:300" json:"meta_key"`
	MetaValue    string `sql:"size:100000" json:"meta_value"`
	Plugin       string `sql:"size:180"`
	//Meta   utils.SQLGenericMap `sql:"type:text"`
}

func (a *Account) Save() error {
	db := database.GetDB(account_db)
	err := db.Save(a).Error
	return err
}

func AccountGet(id int64) (*Account, error) {
	db := database.GetDB(account_db)
	acct := &Account{Id: id}
	err := db.First(acct).Error
	return acct, err
}

func AccountList(page, size int, filter map[string]interface{}) ([]Account, error) {
	db := database.GetDB(account_db)
	var accts []Account

	err := db.Offset(page * size).Limit(size).Find(&accts).Error
	return accts, err
}

func (a *Account) Delete() error {
	db := database.GetDB(account_db)
	return db.Delete(a).Error
}
