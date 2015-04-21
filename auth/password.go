package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
)

var (
	_all_ch = []byte(`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`)
	_secret = []byte(`e9DuuAVpK1OuBKHOpv0C0mFt9xfb3AlReuyGbetkXBtb2Ukz`)
	_ch_len = len(_all_ch)
)

const (
	_PASS_TOKEN_LEN = 12
)

func createtoken(length int) []byte {
	var res []byte

	for i := 0; i < length; i++ {
		res = append(res, _all_ch[rand.Intn(_ch_len)])
	}

	return res
}

// 返回值：长度120
func _createpasswd(salt, passwd []byte) string {
	var orig []byte

	orig = append(salt, passwd...)
	orig = append(orig, _secret...)
	hsh := sha512.Sum384(orig)

	return hex.EncodeToString(append(salt, hsh[:]...))
}

func createpasswd(passwd string) string {
	salt := createtoken(_PASS_TOKEN_LEN)
	return _createpasswd(salt, []byte(passwd))
}

func checkpasswd(passwd string, input string) bool {
	res, err := hex.DecodeString(passwd)
	if err != nil {
		panic(err.Error())
	}
	salt := res[0:_PASS_TOKEN_LEN]

	pwinput := _createpasswd(salt, []byte(input))
	//println(salt, pwinput)
	return passwd == pwinput
}
