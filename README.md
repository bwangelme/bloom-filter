# Golang 布隆过滤器实现

这是一个用 Golang 实现的布隆过滤器（Bloom Filter）数据结构。布隆过滤器是一种空间效率高的概率性数据结构，用于判断一个元素是否在集合中。

## 特性

- 空间效率高：使用位数组存储数据
- 时间效率高：添加和查询操作的时间复杂度都是 O(k)，其中 k 是哈希函数的数量
- 可配置性：可以根据预期元素数量和期望的误判率来优化性能
- 使用 FNV 哈希算法：提供良好的哈希分布

## 安装

```bash
go get github.com/bwangelme/bloom-filter
```

## 使用方法

### 创建布隆过滤器

```go
import "github.com/bwangelme/bloomfilter"

// 创建一个布隆过滤器
// 参数1：预期要存储的元素数量
// 参数2：期望的误判率（例如：0.01 表示 1% 的误判率）
bf := bloomfilter.NewBloomFilter(1000, 0.01)
```

### 添加元素

```go
// 添加元素到布隆过滤器
bf.Add([]byte("hello"))
```

### 查询元素

```go
// 检查元素是否可能在布隆过滤器中
exists := bf.Contains([]byte("hello"))
```

## 注意事项

1. 布隆过滤器有以下特性：
   - 不会漏报（false negative）：如果过滤器说元素不存在，那么元素一定不存在
   - 可能误报（false positive）：如果过滤器说元素存在，元素可能不存在
   - 不支持删除操作：一旦添加元素，就不能删除

2. 误判率与空间使用：
   - 误判率越低，需要的空间越大
   - 预期元素数量越多，需要的空间越大

## 示例代码

```go
package main

import (
    "fmt"
    "github.com/bwangelme/bloomfilter"
)

func main() {
    // 创建布隆过滤器
    bf := bloomfilter.NewBloomFilter(1000, 0.01)

    // 添加一些元素
    items := []string{"hello", "world", "golang"}
    for _, item := range items {
        bf.Add([]byte(item))
    }

    // 检查元素是否存在
    for _, item := range items {
        exists := bf.Contains([]byte(item))
        fmt.Printf("检查 %s: %v\n", item, exists)
    }

    // 检查不存在的元素
    nonExistent := "notexist"
    exists := bf.Contains([]byte(nonExistent))
    fmt.Printf("检查 %s: %v\n", nonExistent, exists)
}
```

## 测试

运行测试：

```bash
go test ./bloomfilter
```

## 性能考虑

- 添加和查询操作的时间复杂度：O(k)，其中 k 是哈希函数的数量
- 空间复杂度：O(m)，其中 m 是位数组的大小
- 哈希函数的数量会根据预期元素数量和位数组大小自动优化

## 许可证

MIT License
