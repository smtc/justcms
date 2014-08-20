package models

import (
	//"github.com/smtc//utils"
	"time"

	"github.com/smtc/JustCms/database"
)

// 账号管理

type Account struct {
	Id        int64
	Name      string `sql:"size:40"`
	Email     string `sql:"size:100"`
	Avatar    string `sql:"size:120"`
	Msisdn    string `sql:"size:20"`
	Telephone string `sql:"size:20"`
	Password  string `sql:"size:80"`
	Roles     string `sql:"type:text"` // 这是一个string数组, 以,分割
	City      string `sql:"size:40"`
	MainUser  bool   // 是否是主用户
	MainId    int64  // 如果不是主用户, 主用户id；否则为0
	Approved  bool
	Activing  bool
	ApproveId string `sql:"size:20"`
	IpAddr    string `sql:"size:30"`
	DaysLogin int

	Birthday time.Time

	BannedAt    int64
	BannedTill  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLoginAt int64
	LastPostAt  int64
	ApprovedAt  int64

	Goldcoins     int64 // 金币
	Coppercoins   int64 // 铜币
	Reputation    int   // 声望值
	Credits       int   // 信用等级
	Experience    int   // 经验值
	Activation    int   // 活跃度
	Redcards      int
	Yellowcards   int
	Notifications int
	Messages      int
}

type AccountMeta struct {
	Id     int64
	Aid    int64
	Plugin string `sql:"size:180"`
	//Meta   utils.SQLGenericMap `sql:"type:text"`
}

var (
	account_db = ""
)

func AccountSave(acct *Account) error {
	return nil
}

func AccountGet(id int64) (Account, error) {
	return Account{}, nil
}

func AccountList(page, size int, filter map[string]interface{}) ([]Account, error) {
	db := database.GetDB(account_db)
	var accts []Account

	err := db.Offset(page * size).Limit(size).Find(&accts).Error
	return accts, err
}

func AccountDelete(id int64) error {
	return nil
}
