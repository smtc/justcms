package utils

import (
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

// 生成universal distribute id
var (
	pid        uint64 = 0
	machineid  uint64 = 0
	atmid      uint32
	randString = `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	randLen    = len(randString)
	tmBegin    = time.Date(2010, 2, 9, 1, 0, 0, 0, time.Local).Unix()
)

func init() {
	_genmachineid()
	_p2()
}

// ObjectID —— 类似于monogodb中的objectID
//
// 由一个int32和一个int64编码为36进制
// 0xc0000000编码为36进制: 1h9u1hc
// 0xffffffff编码为36进制: 1z141z3
//
// 0xf000000000000000编码为36进制: 3w5e11264sgsf
// 0xffffffffffffffff编码为36进制: 3w5e11264sgsf
//
// 第一部分int32的最高2位固定为1, 则该int32的36位编码长度固定位7位
// 前31位: 时间, 当前秒数-tmBegin
//
// 第二部门为int64构成如下：
// atomicID: 20位
// machineID: 24位
// pid: 20位
func ObjectId() string {
	return _tmid() + _p2()
}

// 返回random string, 长度为l
func RandomString(l int) string {
	var res string

	rand.Seed(time.Now().Unix())
	for i := 0; i < l; i++ {
		res += string(randString[rand.Intn(randLen)])
	}

	return res
}

func _tmid() string {
	now := time.Now().Unix() - tmBegin
	now = (now & 0xc0000000) + 0xc0000000
	return strconv.FormatUint(uint64(now), 36)
}

func _p2() string {
	var res uint64

	res = (uint64(atmid) << 44) + machineid + _genpid()
	atmid = atomic.AddUint32(&atmid, 1)
	if atmid >= 0xfffff {
		atmid = 0
	}

	return strconv.FormatUint(res, 36)
}

// 20 bit
func _genpid() uint64 {
	if pid == 0 {
		id := os.Getpid()
		id = ((id & 0xfffff) + (id >> 20)) & 0xfffff

		pid = uint64(id)
	}
	return pid
}

// 24位
func _genmachineid() {
	var s string

	ifs, err := net.Interfaces()
	if err == nil {
		for _, i := range ifs {
			s += i.HardwareAddr.String()
		}
	} else {
		s = err.Error()
	}

	//code := hashcode(s)
	var code uint32 = 1

	for i := 0; i < len(s); i++ {
		code += (code * 37) + uint32(s[i])
	}

	code = ((code & 0xffffff) + (code >> 24)) & 0xffffff

	machineid = uint64(code)
}
