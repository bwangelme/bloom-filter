package bloomfilter

import (
	"context"
	"hash/fnv"
	"math"

	"github.com/redis/go-redis/v9"
)

// BloomFilter 表示布隆过滤器
type BloomFilter struct {
	client    *redis.Client
	key       string
	size      uint
	hashFuncs []func([]byte) uint
}

// NewBloomFilter 创建一个新的布隆过滤器
// expectedItems: 预期要存储的元素数量
// falsePositiveRate: 期望的误判率
// redisAddr: Redis 服务器地址
// key: Redis 中存储位数组的键名
func NewBloomFilter(expectedItems int, falsePositiveRate float64, redisAddr, key string) (*BloomFilter, error) {
	// 计算最优的位数组大小
	size := calculateOptimalSize(expectedItems, falsePositiveRate)

	// 计算最优的哈希函数数量
	numHashes := calculateOptimalHashFunctions(expectedItems, size)

	// 创建哈希函数
	hashFuncs := make([]func([]byte) uint, numHashes)
	for i := 0; i < numHashes; i++ {
		seed := i
		hashFuncs[i] = func(data []byte) uint {
			h := fnv.New64a()
			h.Write(data)
			h.Write([]byte{byte(seed)})
			return uint(h.Sum64() % uint64(size))
		}
	}

	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// 测试 Redis 连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	// 初始化位数组
	bf := &BloomFilter{
		client:    client,
		key:       key,
		size:      uint(size),
		hashFuncs: hashFuncs,
	}

	// 确保位数组已初始化
	exists, err := client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		// 使用 Redis 的 SETBIT 命令初始化位数组
		for i := uint(0); i < bf.size; i++ {
			if err := client.SetBit(ctx, key, int64(i), 0).Err(); err != nil {
				return nil, err
			}
		}
	}

	return bf, nil
}

// Add 向布隆过滤器中添加元素
func (bf *BloomFilter) Add(item []byte) error {
	ctx := context.Background()
	for _, hashFunc := range bf.hashFuncs {
		position := hashFunc(item)
		if err := bf.client.SetBit(ctx, bf.key, int64(position), 1).Err(); err != nil {
			return err
		}
	}
	return nil
}

// Contains 检查元素是否可能在布隆过滤器中
func (bf *BloomFilter) Contains(item []byte) (bool, error) {
	ctx := context.Background()
	for _, hashFunc := range bf.hashFuncs {
		position := hashFunc(item)
		bit, err := bf.client.GetBit(ctx, bf.key, int64(position)).Result()
		if err != nil {
			return false, err
		}
		if bit == 0 {
			return false, nil
		}
	}
	return true, nil
}

// Close 关闭 Redis 连接
func (bf *BloomFilter) Close() error {
	return bf.client.Close()
}

// calculateOptimalSize 计算最优的位数组大小
// 参数:
//   - n: 预期要存储的元素数量
//   - p: 期望的误判率（例如：0.01 表示 1% 的误判率）
//
// 返回值:
//   - 计算出的最优位数组大小（向上取整）
//
// 计算公式: m = -1 * (n * ln(p)) / (ln(2)^2)
// 说明:
//  1. 位数组大小 m 与预期元素数量 n 和误判率 p 有关
//  2. 当误判率 p 越小时，需要的位数组大小 m 越大
//  3. 当预期元素数量 n 越大时，需要的位数组大小 m 越大
//  4. 使用向上取整确保有足够的空间存储所有元素
func calculateOptimalSize(n int, p float64) int {
	// m = -1 * (n * ln(p)) / (ln(2)^2)
	m := -1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)
	return int(math.Ceil(m))
}

// calculateOptimalHashFunctions 计算最优的哈希函数数量
// 参数:
//   - n: 预期要存储的元素数量
//   - m: 位数组的大小
//
// 返回值:
//   - 计算出的最优哈希函数数量（向上取整）
//
// 计算公式: k = (m/n) * ln(2)
// 说明:
//  1. 哈希函数数量 k 与位数组大小 m 和预期元素数量 n 有关
//  2. 当位数组大小 m 越大时，可以使用更多的哈希函数
//  3. 当预期元素数量 n 越大时，应该使用更少的哈希函数
//  4. 使用向上取整确保有足够的哈希函数来降低误判率
//  5. 哈希函数数量会影响：
//     - 添加和查询操作的时间复杂度（与 k 成正比）
//     - 误判率（k 过大或过小都会增加误判率）
func calculateOptimalHashFunctions(n, m int) int {
	// k = (m/n) * ln(2)
	k := float64(m) / float64(n) * math.Log(2)
	return int(math.Ceil(k))
}
