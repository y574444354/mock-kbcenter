# Go WebServer 国际化(i18n)开发指南

## 1. 概述

本项目使用`go-i18n`库实现了国际化支持，所有用户可见的文本（包括日志和API响应）都通过i18n系统处理。

## 2. 语言文件结构

语言文件存放在`i18n/locales/`目录下，每个语言一个YAML文件，例如：
- `zh-CN.yaml` - 简体中文
- `en.yaml` - 英文

文件格式示例：
```yaml
common:
  success: "成功"
  error: "错误"
user:
  login:
    success: "登录成功"
    failed: "登录失败"
```

## 3. 添加新语言

1. 在`i18n/locales/`目录下创建新的YAML文件，如`fr.yaml`（法语）
2. 在`config/config.yaml`中添加新语言的配置
3. 在`i18n/i18n.go`的`InitI18n`函数中会自动加载新语言文件

## 4. 使用翻译功能

### 在API处理器中使用

```go
// 使用消息ID而不是硬编码文本
api.BadRequest(c, "user.login.invalid_params")

// 带模板数据的翻译
data := map[string]interface{}{"count": 5}
message := i18n.Translate("items.count", locale, data)
```

### 在日志中使用

```go
// 使用i18nlogger包记录国际化日志
locale := i18nlogger.GetLocaleFromContext(c)
i18nlogger.Error("user.login.failed", locale, nil, "error", err)
```

## 5. 中间件集成

`I18n`中间件会自动从请求中获取语言设置，并添加到Gin上下文中：

```go
// 注册中间件
router.Use(middleware.I18n())

// 获取当前语言
locale := c.GetString("locale")
```

## 6. 最佳实践

1. 所有用户可见的文本都应该使用消息ID，而不是硬编码
2. 消息ID应该按功能模块组织，如`user.login.success`
3. 添加新功能时，同时更新所有语言文件
4. 日志消息也应该国际化，使用i18nlogger包

## 7. 测试国际化

1. 通过`locale`查询参数切换语言：`?locale=en`
2. 通过`Accept-Language`头指定语言
3. 通过Cookie设置语言

## 8. 示例

```go
// 获取翻译
message := i18n.Translate("welcome.message", "zh-CN", nil)

// 记录日志
i18nlogger.Info("app.startup", "en", nil)
```

## 9. 注意事项

1. 确保所有消息ID都有对应的翻译
2. 更新语言文件后不需要重启服务
3. 默认语言在配置文件中设置