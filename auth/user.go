package auth

import (
	"encoding/json"
	"errors"
	"time"

	//"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
	"github.com/smtc/justcms/session"
	"github.com/smtc/justcms/utils"
)

const KEY_TMP_USERS = "tmp_users"

var (
	NoUserInCookie = errors.New("Not found user id in cookie.")
	VisitUser      = &User{}
)

// 账号管理

type User struct {
	Id         int64  `json:"id"`
	SiteId     int64  `json:"site_id"`
	ObjectId   string `sql:"size:64" json:"object_id"`
	Name       string `sql:"size:40" json:"name"`
	Email      string `sql:"size:100" json:"email"`
	Avatar     string `sql:"size:120" json:"avatar"`
	Msisdn     string `sql:"size:20" json:"msisdn"`
	Password   string `sql:"size:200" json:"password"`
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
	SilverCoins   int64 `json:"silver_coins"` // 银币
	Coppercoins   int64 `json:"copper_coins"` // 铜币
	Reputation    int   `json:"reputation"`   // 声望值
	Credits       int   `json:"credits"`      // 信用等级
	Experience    int   `json:"experience"`   // 经验值
	Activation    int   `json:"activation"`   // 活跃度
	Redcards      int   `json:"red_cards"`
	Yellowcards   int   `json:"yellow_cards"`
	Notifications int   `json:"notifications"`
	Messages      int   `json:"messages"`

	// user capability
	Capability *UserCap `json:"capability" sql:"-"`
	// other meta data
	metaData map[string]interface{} `json:"-" sql:"-"`
}

//
// 创建用户
func createUser(msisdn, name, passwd string) (*User, error) {
	objId := utils.ObjectId()
	user := User{
		ObjectId:  objId,
		Msisdn:    msisdn,
		Name:      name,
		Password:  createpasswd(passwd),
		CreatedAt: time.Now(),
	}

	err := database.GetDB("").Save(&user).Error
	return &user, err
}

// 创建临时用户
func CreateTmpUser() *User {
	return &User{
		Id:        -1,
		SiteId:    0,
		ObjectId:  utils.ObjectId(),
		CreatedAt: time.Now(),
		LastLogin: time.Now(),
	}
}

// 将临时用户保存在redis中
func SaveTmpUser(u *User, secs int) error {
	if secs == 0 {
		secs = 86400 * 7
	}
	buf, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = database.SET(u.ObjectId, buf)

	return err
}

// 获取用户:
//    当Id为-1时，直接从redis中获取用户；
//    当Id>0时，从数据库中获取用户
func GetUser(sess session.Session) (user *User, err error) {
	// 从cookie中获取用户id
	uid, ok := sess.Get("uid").(int64)
	if !ok {
		return nil, NoUserInCookie
	}

	user = &User{}
	if uid > 0 {
		// 从数据库中获取用户
		user.Id = uid
		err = database.GetDB("").Find(user).Error

		return user, err
	}

	// 从redis中获取用户
	objId, ok := sess.Get("user_objid").(string)
	if !ok {
		return nil, NoUserInCookie
	}
	val, err := database.GET(objId)
	if err != nil {
		return
	}
	err = json.Unmarshal(val, user)
	return
}
