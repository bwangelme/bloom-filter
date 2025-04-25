# Golang 布隆过滤器实现

这是一个用 Golang 实现的布隆过滤器（Bloom Filter）数据结构，使用 Redis 作为存储后端。布隆过滤器是一种空间效率高的概率性数据结构，用于判断一个元素是否在集合中。

## 特性

- 基于 Redis 存储：使用 Redis 的位数组存储布隆过滤器的数据，支持分布式部署
- 空间效率高：使用位数组存储数据
- 时间效率高：添加和查询操作的时间复杂度都是 O(k)，其中 k 是哈希函数的数量
- 可配置性：可以根据预期元素数量和期望的误判率来优化性能
- 使用 FNV 哈希算法：提供良好的哈希分布
- Web API 支持：提供 RESTful API 接口
- 前端演示页面：直观的 Web 界面展示布隆过滤器的功能

## 系统要求

- Go 1.20 或更高版本
- Redis 6.0 或更高版本

## 安装

1. 克隆仓库：
```bash
git clone https://github.com/yourusername/bloom-filter.git
cd bloom-filter
```

2. 安装依赖：
```bash
go mod download
```

3. 配置 Redis：
   - 确保 Redis 服务器已启动
   - 修改 `config.yaml` 文件中的 Redis 配置：
     ```yaml
     redis:
       addr: "localhost:6379"    # Redis 服务器地址
       key: "bloomfilter"        # 用于存储位数组的键名
       password: ""              # Redis 密码（如果有）
       db: 0                     # Redis 数据库编号
     ```

4. 配置布隆过滤器参数：
   - 在 `config.yaml` 文件中设置：
     ```yaml
     bloom:
       expected_items: 1000      # 预期要存储的元素数量
       false_positive_rate: 0.01 # 期望的误判率（1%）
     ```

5. 配置服务器：
   - 在 `config.yaml` 文件中设置：
     ```yaml
     server:
       port: 8080               # 服务器端口
       host: "0.0.0.0"         # 服务器主机地址
     ```

## 运行

```bash
go run main.go
```

服务器启动后，可以通过浏览器访问 `http://localhost:8080` 来使用 Web 界面。

## API 接口

所有 API 接口都以 `/api` 为前缀：

### 添加元素

- 添加单个元素
  ```http
  POST /api/add
  Content-Type: application/json

  {
    "item": "要添加的元素"
  }
  ```

- 批量添加元素
  ```http
  POST /api/batch-add
  Content-Type: application/json

  {
    "items": ["元素1", "元素2", "元素3"]
  }
  ```

### 检查元素

- 检查单个元素
  ```http
  POST /api/check
  Content-Type: application/json

  {
    "item": "要检查的元素"
  }
  ```

- 批量检查元素
  ```http
  POST /api/batch/contains
  Content-Type: application/json

  {
    "items": ["元素1", "元素2", "元素3"]
  }
  ```

## 注意事项

1. Redis 相关：
   - 确保 Redis 服务器有足够的内存来存储位数组
   - 建议为布隆过滤器使用独立的 Redis 数据库或键名前缀
   - Redis 连接失败时服务将无法启动
   - Redis 中的数据是持久化的，重启服务不会丢失数据

2. 布隆过滤器特性：
   - 不会漏报（false negative）：如果过滤器说元素不存在，那么元素一定不存在
   - 可能误报（false positive）：如果过滤器说元素存在，元素可能不存在
   - 不支持删除操作：一旦添加元素，就不能删除

3. 性能考虑：
   - 误判率越低，需要的 Redis 存储空间越大
   - 预期元素数量越多，需要的存储空间越大
   - 哈希函数数量会影响操作的响应时间

## 性能指标

在默认配置下（预期 1000 个元素，1% 误判率）：
- 位数组大小：约 9.6 KB
- 哈希函数数量：7 个
- 单个元素添加/查询时间：< 5ms（取决于 Redis 网络延迟）

## 许可证

MIT License
