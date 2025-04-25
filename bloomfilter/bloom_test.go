package bloomfilter

import (
	"testing"
)

func TestBloomFilter(t *testing.T) {
	// 创建一个布隆过滤器，预期存储1000个元素，误判率为0.01
	bf := NewBloomFilter(1000, 0.01)

	// 测试添加和查询
	testCases := []struct {
		name     string
		item     []byte
		expected bool
	}{
		{"测试1", []byte("hello"), true},
		{"测试2", []byte("world"), true},
		{"测试3", []byte("golang"), true},
	}

	// 添加元素
	for _, tc := range testCases {
		bf.Add(tc.item)
	}

	// 验证已添加的元素
	for _, tc := range testCases {
		if !bf.Contains(tc.item) {
			t.Errorf("%s: 期望找到元素 %s，但未找到", tc.name, tc.item)
		}
	}

	// 验证未添加的元素（可能存在误判）
	nonExistentItems := []struct {
		name string
		item []byte
	}{
		{"不存在1", []byte("notexist1")},
		{"不存在2", []byte("notexist2")},
	}

	for _, tc := range nonExistentItems {
		// 由于布隆过滤器的特性，这里的结果可能是true（误判）或false
		_ = bf.Contains(tc.item)
	}
}

func TestOptimalSize(t *testing.T) {
	size := calculateOptimalSize(1000, 0.01)
	if size <= 0 {
		t.Error("计算出的位数组大小应该大于0")
	}
}

func TestOptimalHashFunctions(t *testing.T) {
	numHashes := calculateOptimalHashFunctions(1000, 10000)
	if numHashes <= 0 {
		t.Error("计算出的哈希函数数量应该大于0")
	}
}
