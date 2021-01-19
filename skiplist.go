package skiplist

import (
	"math/rand"
	"time"
)

type Node struct {
	// 跳表排序key
	key int64
	// 跳表节点存储value
	value interface{}
	// 向前链表
	forward []*Node
}

type SkipList struct {
	// 跳表设置最高层级
	maxLevel int
	// the P value for the SkipList
	probability float32

	// 跳表头节点
	header *Node
	// 跳表当前层级, 初始化level=1
	level int
	// 跳表底层链表长度
	length int
}

// level 从1开始
func New(maxLevel int, probability float32) *SkipList {
	return &SkipList{
		maxLevel:    maxLevel,
		probability: probability,
		header: &Node{
			key:     0,
			value:   nil,
			forward: make([]*Node, maxLevel),
		},
		level:  1,
		length: 0,
	}
}

func (sl *SkipList) Find(key int64) interface{} {
	x := sl.header
	// 从当前最顶层开始
	for i := sl.level - 1; i >= 0; i-- {
		// 遍历每一层的链表
		// 如果比查找的值要小, 就继续遍历; 否则就遍历下一层
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
	}
	// 遍历到最底层, 进行比较
	x = x.forward[0]
	if x != nil && x.key == key {
		return x.value
	}
	return nil
}

func (sl *SkipList) Insert(key int64, value interface{}) {
	updateNodes := make([]*Node, sl.maxLevel)
	// 从当前最顶层开始
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		// 遍历每一层的链表
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		// 找到每一层应该插入的位置
		updateNodes[i] = x
	}
	// 遍历到最底层, 进行比较
	x = x.forward[0]
	// 如果key已经存在, 则更新value
	if x != nil && x.key == key {
		x.value = value
		return
	}
	// 如果不存在, 则先选择索引应该插入的层级
	newLevel := sl.randomLevel()
	// 如果新增索引层级大于当前层级, 则新建一层索引
	if newLevel > sl.level {
		// 从level到newLevel之间
		for i := sl.level; i < newLevel; i++ {
			// 新增需要插入的位置, 插入的位置都是在头节点之后
			updateNodes[i] = sl.header
		}
		// 更新当前层级
		sl.level = newLevel
	}
	newNode := &Node{key: key, value: value, forward: make([]*Node, newLevel)}
	// 更新每一层需要插入元素
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = updateNodes[i].forward[i]
		updateNodes[i].forward[i] = newNode
	}
	sl.length++
}

func (sl *SkipList) Delete(key int64) {
	updateNodes := make([]*Node, sl.maxLevel)
	// 从当前最顶层开始
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		// 遍历每一层的链表
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		// 找到每一层应该插入的位置
		updateNodes[i] = x
	}
	// 遍历到最底层, 进行比较
	x = x.forward[0]
	// 如果key存在，则进行删除
	if x != nil && x.key == key {
		// 从第1层开始删除
		for i := 0; i < sl.level; i++ {
			// 遍历每层链表, 直到遍历链表尾部退出
			if updateNodes[i].forward[i] != nil && updateNodes[i].forward[i].key != x.key {
				break
			}
			// 如果链表中有这个元素，则删除, Golang GC 不需要显示删除
			if updateNodes[i].forward[i] != nil {
				updateNodes[i].forward[i] = x.forward[i]
			}
		}
		// 判断一下删除之后每层元素是否只有头节点
		for sl.level >= 1 && sl.header.forward[sl.level-1] == nil {
			sl.level--
		}
		sl.length--
	}
}

func (sl *SkipList) randomLevel() int {
	newLevel := 1
	for sl.randP() >= sl.probability && newLevel < sl.maxLevel {
		newLevel++
	}
	return newLevel
}

func (sl *SkipList) randP() float32 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Float32()
}
