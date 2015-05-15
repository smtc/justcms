package main

import (
	"net/http"

	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/justcms/auth"
	"github.com/smtc/justcms/session"
	"github.com/zenazn/goji/web"
)

var (
	// 是否创建临时用户
	_createTmpUser bool
	// 临时用户默认超时时长，7天
	_tmpUserTimeout int
)

func init() {
	deferinit.AddInit(func() {
		_createTmpUser = config.GetBooleanDefault("CreateTmpUser", false)
		_tmpUserTimeout = config.GetIntDefault("TmpUserTimeout", 86400*7)
	}, nil, 10)
}

// 从web.C中获得session，如果没有，创建session，并保存在web.C中
func getSession(c web.C, r *http.Request) session.Session {
	if sess, ok := c.Env["session"]; ok {
		return sess.(session.Session)
	}

	// 更加request创建session
	sess := session.NewSession(r)
	if ok := sess.Init(); !ok {
		//println("session init failed!")
		sess.Create(_tmpUserTimeout, nil)
	}
	c.Env["session"] = sess

	return sess
}

func putSession(sess session.Session, w http.ResponseWriter) {
	sess.SetStore()
	sess.Save(w)
}

// 从session中获取用户，或者没有用户时，给其创建一个临时用户, 并设置cookie
func GetOrCreateUser(c web.C, w http.ResponseWriter, r *http.Request) (user *auth.User, err error) {
	sess := getSession(c, r)
	user, err = auth.GetUser(sess)
	if err == nil {
		return
	}
	// 不允许创建临时用户的情况下
	if _createTmpUser == false {
		return nil, err
	}

	//println("create tmp user ....")
	// 创建临时用户
	user = auth.CreateTmpUser()
	// 保存临时用户
	if err = auth.SaveTmpUser(user, _tmpUserTimeout); err != nil {
		println("save tmp user to redis failed", err.Error())
		return
	}

	auth.SaveUserSession(w, r, user, _tmpUserTimeout)
	return
}
