# FlowGo Command Pattern Implementation

## 概述

本文档详细说明了 FlowGo 中 Command 模式的实现，这是从 Flowable/Activiti 移植过来的核心架构设计。

## 什么是 Command 模式？

Command 模式是一种行为设计模式，它将请求封装为对象，从而使你可以用不同的请求对客户端进行参数化。

### 在 Flowable 中的应用

在 Flowable/Activiti 中，所有对引擎的操作都被封装为 Command：

```java
// Java - Flowable
public interface Command<T> {
    T execute(CommandContext commandContext);
}

// 执行命令
processEngineConfiguration.getCommandExecutor()
    .execute(new DeployCmd(deployment));
```

### 在 FlowGo 中的实现

我们使用 Go 的泛型来实现类型安全的 Command：

```go
// Go - FlowGo
type Command[T any] interface {
    Execute(ctx context.Context, commandContext *CommandContext) (T, error)
}

// 执行命令
result, err := engine.Execute(ctx, command)
```

## 核心组件

### 1. Command - 命令接口

**位置**: `engine/command.go`

```go
type Command[T any] interface {
    Execute(ctx context.Context, commandContext *CommandContext) (T, error)
}
```

**特点**:
- 使用 Go 泛型实现类型安全
- 接受 `context.Context` 用于取消和超时
- 返回特定类型的结果和错误

### 2. CommandContext - 执行上下文

**位置**: `engine/command.go`

```go
type CommandContext struct {
    Context    context.Context           // Go 上下文
    Engine     *ProcessEngineImpl        // 引擎实例
    Session    interface{}               // 数据库会话
    Attributes map[string]interface{}    // 自定义属性
    Exception  error                     // 异常
    Result     interface{}               // 结果
}
```

**用途**:
- 为命令提供执行环境
- 管理数据库会话和事务
- 存储命令间共享的数据
- 追踪执行状态

### 3. CommandExecutor - 命令执行器

**位置**: `engine/command_executor.go`

```go
type CommandExecutor interface {
    Execute(ctx context.Context, command Command[any]) (any, error)
}
```

**实现**: `CommandExecutorImpl`

**职责**:
- 管理拦截器链
- 执行命令
- 处理异常

### 4. CommandInterceptor - 命令拦截器

**位置**: `engine/interceptor.go`

```go
type CommandInterceptor interface {
    Execute(ctx context.Context, command Command[any], next CommandExecutor) (any, error)
    GetNext() CommandInterceptor
    SetNext(next CommandInterceptor)
}
```

**内置拦截器**:

#### LoggingInterceptor
记录命令执行日志和耗时

```go
[FlowGo] Executing command: *commands.DeployCommand
[FlowGo] Command *commands.DeployCommand completed successfully in 15ms
```

#### TransactionInterceptor
管理数据库事务
- 开始事务
- 提交成功
- 失败回滚

#### ContextInterceptor
管理 CommandContext 生命周期
- 创建上下文
- 注入到 Go context
- 清理资源

#### RetryInterceptor
提供重试机制
- 可配置重试次数
- 可配置重试间隔
- 区分可重试错误

#### CommandInvoker
最终执行命令
- 获取 CommandContext
- 调用 Command.Execute

## 拦截器链执行流程

```
┌──────────────────────────────────────┐
│     CommandExecutor.Execute()        │
└────────────────┬─────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────┐
│      LoggingInterceptor              │  ← 记录开始时间
│  [FlowGo] Executing command...      │
└────────────────┬─────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────┐
│      RetryInterceptor (可选)         │  ← 提供重试逻辑
└────────────────┬─────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────┐
│    TransactionInterceptor            │  ← 开始事务
└────────────────┬─────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────┐
│      ContextInterceptor              │  ← 创建 CommandContext
└────────────────┬─────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────┐
│        CommandInvoker                │  ← 执行实际命令
│    command.Execute(ctx, cmdCtx)     │
└────────────────┬─────────────────────┘
                 │
                 ▼
           返回结果/错误
                 │
                 ▼
     逆序通过拦截器链返回
```

## 具体 Command 实现

### DeployCommand - 部署流程定义

**位置**: `engine/commands/deploy_command.go`

```go
type DeployCommand struct {
    DeploymentName   string
    Category         string
    TenantID         string
    ResourceName     string
    ResourceContent  []byte
}

func (c *DeployCommand) Execute(ctx context.Context, commandContext *CommandContext) (*repository.Deployment, error) {
    // 验证输入
    // 获取 RepositoryService
    // 执行部署
    // 返回结果
}
```

**使用示例**:
```go
cmd := commands.NewDeployCommand("My Process", "process.json", content)
result, err := engine.Execute(ctx, cmd)
deployment := result.(*repository.Deployment)
```

### StartProcessInstanceCommand - 启动流程实例

**位置**: `engine/commands/start_process_instance_command.go`

```go
type StartProcessInstanceCommand struct {
    ProcessDefinitionID  string
    ProcessDefinitionKey string
    BusinessKey          string
    Variables            map[string]interface{}
}
```

**特点**:
- 支持按 ID 或 Key 启动
- 支持业务键（Business Key）
- 验证流程定义状态（是否挂起）
- 记录到历史表

### CompleteTaskCommand - 完成任务

**位置**: `engine/commands/complete_task_command.go`

```go
type CompleteTaskCommand struct {
    TaskID    string
    Variables map[string]interface{}
}
```

**执行流程**:
1. 验证任务存在
2. 检查任务状态（是否挂起）
3. 设置输出变量
4. 完成任务
5. 记录历史
6. 触发流程继续执行

### ClaimTaskCommand - 认领任务

**位置**: `engine/commands/claim_task_command.go`

```go
type ClaimTaskCommand struct {
    TaskID string
    UserID string
}
```

**验证逻辑**:
- 任务是否已被其他用户认领
- 用户是否是候选人
- 记录认领时间

## 使用示例

### 基本用法

```go
package main

import (
    "context"
    "github.com/muixstudio/flowgo/engine"
    "github.com/muixstudio/flowgo/engine/commands"
)

func main() {
    // 创建引擎
    engine, _ := engine.NewProcessEngineBuilder().Build()
    engine.Start(context.Background())
    defer engine.Stop(context.Background())

    ctx := context.Background()

    // 1. 部署流程
    deployCmd := commands.NewDeployCommand("My Process", "process.json", content)
    deployment, _ := engine.Execute(ctx, deployCmd)

    // 2. 启动实例
    startCmd := commands.NewStartProcessInstanceByKeyCommand("my-process", variables)
    instance, _ := engine.Execute(ctx, startCmd)

    // 3. 认领任务
    claimCmd := commands.NewClaimTaskCommand(taskID, "user123")
    _, _ = engine.Execute(ctx, claimCmd)

    // 4. 完成任务
    completeCmd := commands.NewCompleteTaskCommand(taskID, outputVars)
    _, _ = engine.Execute(ctx, completeCmd)
}
```

### 自定义 Command

```go
package mycommands

import (
    "context"
    "github.com/muixstudio/flowgo/engine"
)

// 自定义命令：批量启动流程实例
type BatchStartProcessCommand struct {
    ProcessDefinitionKey string
    InstanceCount        int
    VariablesTemplate    map[string]interface{}
}

func (c *BatchStartProcessCommand) Execute(
    ctx context.Context,
    commandContext *engine.CommandContext,
) ([]string, error) {
    
    instanceIDs := make([]string, 0, c.InstanceCount)
    runtimeService := commandContext.Engine.GetRuntimeService()

    for i := 0; i < c.InstanceCount; i++ {
        // 为每个实例准备变量
        vars := make(map[string]interface{})
        for k, v := range c.VariablesTemplate {
            vars[k] = v
        }
        vars["batchIndex"] = i

        // 启动实例
        instance, err := runtimeService.StartProcessInstanceByKey(
            ctx, c.ProcessDefinitionKey, vars)
        if err != nil {
            return nil, err
        }
        
        instanceIDs = append(instanceIDs, instance.ID)
    }

    return instanceIDs, nil
}
```

### 自定义拦截器

```go
package myinterceptors

import (
    "context"
    "time"
    "github.com/muixstudio/flowgo/engine"
)

// 性能监控拦截器
type PerformanceInterceptor struct {
    engine.BaseCommandInterceptor
    slowThreshold time.Duration
}

func (i *PerformanceInterceptor) Execute(
    ctx context.Context,
    command engine.Command[any],
    next engine.CommandExecutor,
) (any, error) {
    
    start := time.Now()
    result, err := i.next.Execute(ctx, command, next)
    duration := time.Since(start)

    // 记录慢命令
    if duration > i.slowThreshold {
        log.Printf("[SLOW] Command %T took %v", command, duration)
    }

    return result, err
}

// 使用自定义拦截器
executor := engine.NewDefaultCommandExecutorBuilder(engineImpl).
    AddInterceptor(&PerformanceInterceptor{slowThreshold: 100 * time.Millisecond}).
    Build()
```

## 与 Flowable 的对比

| 特性 | Flowable (Java) | FlowGo (Go) |
|------|-----------------|-------------|
| Command 接口 | `Command<T>` | `Command[T]` (泛型) |
| 返回类型 | 需要类型转换 | 类型安全 |
| 上下文 | `CommandContext` | `CommandContext` + Go `context.Context` |
| 拦截器 | `CommandInterceptor` | `CommandInterceptor` |
| 事务管理 | Spring `@Transactional` | `TransactionInterceptor` |
| 异步执行 | `AsyncExecutor` | Go goroutines + channels |
| 错误处理 | Java Exceptions | Go error values |

## 优势

### 1. 统一的执行模型
所有操作都是 Command，执行方式一致。

### 2. 关注点分离
业务逻辑（Command）和横切关注点（Interceptor）分离。

### 3. 可测试性
Command 可以独立测试，无需完整的引擎环境。

### 4. 可扩展性
轻松添加新的 Command 和 Interceptor。

### 5. 事务一致性
通过 TransactionInterceptor 自动管理事务。

### 6. 可观测性
通过 LoggingInterceptor 统一记录所有操作。

### 7. 类型安全
Go 泛型确保编译时类型检查。

## 最佳实践

### 1. Command 应该是幂等的
尽可能使 Command 可以安全地重试。

### 2. 使用有意义的命名
```go
✅ Good: StartProcessInstanceByKeyCommand
❌ Bad: StartCmd
```

### 3. 在 Command 中进行验证
```go
func (c *DeployCommand) Execute(...) {
    if c.ResourceContent == nil {
        return nil, fmt.Errorf("content cannot be nil")
    }
    // ...
}
```

### 4. 合理使用 CommandContext
```go
// 存储用户信息
commandContext.SetAttribute("userId", userId)

// 在后续 Interceptor 中获取
userID := commandContext.GetAttribute("userId")
```

### 5. 错误应该包含上下文
```go
return nil, fmt.Errorf("failed to start process '%s': %w", key, err)
```

## 未来增强

- [ ] Command 序列化（用于持久化和远程执行）
- [ ] Command 队列（异步执行）
- [ ] Command 审计日志
- [ ] Command 调度（定时执行）
- [ ] 分布式 Command 执行
- [ ] Command 补偿机制（Saga 模式）
- [ ] Command 批处理优化

## 总结

FlowGo 的 Command 模式实现提供了：
- ✅ 清晰的架构
- ✅ 强大的扩展性
- ✅ 类型安全
- ✅ 易于测试
- ✅ 统一的执行模型
- ✅ 完善的拦截器机制

这是构建企业级工作流引擎的坚实基础。
