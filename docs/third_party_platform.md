# 第三方平台开发指南

## 开发新平台服务

1. 创建服务结构体
```go
package thirdPlatform

import (
    "github.com/zgsm/review-manager/pkg/httpclient"
)

// MyPlatformUser 用户数据结构
type MyPlatformUser struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    // 其他字段...
}

// MyPlatformService 自定义平台服务
type MyPlatformService struct {
    *Service // 继承基础Service
}

// NewMyPlatformService 创建服务实例
func NewMyPlatformService(clientConfig *httpclient.HttpServiceConfig) (*MyPlatformService, error) {
    client, err := httpclient.NewClient(clientConfig)
    if err != nil {
        return nil, err
    }

    return &MyPlatformService{
        Service: &Service{
            client: client,
        },
    }, nil
}
```

2. 实现服务方法
```go
// GetUser 获取用户信息
func (s *MyPlatformService) GetUser(ctx context.Context, userID string) (*MyPlatformUser, error) {
    var response struct {
        Code    int           `json:"code"`
        Message string        `json:"message"`
        Data    MyPlatformUser `json:"data"`
    }

    err := s.client.GetJSON(ctx, "/users/"+userID, nil, &response)
    if err != nil {
        return nil, fmt.Errorf("获取用户失败: %w", err)
    }

    if response.Code != 0 {
        return nil, fmt.Errorf("API错误: %s", response.Message)
    }

    return &response.Data, nil
}
```

3. 注册服务
在init.go中添加服务注册:
```go
func InitHTTPClient() error {
    // 原有代码...
    
    // 添加新服务
    myPlatformConfig, err := GetServiceConfig("my_platform")
    if err != nil {
        return err
    }
    myPlatformService, err := NewMyPlatformService(myPlatformConfig)
    if err != nil {
        return err
    }

    serverManager = &HttpServices{
        Example: *exampleService,
        MyPlatform: *myPlatformService, // 添加新服务
    }

    return nil
}
```

## HTTP请求发送指南

1. 基础请求方法
```go
// GET请求
err := s.client.GetJSON(ctx, "/path", params, &response)

// POST请求
err := s.client.PostJSON(ctx, "/path", requestBody, nil, &response)

// PUT请求
err := s.client.PutJSON(ctx, "/path", requestBody, nil, &response)

// DELETE请求
err := s.client.Delete(ctx, "/path", nil, &response)
```

2. 请求参数处理
```go
// 查询参数
params := map[string]string{
    "page": "1",
    "size": "10",
}

// 请求体
requestBody := struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}{
    Name: "John",
    Age:  30,
}
```

3. 响应处理
```go
var response struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    YourDataType `json:"data"`
}

if response.Code != 0 {
    return fmt.Errorf("API错误: %s", response.Message)
}
```

## 配置说明

1. 在config.yaml中添加服务配置
```yaml
httpclient:
  services:
    my_platform:
      base_url: "https://api.myplatform.com"
      timeout: 30
      max_retries: 3
      auth_type: "token"
      token: "your_api_token"
```

2. 配置参数说明
| 参数 | 说明 |
|---|---|
| base_url | 服务基础URL |
| timeout | 请求超时时间(秒) |
| max_retries | 最大重试次数 |
| auth_type | 认证类型(basic/token) |
| token | API访问令牌 |