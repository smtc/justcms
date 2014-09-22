package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smtc/goutils"
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
	Id        int64
	AccountId int64  `json:"account_id"`
	ObjectId  string `json:"object_id"`
	MetaKey   string `sql:"size:300" json:"meta_key"`
	MetaValue string `sql:"size:100000" json:"meta_value"`
}

func getAccountDB() *gorm.DB {
	return GetDB(ACCOUNT_DB)
}

func (a *Account) Get(id int64) error {
	db := getAccountDB()
	return db.First(a, id).Error
}

func (a *Account) Save() error {
	db := getAccountDB()
	return db.Save(a).Error
}

func (a *Account) Delete() error {
	db := getAccountDB()
	return db.Delete(a).Error
}

func AccountDelete(where string) {
	db := getAccountDB()
	db.Where(where).Delete(&Account{})
}

func AccountList(page, size int, filter *map[string]interface{}) ([]Account, error) {
	db := getAccountDB()
	var accts []Account

	err := db.Offset(page * size).Limit(size).Find(&accts).Error
	return accts, err
}

// get account id by name
func getAuthorIdByName(name string) (id int64, err error) {
	var acct Account
	db := getAccountDB()
	err = db.Where("name=?", name).First(&acct).Error
	id = acct.Id
	return
}

// author, author__in, author__not_in
func getAuthorSql(opt map[string]interface{}) (qc queryClause, err error) {
	if opt["author"] != nil {
		aid := opt["author"].(int64)
		qc.where = fmt.Sprintf(" AND posts.post_author = %d ", aid)
	} else if opt["author__in"] != nil {
		aid := goutils.ToString(opt["author__in"].([]int64), "")
		qc.where = fmt.Sprintf(" AND posts.post_author IN (%s) ", aid)
	} else if opt["author__not_in"] != nil {
		//$where .= " AND {$wpdb->posts}.post_author NOT IN ($author__not_in) ";
		aid := goutils.ToString(opt["author__not_in"].([]int64), "")
		qc.where = fmt.Sprintf(" AND posts.post_author NOT IN (%s) ", aid)
	}
	return
}
