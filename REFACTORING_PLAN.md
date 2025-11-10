# FlowGo 项目重构方案

## 设计原则

1. **根包作为入口** - 对外暴露主要 API
2. **internal/** - 内部实现细节，Go 编译器保证不可从外部导入
3. **api/** - 公共接口定义，可被外部导入
4. **pkg/** - 可复用的公共工具库（与流程引擎主题无关）
5. **最小暴露** - 只暴露必要的接口和类型

## 新的目录结构

```
flowgo/
├── engine.go                      # 主入口 - ProcessEngine 接口
├── config.go                      # 配置和 Builder
├── errors.go                      # 公共错误类型
│
├── api/                           # 对外暴露的服务接口
│   ├── repository/               
│   │   ├── service.go             # Repository Service 接口
│   │   └── types.go               # 公共类型（Deployment, ProcessDefinition等）
│   ├── runtime/
│   │   ├── service.go             # Runtime Service 接口  
│   │   └── types.go               # 公共类型（ProcessInstance, Execution等）
│   ├── task/
│   │   ├── service.go             # Task Service 接口
│   │   └── types.go               # 公共类型（Task, Comment, Attachment等）
│   └── history/
│       ├── service.go             # History Service 接口
│       └── types.go               # 公共类型（Historic*等）
│
├── internal/                      # 内部实现（外部无法导入）
│   ├── engine/
│   │   ├── engine.go              # ProcessEngine 实现
│   │   ├── command/               # Command 模式
│   │   │   ├── command.go         # Command 接口
│   │   │   ├── context.go         # CommandContext
│   │   │   ├── executor.go        # CommandExecutor
│   │   │   └── interceptor/       # 拦截器
│   │   │       ├── logging.go
│   │   │       ├── transaction.go
│   │   │       ├── context.go
│   │   │       └── retry.go
│   │   └── commands/              # 具体命令实现
│   │       ├── deploy.go
│   │       ├── start_instance.go
│   │       ├── complete_task.go
│   │       └── claim_task.go
│   │
│   ├── repository/
│   │   ├── service.go             # 实现 api/repository.Service
│   │   ├── deployment.go          # Deployment 相关逻辑
│   │   └── definition.go          # ProcessDefinition 相关逻辑
│   │
│   ├── runtime/
│   │   ├── service.go             # 实现 api/runtime.Service
│   │   ├── instance.go            # ProcessInstance 管理
│   │   ├── execution.go           # Execution 管理
│   │   └── variable.go            # 变量管理
│   │
│   ├── task/
│   │   ├── service.go             # 实现 api/task.Service
│   │   ├── task.go                # Task 管理
│   │   ├── assignment.go          # 任务分配逻辑
│   │   └── comment.go             # 评论和附件
│   │
│   ├── history/
│   │   ├── service.go             # 实现 api/history.Service
│   │   └── recorder.go            # 历史记录器
│   │
│   ├── persistence/               # 数据持久化层
│   │   ├── db.go                  # 数据库连接管理
│   │   ├── transaction.go         # 事务管理
│   │   └── dao/                   # Data Access Objects
│   │       ├── deployment.go
│   │       ├── process_instance.go
│   │       └── task.go
│   │
│   └── executor/                  # 流程执行器
│       ├── engine.go               # 执行引擎
│       ├── behavior/               # 节点行为
│       │   ├── user_task.go
│       │   ├── service_task.go
│       │   └── gateway.go
│       └── async/                  # 异步执行
│           └── job_executor.go
│
├── pkg/                           # 可复用的公共包
│   ├── expression/                # 表达式引擎
│   │   ├── parser.go
│   │   ├── evaluator.go
│   │   └── functions.go           # 内置函数（now(), duration()等）
│   │
│   ├── parser/                    # 流程定义解析器
│   │   ├── json_parser.go         # JSON 解析
│   │   └── validator.go           # 验证器
│   │
│   └── util/                      # 工具函数
│       ├── id_generator.go        # ID 生成器
│       └── time.go                # 时间处理
│
├── schema/                        # JSON Schema
│   ├── process_definition.schema.json
│   └── README.md
│
├── examples/                      # 示例代码
│   ├── basic/
│   │   └── main.go
│   ├── command_pattern/
│   │   └── main.go
│   └── processes/
│       ├── leave_approval.json
│       └── expense_approval.json
│
├── docs/                          # 文档
│   ├── architecture.md
│   ├── command_pattern.md
│   └── api/                       # API 文档
│
├── README.md
├── go.mod
└── go.sum
```

## 包的职责划分

### 根包 (flowgo)

**暴露内容**:
- `ProcessEngine` 接口
- `Configuration` 和 `Builder`
- 工厂函数: `NewProcessEngine()`, `NewProcessEngineBuilder()`
- 公共错误类型

**不暴露**:
- 任何实现细节
- Command 模式相关
- 数据持久化逻辑

```go
// 用户这样使用
import "github.com/muixstudio/flowgo"

engine, err := flowgo.NewProcessEngineBuilder().
    WithEngineName("my-engine").
    WithDatabase("postgres", "...").
    Build()
```

### api/ - 服务接口包

**暴露内容**:
- 各个服务的接口定义
- 领域模型（Deployment, ProcessInstance, Task 等）
- Query 构建器

**示例**:
```go
// api/repository/service.go
package repository

type Service interface {
    CreateDeployment() *DeploymentBuilder
    GetProcessDefinition(ctx context.Context, id string) (*ProcessDefinition, error)
    // ...
}
```

### internal/ - 内部实现

**特点**:
- Go 编译器保证外部包无法导入
- 包含所有实现细节
- 可以随意重构，不影响公共 API

**职责**:
- ProcessEngine 的具体实现
- Command 模式实现
- 服务实现
- 数据访问层
- 流程执行引擎

### pkg/ - 公共工具包

**特点**:
- 可以被外部项目导入
- 与流程引擎主题无关
- 通用、可复用

**示例**:
```go
// pkg/expression/evaluator.go
package expression

// Evaluate 计算表达式
func Evaluate(expr string, variables map[string]interface{}) (interface{}, error) {
    // ...
}
```

## 迁移步骤

### 阶段 1: 创建新结构

1. ✅ 创建 `api/` 目录和服务接口
2. ✅ 创建根包的 `engine.go` 和 `config.go`
3. ⬜ 创建 `internal/` 目录
4. ⬜ 创建 `pkg/` 目录

### 阶段 2: 迁移现有代码

1. ⬜ 将 `engine/` 下的实现移到 `internal/engine/`
2. ⬜ 将 `repository/`, `runtime/`, `task/`, `history/` 的实现移到 `internal/`
3. ⬜ 将接口定义提取到 `api/` 下
4. ⬜ 删除旧的根目录文件（engine.go, engine_impl.go）

### 阶段 3: 更新导入路径

1. ⬜ 更新 examples
2. ⬜ 更新文档
3. ⬜ 更新测试

### 阶段 4: 添加新功能

1. ⬜ 实现表达式引擎 (`pkg/expression/`)
2. ⬜ 实现流程执行引擎 (`internal/executor/`)
3. ⬜ 实现数据持久化层 (`internal/persistence/`)

## 使用示例对比

### 重构前

```go
import (
    "github.com/muixstudio/flowgo"
    "github.com/muixstudio/flowgo/repository"  // 直接导入实现
    "github.com/muixstudio/flowgo/runtime"
)

engine := flowgo.NewProcessEngineBuilder().Build()
repoService := engine.GetRepositoryService()
```

### 重构后

```go
import "github.com/muixstudio/flowgo"
// 不需要导入其他包，所有 API 通过 engine 获取

engine, _ := flowgo.NewProcessEngineBuilder().
    WithEngineName("my-engine").
    Build()

repoService := engine.GetRepositoryService()
// repoService 类型是 api/repository.Service 接口
// 用户不知道也不需要知道具体实现
```

## 好处

### 1. 清晰的 API 边界

✅ 用户只需导入 `github.com/muixstudio/flowgo`
✅ 所有公共 API 都在根包或 `api/` 下
✅ `internal/` 下的代码无法被外部导入

### 2. 易于维护

✅ 内部实现可以随意重构
✅ 不会破坏公共 API
✅ 代码组织清晰

### 3. 符合 Go 惯例

✅ 使用 `internal/` 包
✅ 根包作为入口
✅ `pkg/` 下是可复用的工具

### 4. 更好的封装

✅ Command 模式细节对用户透明
✅ 用户不需要了解拦截器
✅ 实现细节隐藏在 `internal/` 中

### 5. 灵活性

✅ 可以提供多种实现（如不同的数据库适配器）
✅ 可以通过插件扩展功能
✅ 便于单元测试（mock 接口）

## 下一步行动

1. **审查方案** - 确认这个结构符合项目需求
2. **创建目录** - 创建新的目录结构
3. **迁移代码** - 逐步迁移现有代码
4. **更新文档** - 更新所有文档和示例
5. **测试验证** - 确保所有功能正常

## 注意事项

⚠️ **Breaking Changes**: 这是一次大的重构，会破坏现有的导入路径
⚠️ **逐步迁移**: 建议分阶段完成，每个阶段都确保代码可运行
⚠️ **保持兼容**: 考虑在根包提供兼容性别名，方便迁移
