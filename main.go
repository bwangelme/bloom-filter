package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/bwangelme/bloom-filter/bloomfilter"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// Config 配置结构
type Config struct {
	Redis struct {
		Addr     string `yaml:"addr"`
		Key      string `yaml:"key"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	Bloom struct {
		ExpectedItems     int     `yaml:"expected_items"`
		FalsePositiveRate float64 `yaml:"false_positive_rate"`
	} `yaml:"bloom"`
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
}

// Request 请求结构
type AddRequest struct {
	Item  string   `json:"item"`
	Items []string `json:"items"`
}

// Response 响应结构
type ContainsResponse struct {
	Results map[string]bool `json:"results"`
}

var (
	config *Config
	bf     *bloomfilter.BloomFilter
	once   sync.Once
)

// 初始化配置
func initConfig() error {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	config = &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

// 初始化布隆过滤器
func initBloomFilter() error {
	var err error
	once.Do(func() {
		bf, err = bloomfilter.NewBloomFilter(
			config.Bloom.ExpectedItems,
			config.Bloom.FalsePositiveRate,
			config.Redis.Addr,
			config.Redis.Key,
		)
	})
	return err
}

// 处理单个元素添加请求
func handleAdd(c *gin.Context) {
	var req struct {
		Item string `json:"item"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bf.Add([]byte(req.Item))
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 处理批量添加请求
func handleBatchAdd(c *gin.Context) {
	var req struct {
		Items []string `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, item := range req.Items {
		bf.Add([]byte(item))
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 处理检查元素是否存在的请求
func handleCheck(c *gin.Context) {
	var req struct {
		Item string `json:"item"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := bf.Contains([]byte(req.Item))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

// 处理批量查询元素的请求
func handleBatchContains(c *gin.Context) {
	var req AddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "items 不能为空"})
		return
	}

	results := make(map[string]bool)
	for _, item := range req.Items {
		exists, err := bf.Contains([]byte(item))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("查询元素失败: %v", err)})
			return
		}
		results[item] = exists
	}

	c.JSON(http.StatusOK, ContainsResponse{Results: results})
}

func main() {
	// 初始化配置
	if err := initConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化布隆过滤器
	if err := initBloomFilter(); err != nil {
		log.Fatalf("初始化布隆过滤器失败: %v", err)
	}
	defer bf.Close()

	// 创建 Gin 引擎
	r := gin.Default()

	// 设置静态文件服务
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// 设置路由
	r.POST("/api/add", handleAdd)
	r.POST("/api/check", handleCheck)
	r.POST("/api/batch-add", handleBatchAdd)
	r.POST("/api/batch/contains", handleBatchContains)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
