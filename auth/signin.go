package auth

import (
	"errors"

	"github.com/smtc/justcms/database"
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
