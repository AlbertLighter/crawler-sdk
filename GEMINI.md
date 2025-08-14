# Golang Crawler SDK

## 1. 项目目标

创建一个模块化、可扩展的Golang爬虫SDK，旨在简化从抖音、小红书等多个平台抓取数据的过程。该SDK将作为其他Golang项目的依赖库。

## 2. 核心设计原则

*   **模块化:** 每个平台（如抖音、小红书）都是一个独立的模块，易于维护和扩展。
*   **可扩展性:** 添加对新平台的支持应该简单直观。
*   **易用性:** 为最终用户提供一个干净、简洁的API。

## 3. 项目架构

```
/crawler-sdk
|-- go.mod
|-- pkg/
|   `-- http/
|       `-- client.go  // Resty v3 的封装，处理通用HTTP逻辑
|-- platforms/
|   |-- sdk/
|   |   `-- sdk.go     // 定义所有平台都需要实现的通用SDK接口
|   |-- douyin/
|   |   `-- douyin.go  // 抖音平台的具体实现
|   `-- xhs/
|       `-- xhs.go     // 小红书平台的具体实现
|-- examples/
|   `-- main.go        // SDK使用示例
`-- GEMINI.md
```

### 3.1. 目录说明

*   **`pkg/`**: 存放项目内部共享的包。
    *   **`http/`**: 对 `resty v3` 进行封装，统一处理客户端的通用逻辑，例如：设置User-Agent、处理Cookies、管理代理、请求重试等。
*   **`platforms/`**: 存放各个平台的核心实现。
    *   **`sdk/`**: 定义了 `SDK` 的核心接口。这是实现多平台支持的关键，它确保了所有平台都提供一致的API。
    *   **`douyin/`**: 抖音平台的具体实现，它将实现 `sdk.SDK` 接口。
    *   **`xhs/`**: 小红书平台的具体实现，同样实现 `sdk.SDK` 接口。
*   **`examples/`**: 提供如何使用本SDK的示例代码。

## 4. 工作流程

### 4.1. 如何添加一个新平台 (例如: `weibo`)

1.  在 `platforms/` 目录下创建一个新目录: `weibo`。
2.  在 `platforms/weibo/` 目录下创建一个 `weibo.go` 文件。
3.  在 `weibo.go` 中，创建一个结构体并实现 `sdk.SDK` 接口中定义的所有方法。
4.  在 `platforms/weibo/` 目录中添加该平台特定的其他逻辑。

### 4.2. 开发步骤

1.  **定义接口**: 在 `platforms/sdk/sdk.go` 中定义通用的 `SDK` 接口。
2.  **封装HTTP客户端**: 在 `pkg/http/client.go` 中封装 `resty v3`。
3.  **实现平台**: 开始在 `douyin` 和 `xhs` 目录中实现具体的平台逻辑。
4.  **编写示例**: 在 `examples/` 目录中提供清晰的SDK用法示例。
5.  **单元测试**: 为关键功能添加单元测试，确保代码质量和稳定性。
