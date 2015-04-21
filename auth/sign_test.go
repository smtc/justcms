package auth

import (
	"testing"
)

func TestSign(t *testing.T) {
	testSignup(t)
	testSignin(t)
}

func assert(t *testing.T, res bool, errmsg string) {
	if res == false {
		t.Error(errmsg)
	}
}

// 测试注册用户
func testSignup(t *testing.T) {
	var err error

	// password case 1
	_, err = Signup("1234", "guotie", "abcd", "efg")
	assert(t, err == PasswdNotMatch, "error should be PasswdNotMatch\n")

	// password case 2
	_, err = Signup("1234", "guotie", "abcd", "abcd")
	assert(t, err == PasswdTooShort, "error should be PasswdTooShort\n")

	// password case 3
	_, err = Signup("1234", "guotie", "ABCDEFGH", "ABCDEFGH")
	assert(t, err == PasswdNoChar, "error should be PasswdNoChar")
	_, err = Signup("1234", "guotie", "12345678", "12345678")
	assert(t, err == PasswdNoChar, "error should be PasswdNoChar")

	// msisdn case 1
	_, err = Signup("1234", "guotie", "a12345678", "a12345678")
	assert(t, err == MsisdnInvalid, "error should be MsisdnInvalid")

	// msisdn case 2
	_, err = Signup("32345678", "guotie", "a12345678", "a12345678")
	assert(t, err == MsisdnInvalid, "error should be MsisdnInvalid")

	// name case 1
	_, err = Signup("15612345678", "铁ge", "a12345678", "a12345678")
	assert(t, err == NameTooShort, "error should be NameTooShort")

	// name case 2
	_, err = Signup("15612345678", "版主", "a12345678", "a12345678")
	assert(t, err == NameIsKept, "error should be NameIsKept")

	// name case 3
	_, err = Signup("15612345678", "ro--ot", "a12345678", "a12345678")
	assert(t, err == NameCharInvalid, "error should be NameCharInvalid")

	// name case 4
	_, err = Signup("15612345678", "guotie", "a12345678", "a12345678")
	assert(t, err == nil, "error should be nil")

	// name case 5
	_, err = Signup("15612345678", "guotie", "a12345678", "a12345678")
	assert(t, err == MsisdnExist, "error should be MsisdnExist")
	// name case 6
	_, err = Signup("1561234567", "guotie", "a12345678", "a12345678")
	assert(t, err == NameExist, "error should be NameExist")
	//println(err.Error())
}

func testSignin(t *testing.T) {
	_, err := Signin("15612345678", "1234")
	assert(t, err == PasswdWrong, "error should be PasswdWrong")
	//println(err.Error())

	_, err = Signin("15612345678", "a12345678")
	assert(t, err == nil, "error should be nil")

	_, err = Signin("15612344321", "abcd")
	assert(t, err == PasswdWrong, "error should be PasswdWrong")
	//println(err.Error())
}
