package skiplist

import (
	"testing"

	"github.com/go-playground/assert"
)

func TestSkipList(t *testing.T) {
	insertData := map[int64]int64{
		1:   1,
		5:   5,
		3:   3,
		6:   6,
		9:   9,
		111: 111,
		123: 123,
		10:  10,
		37:  37,
		22:  22,
		25:  25,
		28:  28,
	}

	deleteData := map[int64]int64{
		1:   1,
		5:   5,
		3:   3,
		6:   6,
		9:   9,
		111: 111,
		123: 123,
		10:  10,
		37:  37,
		22:  22,
		25:  25,
	}

	updateData := map[int64]int64{
		1:   11,
		5:   15,
		3:   13,
		6:   16,
		9:   19,
		111: 1111,
		123: 1123,
		10:  110,
		37:  137,
		22:  122,
		25:  125,
		28:  128,
		66:  1,
		77:  2,
		88:  3,
	}

	sl := New(32, 0.5)
	for k, v := range insertData {
		sl.Insert(k, v)
	}
	assert.Equal(t, sl.length, len(insertData))
	for k, v := range insertData {
		assert.Equal(t, sl.Find(k), v)
	}

	ll := len(insertData)
	for k := range deleteData {
		sl.Delete(k)
		assert.Equal(t, sl.Find(k), nil)
		ll--
		assert.Equal(t, sl.length, ll)
	}
	assert.Equal(t, sl.Find(28), int64(28))

	// 删除剩余1个元素之后重新插入
	for k, v := range insertData {
		sl.Insert(k, v)
	}
	assert.Equal(t, sl.length, len(insertData))
	for k, v := range insertData {
		assert.Equal(t, sl.Find(k), v)
	}

	// 更新测试
	for k, v := range updateData {
		sl.Insert(k, v)
	}
	assert.Equal(t, sl.length, len(updateData))
	for k, v := range updateData {
		assert.Equal(t, sl.Find(k), v)
	}
}
