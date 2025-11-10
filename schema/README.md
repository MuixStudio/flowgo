# FlowGo Process Definition Schema

## 概述

FlowGo 使用 JSON 格式定义工作流程，替代 BPMN 的 XML 格式。本 Schema 参考了 BPMN 2.0 的核心概念，并针对 Go 语言进行了优化。

## 核心概念

### 1. 流程定义 (Process Definition)

顶层对象，包含完整的流程定义：

```json
{
  "id": "process-id",
  "name": "Process Name",
  "version": 1,
  "nodes": [...],
  "edges": [...],
  "variables": {...}
}
```

### 2. 节点类型 (Node Types)

#### 事件 (Events)
- **startEvent**: 流程开始事件
- **endEvent**: 流程结束事件
- **intermediateEvent**: 中间事件（消息、定时器等）
- **boundaryEvent**: 边界事件（附加到活动上的事件）

#### 任务 (Tasks)
- **userTask**: 用户任务，需要人工处理
- **serviceTask**: 服务任务，自动执行业务逻辑
- **scriptTask**: 脚本任务，执行脚本代码
- **callActivity**: 调用子流程
- **subProcess**: 嵌入式子流程

#### 网关 (Gateways)
- **exclusiveGateway**: 排他网关（XOR），只选择一条路径
- **parallelGateway**: 并行网关（AND），所有路径同时执行
- **inclusiveGateway**: 包容网关（OR），多条路径可选
- **eventBasedGateway**: 基于事件的网关

### 3. 序列流 (Edges/Sequence Flows)

连接节点的有向边：

```json
{
  "id": "flow-id",
  "source": "source-node-id",
  "target": "target-node-id",
  "condition": "${expression}",
  "isDefault": false
}
```

## 表达式语言

FlowGo 使用类似 SpEL (Spring Expression Language) 的表达式语法：

- **变量引用**: `${variableName}`
- **比较运算**: `${value > 10}`, `${status == 'approved'}`
- **逻辑运算**: `${condition1 && condition2}`, `${!flag}`
- **函数调用**: `${now()}`, `${duration('P2D')}`

## 用户任务属性

### 分配策略

1. **直接分配**: `assignee: "user123"`
2. **候选用户**: `candidateUsers: ["user1", "user2"]`
3. **候选组**: `candidateGroups: ["managers", "hr"]`
4. **表达式分配**: `assignee: "${processInitiator}"`

### 其他属性

- **formKey**: 关联的表单定义
- **dueDate**: 截止日期
- **priority**: 优先级（1-10）

## 服务任务

### 实现方式

1. **类实现**: `implementation: "com.example.MyService"`
2. **委托表达式**: `delegateExpression: "${myDelegate}"`
3. **表达式**: 直接使用表达式执行逻辑

### 异步执行

```json
{
  "async": true,
  "retries": 3,
  "retryInterval": "PT5M"
}
```

## 网关路由

### 排他网关示例

```json
{
  "id": "gateway1",
  "type": "exclusiveGateway"
}
```

出边需要配置条件：

```json
[
  {
    "source": "gateway1",
    "target": "task1",
    "condition": "${amount > 1000}"
  },
  {
    "source": "gateway1",
    "target": "task2",
    "condition": "${amount <= 1000}"
  },
  {
    "source": "gateway1",
    "target": "task3",
    "isDefault": true
  }
]
```

## 变量映射

### 输入映射 (inputMappings)

将流程变量映射到任务局部变量：

```json
{
  "inputMappings": {
    "localVar": "${processVar}",
    "amount": "${order.totalAmount}"
  }
}
```

### 输出映射 (outputMappings)

将任务结果映射回流程变量：

```json
{
  "outputMappings": {
    "approved": "${taskResult.isApproved}",
    "comment": "${taskResult.reviewComment}"
  }
}
```

## 事件定义

### 定时器事件

```json
{
  "type": "intermediateEvent",
  "properties": {
    "eventType": "timer",
    "eventDefinition": {
      "timerType": "duration",
      "timerValue": "PT2H"
    }
  }
}
```

### 消息事件

```json
{
  "type": "intermediateEvent",
  "properties": {
    "eventType": "message",
    "eventDefinition": {
      "messageName": "payment-received"
    }
  }
}
```

### 边界事件

```json
{
  "id": "timeout-boundary",
  "type": "boundaryEvent",
  "properties": {
    "eventType": "timer",
    "attachedTo": "user-task-1",
    "cancelActivity": true,
    "eventDefinition": {
      "timerType": "duration",
      "timerValue": "P3D"
    }
  }
}
```

## 扩展元素

使用 `extensionElements` 添加自定义属性：

```json
{
  "extensionElements": {
    "customProperty1": "value1",
    "tags": ["important", "finance"],
    "metadata": {
      "department": "HR"
    }
  }
}
```

## 完整示例

查看 `examples/leave_approval.json` 获取完整的请假审批流程示例。

## 与 BPMN 的对应关系

| BPMN 元素 | FlowGo JSON |
|-----------|-------------|
| Process | 根对象 |
| StartEvent | `{"type": "startEvent"}` |
| EndEvent | `{"type": "endEvent"}` |
| UserTask | `{"type": "userTask"}` |
| ServiceTask | `{"type": "serviceTask"}` |
| ExclusiveGateway | `{"type": "exclusiveGateway"}` |
| ParallelGateway | `{"type": "parallelGateway"}` |
| SequenceFlow | edges 数组中的对象 |
| ConditionExpression | `{"condition": "${expr}"}` |
| DataObject | variables 对象 |

## 验证

可以使用标准的 JSON Schema 验证器验证流程定义文件：

```bash
# 使用 ajv-cli 验证
ajv validate -s process_definition.schema.json -d ../examples/leave_approval.json
```
