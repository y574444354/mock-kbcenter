这是一个 go 项目
项目名为: github.com/zgsm/go-webserver

项目规范：
1、业务过程中涉及的常用结构体优先在 pkg/types/types.go 中定义，除非是非常少用的场景
2、项目中所有涉及到对外用户可见的文案，均需使用 i18n，例如报错内容、日志打印。i18n 定义见 i18n/i18n.go。禁止硬编码中文或英文，除非特殊情况，例如测试文件、文档说明、注释
3、第三方 api 对接，需在 pkg/thirdPlatform/ 中定义对应的 service，具体可参考 pkg/thirdPlatform/issue_manager_service.go 文件。而要获取对应的 service 类时，需通过 pkg/thirdPlatform/init.go 中定义的 GetServerManager 来统一获取。第三方 api 所需的配置在 config/config.go 中由 HTTPClient.Services 结构体进行扩展。第三方 api 中涉及到的跟第三方接口相关的结构体逻辑，应该收敛在该 service 中，避免要求调用方构造第三方特定数据结构
4、涉及到 http 请求时，应使用 pkg/httpclient/client.go 中定义的请求客户端
5、数据库模型定义在 internal/model/ 内（一般情况下模型都需要一个支持自增的 ID），dao 层为 internal/repository/（repository 依赖的 db 可直接使用 pkg/db 中定义的 GetDB 方法，具体参考 internal/repository/review_task_repository.go 文件），service 层为 internal/service（service 依赖的 repository，应在对应 repository 中定义实例化方法，具体可参考 internal/repository/review_task_repository.go 文件）。应限制 model 仅由 repository 进行操作，service 操作 repository 层，禁止外部文件直接调用 repository 与 model 层接口，对外应由 service 层统一提供数据业务操作与更高级的业务逻辑操作。对于具体的业务逻辑操作，也应该优先构造 service 类。
6、涉及到多种数据适配需求的，应在 internal/repository/ 层中，使用适配器模式支持
7、每当新增模型时，需同步在 cmd/dbtools/init.go 与 cmd/dbtools/migrate.go 中注册该模型，用于支持 db 初始化与数据迁移操作
8、涉及到 api 接口定义时，参考 api/v1/review_task_handler.go，注意需遵循 restful api 规范，每种业务资源有对应自己的 handler。response 需参考 api/v1/response.go 中定义的响应体生成方法，所有接口均需编写 swag 注释。
9、所有工具库定义可在 pkg/ 下查找，例如 httpclient、redis 等等
10、config 中的 httpclient.services 中应只配置第三方服务的 http 协议相关参数，若存在一些跟业务细节有关的配置参数，应单独配置，不要混入 httpclient 中
11、每当 model 发生字段变更时，需检查相应的 repository、service 层是否需要做相应的参数调整
12、pkg 中定义的通用包不应引入其他业务文件，保持通用性。有异常可考虑直接抛出异常，由上层调用层进行进一步处理
