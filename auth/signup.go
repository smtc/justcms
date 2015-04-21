package auth

// 注册

import (
	"errors"
	"regexp"

	"github.com/jinzhu/gorm"
	"github.com/smtc/justcms/database"
)

const (
	minNameLength     = 6
	maxNameLength     = 32
	minPasswordLength = 7
)

var (
	hasLetters     = regexp.MustCompile(`[a-z]`) // 用户密码必选包含小写字母
	captainLetters = regexp.MustCompile(`[A-Z]`) // 用户密码必选包含大写字母
	// 无法使用在正则表达式中匹配.和-, 使用\.\-编译报错 2013-12-19
	// fix: 需要使用\\. \\-
	validNameChars = regexp.MustCompile("^[0-9a-zA-z_\\.\u4e00-\u9fa5]+$")
	validPhone     = regexp.MustCompile("^1[0-9]+$")

	MsisdnExist     = errors.New("msisdn has signup")          // 手机号码已注册
	MsisdnInvalid   = errors.New("msisdn invalid")             // 手机号码非法
	NameExist       = errors.New("name has exists")            // 用户名已存在
	NameTooShort    = errors.New("name invalid: too short")    // 用户名太短
	NameTooLong     = errors.New("name invalid: too long")     // 用户名太长
	NameCharInvalid = errors.New("name contains invalid char") // 用户名包含非法字符
	NameIsKept      = errors.New("name is kept")               // 用户名为保留字
	PasswdTooShort  = errors.New("password too short")         // 密码太短
	PasswdTooWeak   = errors.New("password not strong")        // 密码太弱
	PasswdNotMatch  = errors.New("password not match")         // 密码不匹配

	//密码必选包含至少一个小写字母
	PasswdNoChar = errors.New("password must contain at least one character")
	keptNames    = map[string]struct{}{
		"root":          struct{}{},
		"admin":         struct{}{},
		"administrator": struct{}{},
		"管理员":           struct{}{},
		"超级管理员":         struct{}{},
		"版主":            struct{}{},
		"超级版主":          struct{}{},
	}
)

// 注册的四个必选参数：
// msisdn:  手机号码
// name:    用户名
// passwd:  密码
// passwd2: 确认密码
func Signup(msisdn, name, passwd, passwd2 string) (user *User, err error) {
	// 检查密码
	if err = validPasswd(passwd, passwd2); err != nil {
		return
	}
	// 检查手机号码
	if err = ValidMsisdn(msisdn); err != nil {
		return
	}
	if err = HasMsisdnExist(msisdn); err != nil {
		return
	}
	// 检查用户名
	if err = ValidUsername(name); err != nil {
		return
	}

	return createUser(msisdn, name, passwd)
}

// 用户名是否合法
func ValidUsername(name string) (err error) {
	if len(name) < minNameLength {
		return NameTooShort
	}
	if len(name) > maxNameLength {
		return NameTooLong
	}
	if _, ok := keptNames[name]; ok {
		return NameIsKept
	}

	if validNameChars.MatchString(name) == false {
		return NameCharInvalid
	}
	// 检查名字是否在数据库中已经存在
	return HasNameExist(name)
}

// 检查名字是否在数据库中已经存在
func HasNameExist(name string) (err error) {
	var user User

	err = database.GetDB("").Where("name=?", name).Limit(1).Find(&user).Error
	if err == gorm.RecordNotFound {
		return nil
	} else if err == nil {
		return NameExist
	}

	return err
}

// 手机号码是否合法
func ValidMsisdn(msisdn string) (err error) {
	if len(msisdn) < 8 {
		return MsisdnInvalid
	}
	if validPhone.MatchString(msisdn) == false {
		return MsisdnInvalid
	}
	return nil
}

// 检查手机号码是否在数据库中已经存在
func HasMsisdnExist(msisdn string) (err error) {
	var user User

	err = database.GetDB("").Where("msisdn=?", msisdn).Limit(1).Find(&user).Error
	if err == gorm.RecordNotFound {
		return nil
	} else if err == nil {
		return MsisdnExist
	}

	return err
}

// 检查密码是否合法
func validPasswd(passwd, passwd2 string) (err error) {
	if passwd != passwd2 {
		return PasswdNotMatch
	}
	// 长度
	if len(passwd) < minPasswordLength {
		return PasswdTooShort
	}
	// 小写字母
	if hasLetters.MatchString(passwd) == false {
		return PasswdNoChar
	}
	return nil
}
