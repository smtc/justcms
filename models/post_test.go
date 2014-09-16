package models

import (
	"github.com/smtc/justcms/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhere(t *testing.T) {
	var posts []Post

	db := database.GetDB("")
	if err := db.Debug().Where(`1=1 
AND term_relations.term_id IN (2,3,17)
AND posts.post_type = 'post'
AND (posts.post_status = 'publish'
OR posts.post_status = 'private')
GROUP BY posts.id
ORDER BY posts.publish_at DESC
LIMIT 0, 10`).Find(&posts).Error; err != nil {
		assert.Nil(t, err)
	}
}

func TestConvert(t *testing.T) {
	testconvertInt(t)
	testconverbool(t)
	testsplitIds(t)
	testsplitSa(t)
}

func testconvertInt(t *testing.T) {
	var c = map[string]int{
		//"-4096.00": -4096,
		"-245": -245,
		"-13":  -13,
		"-1":   -1,
		"0":    0,
		"1":    1,
		"12":   12,
		"256":  256,
		"8000": 8000,
	}
	for name, value := range c {
		i, err := convertInt("convertInt", name)
		assert.Nil(t, err)
		assert.Equal(t, i, value, "should be equal.")
	}
}

func testconvertInt64(t *testing.T) {

}

func testconverbool(t *testing.T) {
	var c = map[string]bool{
		"false":  false,
		"true":   true,
		"False":  false,
		"FALse":  false,
		"True":   true,
		"TRue":   true,
		"others": false,
	}
	for name, value := range c {
		b, err := convertBool("convertBool", name)
		if name != "others" {
			assert.Nil(t, err)
			assert.Equal(t, b, value, "should be equeal")
		} else {
			assert.NotNil(t, err)
		}
	}
}

func testsplitIds(t *testing.T) {
	type result struct {
		ia  []int64
		ea  []int64
		err error
	}
	var c = map[string]result{
		"1,2,3": result{[]int64{1, 2, 3}, []int64{}, nil},
		"1, 2 ,	3 ": result{[]int64{1, 2, 3}, []int64{}, nil},
		"1, -2 ,	3 ": result{[]int64{1, 3}, []int64{2}, nil},
		//"0, -1, -3, -4": result{[]int64{}, []int64{1, 3, 4}, nil},
	}
	for name, value := range c {
		ia, ea, err := splitIds("splitIds", name)
		assert.Nil(t, err)
		assert.Equal(t, ia, value.ia, "include array should be equal")
		assert.Equal(t, ea, value.ea, "exclude array should be equal")
	}
}

func testsplitSa(t *testing.T) {
	type result struct {
		sa  []string
		rel int
		err error
	}
	var c = map[string]result{
		"aa,bb,cc": result{[]string{"aa", "bb", "cc"}, 0, nil},
		"	aa	,		bb   , cc": result{[]string{"aa", "bb", "cc"}, 0, nil},
		"aa+bb+cc": result{[]string{"aa", "bb", "cc"}, 1, nil},
	}
	for name, value := range c {
		sa, rel, err := splitSa("splitSa", name)
		assert.Nil(t, err)
		assert.Equal(t, sa, value.sa, "string array should be equal")
		assert.Equal(t, rel, value.rel, "rel should be equal")
	}
}

func testparseQuery(t *testing.T) {

}
