package auth

import (
	"errors"
	"net/http"

	"github.com/smtc/justcms/database"
	"github.com/smtc/justcms/session"
)

var (
	PasswdWrong = errors.New("Password wrong or User not exist.")
)

// 登陆
// msisdn:  手机号码
// passwd:  密码
func Signin(msisdn, passwd string) (*User, error) {
	var (
		err  error
		user User
	)

	if err = ValidMsisdn(msisdn); err != nil {
		return nil, err
	}
	err = database.GetDB("").Where("msisdn=?", msisdn).Limit(1).Find(&user).Error
	if err != nil {
		return nil, PasswdWrong
	}

	if checkpasswd(user.Password, passwd) == false {
		return nil, PasswdWrong
	}

	return &user, nil
}

// 保存用户的uid到session中
// secs 为0时，将使用默认值 86400 * 7
func SaveUserSession(w http.ResponseWriter, r *http.Request, user *User, secs int) {
	sess := session.NewSession(r)
	if secs == 0 {
		secs = 86400 * 7
	}
	sess.Create(secs, nil)
	sess.SetKey("uid", user.Id)
	sess.SetKey("user_objid", user.ObjectId)

	println(sess.CookieValue())
	sess.SetStore()
	// save session to client
	sess.Save(w)
}
