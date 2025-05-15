# 分支管理规范

## 分支策略
本项目采用主干开发模式(Trunk-based Development)，所有开发都基于master分支进行。

## 分支命名规范
1. 功能分支：`feature/[功能描述]`，例如 `feature/user-authentication`
2. 修复分支：`fix/[问题描述]`，例如 `fix/login-bug`
3. 发布分支：`release/[版本号]`，例如 `release/v1.0.0`
4. 热修复分支：`hotfix/[问题描述]`，例如 `hotfix/security-patch`

## 提交信息规范
提交信息应遵循以下格式：
```
<类型>(<范围>): <主题>

<正文>

<页脚>
```

### 类型
- feat: 新功能
- fix: 修复bug
- docs: 文档变更
- style: 代码格式变更
- refactor: 代码重构
- test: 测试相关
- chore: 构建或辅助工具变更

### 示例
```
feat(user): 添加用户登录功能

- 实现JWT认证
- 添加登录API

Closes #123
```

## 合并请求流程
1. 从master分支创建新分支
2. 开发完成后创建Pull Request
3. 至少需要1位核心成员Code Review
4. 通过CI/CD流水线检查
5. 合并到master分支