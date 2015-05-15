package database

// 初始化资源

import (
	"github.com/guotie/deferinit"
)

func init() {
	deferinit.AddInit(connectDatabases, nil, 1000)
}

// 连接数据库
// 连接redis
func connectDatabases() {
	OpenDefaultDB()
	createRedisPool("")
}
